// Ниже реализован сервис бронирования номеров в отеле. В предметной области
// выделены два понятия: Order — заказ, который включает в себя даты бронирования
// и контакты пользователя, и RoomAvailability — количество свободных номеров на
// конкретный день.
//
// Задание:
// - провести рефакторинг кода с выделением слоев и абстракций
// - применить best-practices там где это имеет смысл
// - исправить имеющиеся в реализации логические и технические ошибки и неточности
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	httpapi "yanmhlv/test-assignment/internal/api/http"
	"yanmhlv/test-assignment/internal/booking"
)

var (
	logger = log.Default()

	orderRepo        = booking.NewInMemoryOrderRepository()
	availabilityRepo = booking.NewInMemoryAvailabilityRepository()
	bookingService   = booking.NewBookingService(orderRepo, availabilityRepo)
	bookingHandler   = httpapi.NewBookingHandler(bookingService, LogInfof, LogErrorf)
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// hack
	availabilityRepo.Availability = append(
		availabilityRepo.Availability,
		// Radisson?
		booking.RoomAvailability{HotelID: "reddison", RoomID: "lux", Date: booking.NewDate(2024, 1, 1).ToTime(), Quota: 1},
		booking.RoomAvailability{HotelID: "reddison", RoomID: "lux", Date: booking.NewDate(2024, 1, 2).ToTime(), Quota: 1},
		booking.RoomAvailability{HotelID: "reddison", RoomID: "lux", Date: booking.NewDate(2024, 1, 3).ToTime(), Quota: 1},
		booking.RoomAvailability{HotelID: "reddison", RoomID: "lux", Date: booking.NewDate(2024, 1, 4).ToTime(), Quota: 1},
		booking.RoomAvailability{HotelID: "reddison", RoomID: "lux", Date: booking.NewDate(2024, 1, 5).ToTime(), Quota: 0},
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/orders", bookingHandler.CreateOrder)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	var eg errgroup.Group
	eg.Go(func() error {
		LogInfof("Server is listening on localhost:8080")
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			LogInfof("Server closed")
			return nil
		}
		return err
	})
	eg.Go(func() error {
		<-ctx.Done()
		LogInfof("Termination signal received. Server will shut down gracefully within 5 seconds.")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return srv.Shutdown(ctx)
	})

	if err := eg.Wait(); err != nil {
		log.Fatalf("Server stopped unexpectedly: %v", err)
	}
}

func LogErrorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func LogInfof(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}
