package booking

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type BookingService struct {
	orderRepo        OrderRepository
	availabilityRepo AvailabilityRepository

	lock sync.Mutex
}

func NewBookingService(orderRepo OrderRepository, availabilityRepo AvailabilityRepository) *BookingService {
	return &BookingService{orderRepo, availabilityRepo, sync.Mutex{}}
}

func (s *BookingService) CreateOrder(order Order) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	daysToBook := daysBetween(order.From, order.To)
	if len(daysToBook) == 0 {
		return errors.New("invalid date range")
	}

	unavailableDays := make(map[time.Time]struct{})
	for _, day := range daysToBook {
		availability, err := s.availabilityRepo.GetAvailability(order.HotelID, order.RoomID, day)
		if err != nil || availability == nil || availability.Quota < 1 {
			unavailableDays[day] = struct{}{}
			continue
		}

		availability.Quota--
		if err := s.availabilityRepo.UpdateAvailability(*availability); err != nil {
			return fmt.Errorf("failed to update availability: %w", err)
		}
	}

	if len(unavailableDays) > 0 {
		return fmt.Errorf("room is unavailable for these dates: %v", unavailableDays)
	}

	if err := s.orderRepo.Create(order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func daysBetween(from, to time.Time) []time.Time {
	if from.After(to) {
		return nil
	}
	var days []time.Time
	for d := toDay(from); !d.After(toDay(to)); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}
	return days
}

func toDay(ts time.Time) time.Time {
	return time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.UTC)
}
