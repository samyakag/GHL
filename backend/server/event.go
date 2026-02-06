package server

import (
	"context"
	"fmt"
	"log"
	"sync"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	eventv1 "github.com/samyakxd/ghl/backend/gen/event/v1"
)

type EventServer struct {
	mu            sync.RWMutex
	events        map[string]*eventv1.Event
	registrations map[string][]*eventv1.Registration // keyed by event_id
}

func NewEventServer() *EventServer {
	return &EventServer{
		events:        make(map[string]*eventv1.Event),
		registrations: make(map[string][]*eventv1.Registration),
	}
}

func (s *EventServer) CreateEvent(
	ctx context.Context,
	req *connect.Request[eventv1.CreateEventRequest],
) (*connect.Response[eventv1.CreateEventResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event := &eventv1.Event{
		Id:              uuid.New().String(),
		Title:           req.Msg.Title,
		Description:     req.Msg.Description,
		Date:            req.Msg.Date,
		Capacity:        req.Msg.Capacity,
		RegisteredCount: 0,
	}
	s.events[event.Id] = event

	log.Printf("Created event: %s - %s (capacity: %d)", event.Id, event.Title, event.Capacity)
	return connect.NewResponse(&eventv1.CreateEventResponse{Event: event}), nil
}

func (s *EventServer) ListEvents(
	ctx context.Context,
	req *connect.Request[eventv1.ListEventsRequest],
) (*connect.Response[eventv1.ListEventsResponse], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*eventv1.Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}

	log.Printf("Listed %d events", len(events))
	return connect.NewResponse(&eventv1.ListEventsResponse{Events: events}), nil
}

func (s *EventServer) RegisterForEvent(
	ctx context.Context,
	req *connect.Request[eventv1.RegisterForEventRequest],
) (*connect.Response[eventv1.RegisterForEventResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, exists := s.events[req.Msg.EventId]
	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("event not found: %s", req.Msg.EventId))
	}

	if event.RegisteredCount >= event.Capacity {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("event is full: %s", req.Msg.EventId))
	}

	registration := &eventv1.Registration{
		Id:      uuid.New().String(),
		EventId: req.Msg.EventId,
		Name:    req.Msg.Name,
		Email:   req.Msg.Email,
	}

	s.registrations[req.Msg.EventId] = append(s.registrations[req.Msg.EventId], registration)
	event.RegisteredCount++

	log.Printf("Registered %s (%s) for event %s (%d/%d)", registration.Name, registration.Email, event.Title, event.RegisteredCount, event.Capacity)
	return connect.NewResponse(&eventv1.RegisterForEventResponse{Registration: registration}), nil
}

func (s *EventServer) GetEventRegistrations(
	ctx context.Context,
	req *connect.Request[eventv1.GetEventRegistrationsRequest],
) (*connect.Response[eventv1.GetEventRegistrationsResponse], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.events[req.Msg.EventId]; !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("event not found: %s", req.Msg.EventId))
	}

	registrations := s.registrations[req.Msg.EventId]
	if registrations == nil {
		registrations = []*eventv1.Registration{}
	}

	log.Printf("Listed %d registrations for event %s", len(registrations), req.Msg.EventId)
	return connect.NewResponse(&eventv1.GetEventRegistrationsResponse{Registrations: registrations}), nil
}
