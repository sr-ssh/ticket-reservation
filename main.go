package ca2

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	ID               string
	Name             string
	Date             time.Time
	TotalTickets     int
	AvailableTickets int
}

type Ticket struct {
	ID      string
	EventID string
}

type TicketService struct {
	events  sync.Map
	tickets sync.Map
	mu      sync.Mutex
}

func (ts *TicketService) CreateEvent(name string, date time.Time, totalTickets int) (*Event, error) {
	event := &Event{
		ID:               generateUUID(),
		Name:             name,
		Date:             date,
		TotalTickets:     totalTickets,
		AvailableTickets: totalTickets,
	}

	ts.events.Store(event.ID, event)
	return event, nil
}

func (ts *TicketService) ListEvents() []*Event {
	var events []*Event

	ts.events.Range(func(key, value interface{}) bool {
		event := value.(*Event)
		events = append(events, event)
		return true
	})

	return events
}

func (ts *TicketService) BookTickets(eventID string, numTickets int) ([]string, error) {
	// Implement concurrency control here (Step 3)
	ts.mu.Lock() // Lock the mutex to ensure atomicity
	defer ts.mu.Unlock()

	// Load the event from the store
	event, ok := ts.events.Load(eventID)
	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	// Check if enough tickets are available
	eventValue := event.(*Event)
	if eventValue.AvailableTickets < numTickets {
		return nil, fmt.Errorf("not enough tickets available")
	}

	var ticketIDs []string
	for i := 0; i < numTickets; i++ {
		// Generate a ticket ID (implementation missing)
		ticketID := generateUUID() // Likely replaced with actual ticketID generation logic
		ticketIDs = append(ticketIDs, ticketID)

		// Store the ticket in a separate data structure (placeholder)
		ticket := &Ticket{
			ID:      ticketID,
			EventID: eventID,
		}
		ts.tickets.Store(ticketID, ticket)
	}

	// Update the available tickets and store the event
	eventValue.AvailableTickets -= numTickets
	ts.events.Store(eventID, event)

	return ticketIDs, nil
}
