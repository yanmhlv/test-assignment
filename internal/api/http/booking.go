package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"yanmhlv/test-assignment/internal/booking"
)

type logger interface {
	Error(string, ...any)
	Info(string, ...any)
}

type BookingHandler struct {
	service *booking.BookingService

	infoLogger  func(string, ...any)
	errorLogger func(string, ...any)
}

func NewBookingHandler(service *booking.BookingService, infoLog, errorLog func(string, ...any)) *BookingHandler {
	if infoLog == nil {
		infoLog = log.Default().Printf
	}

	if errorLog == nil {
		errorLog = log.Default().Printf
	}
	return &BookingHandler{service, infoLog, errorLog}
}

func (h *BookingHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	newOrder := struct {
		HotelID   string    `json:"hotel_id"`
		RoomID    string    `json:"room_id"`
		UserEmail string    `json:"email"`
		From      time.Time `json:"from"`
		To        time.Time `json:"to"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		h.errorLogger("failed to decode request", err)
		return
	}

	if err := h.service.CreateOrder(booking.Order(newOrder)); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		h.errorLogger("failed to create order", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
	h.infoLogger("order created successfully: %v", newOrder)
}
