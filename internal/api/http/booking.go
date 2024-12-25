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
	createOrderRequest := struct {
		HotelID   string    `json:"hotel_id"`
		RoomID    string    `json:"room_id"`
		UserEmail string    `json:"email"`
		From      time.Time `json:"from"`
		To        time.Time `json:"to"`

		// PromoCode           float64
		// Discount            float64
		// LoyaltyPointsearned int
		// Rooms               []RoomDetails
	}{}

	if err := json.NewDecoder(r.Body).Decode(&createOrderRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		h.errorLogger("failed to decode request", err)
		return
	}

	if err := h.service.CreateOrder(booking.Order(createOrderRequest)); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		h.errorLogger("failed to create order", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createOrderRequest); err != nil {
		h.errorLogger("failed to write response: %v", err)
		return
	}
	h.infoLogger("order created successfully: %v", createOrderRequest)
}
