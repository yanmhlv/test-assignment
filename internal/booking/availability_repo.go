package booking

import (
	"errors"
	"sync"
	"time"
)

var _ AvailabilityRepository = (*InMemoryAvailabilityRepository)(nil)

type AvailabilityRepository interface {
	GetAvailability(hotelID, roomID string, date time.Time) (*RoomAvailability, error)
	UpdateAvailability(avail RoomAvailability) error
}

type InMemoryAvailabilityRepository struct {
	lock         sync.RWMutex
	Availability []RoomAvailability // TODO: use map[hotelID][]rooms
}

func NewInMemoryAvailabilityRepository() *InMemoryAvailabilityRepository {
	return &InMemoryAvailabilityRepository{}
}

func (r *InMemoryAvailabilityRepository) GetAvailability(hotelID, roomID string, date time.Time) (*RoomAvailability, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for i, avail := range r.Availability {
		if avail.HotelID == hotelID && avail.RoomID == roomID && avail.Date.Equal(date) {
			return &r.Availability[i], nil
		}
	}
	return nil, nil
}

func (r *InMemoryAvailabilityRepository) UpdateAvailability(avail RoomAvailability) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, a := range r.Availability {
		if a.HotelID == avail.HotelID && a.RoomID == avail.RoomID && a.Date.Equal(avail.Date) {
			r.Availability[i] = avail
			return nil
		}
	}
	return errors.New("room is not found")
}
