package domain

import (
	"testing"
)

// Тест 1 — создание отправления
func TestNewShipment_Success(t *testing.T) {
	s, err := NewShipment("REF-001", "Almaty", "Astana", "John", "TRUCK-1", 1000, 500)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if s.ID == "" {
		t.Error("expected ID to be set")
	}
	if s.Status != StatusPending {
		t.Errorf("expected status pending, got: %v", s.Status)
	}
	if s.ReferenceNumber != "REF-001" {
		t.Errorf("expected REF-001, got: %v", s.ReferenceNumber)
	}
}

// Тест 2 — создание без обязательных полей
func TestNewShipment_ValidationError(t *testing.T) {
	tests := []struct {
		name          string
		refNumber     string
		origin        string
		destination   string
		expectedError string
	}{
		{"empty ref", "", "Almaty", "Astana", "reference number is required"},
		{"empty origin", "REF-001", "", "Astana", "origin is required"},
		{"empty destination", "REF-001", "Almaty", "", "destination is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewShipment(tt.refNumber, tt.origin, tt.destination, "John", "TRUCK-1", 1000, 500)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if err.Error() != tt.expectedError {
				t.Errorf("expected '%v', got '%v'", tt.expectedError, err.Error())
			}
		})
	}
}

// Тест 3 — валидные переходы статусов
func TestShipment_ValidTransitions(t *testing.T) {
	tests := []struct {
		from Status
		to   Status
	}{
		{StatusPending, StatusPickedUp},
		{StatusPickedUp, StatusInTransit},
		{StatusInTransit, StatusDelivered},
		{StatusPending, StatusCancelled},
		{StatusPickedUp, StatusCancelled},
		{StatusInTransit, StatusCancelled},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			s, _ := NewShipment("REF-001", "Almaty", "Astana", "John", "TRUCK-1", 1000, 500)
			s.Status = tt.from

			event, err := s.Transition(tt.to)

			if err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if s.Status != tt.to {
				t.Errorf("expected status %v, got %v", tt.to, s.Status)
			}
			if event == nil {
				t.Error("expected event to be created")
			}
		})
	}
}

// Тест 4 — невалидные переходы статусов
func TestShipment_InvalidTransitions(t *testing.T) {
	tests := []struct {
		from Status
		to   Status
	}{
		{StatusPending, StatusInTransit},
		{StatusPending, StatusDelivered},
		{StatusDelivered, StatusPending},
		{StatusDelivered, StatusPickedUp},
		{StatusCancelled, StatusPending},
		{StatusCancelled, StatusPickedUp},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			s, _ := NewShipment("REF-001", "Almaty", "Astana", "John", "TRUCK-1", 1000, 500)
			s.Status = tt.from

			_, err := s.Transition(tt.to)

			if err == nil {
				t.Errorf("expected error for transition %v -> %v", tt.from, tt.to)
			}
			if s.Status != tt.from {
				t.Errorf("status should not change on invalid transition")
			}
		})
	}
}
