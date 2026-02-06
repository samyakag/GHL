package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	todov1 "github.com/samyakxd/ghl/backend/gen/todo/v1"
)

type TodoServer struct {
	mu    sync.RWMutex
	todos map[string]*todov1.Todo
}

func NewTodoServer() *TodoServer {
	return &TodoServer{
		todos: make(map[string]*todov1.Todo),
	}
}

func (s *TodoServer) CreateTodo(
	ctx context.Context,
	req *connect.Request[todov1.CreateTodoRequest],
) (*connect.Response[todov1.CreateTodoResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := &todov1.Todo{
		Id:        uuid.New().String(),
		Title:     req.Msg.Title,
		Completed: false,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	s.todos[todo.Id] = todo

	log.Printf("Created todo: %s - %s", todo.Id, todo.Title)
	return connect.NewResponse(&todov1.CreateTodoResponse{Todo: todo}), nil
}

func (s *TodoServer) ListTodos(
	ctx context.Context,
	req *connect.Request[todov1.ListTodosRequest],
) (*connect.Response[todov1.ListTodosResponse], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]*todov1.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}

	log.Printf("Listed %d todos", len(todos))
	return connect.NewResponse(&todov1.ListTodosResponse{Todos: todos}), nil
}

func (s *TodoServer) UpdateTodo(
	ctx context.Context,
	req *connect.Request[todov1.UpdateTodoRequest],
) (*connect.Response[todov1.UpdateTodoResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, exists := s.todos[req.Msg.Id]
	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("todo not found: %s", req.Msg.Id))
	}

	todo.Title = req.Msg.Title
	todo.Completed = req.Msg.Completed

	log.Printf("Updated todo: %s - %s (completed: %v)", todo.Id, todo.Title, todo.Completed)
	return connect.NewResponse(&todov1.UpdateTodoResponse{Todo: todo}), nil
}

func (s *TodoServer) DeleteTodo(
	ctx context.Context,
	req *connect.Request[todov1.DeleteTodoRequest],
) (*connect.Response[todov1.DeleteTodoResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.todos[req.Msg.Id]; !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("todo not found: %s", req.Msg.Id))
	}

	delete(s.todos, req.Msg.Id)
	log.Printf("Deleted todo: %s", req.Msg.Id)
	return connect.NewResponse(&todov1.DeleteTodoResponse{}), nil
}
