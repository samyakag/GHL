package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/samyakxd/ghl/backend/gen/event/v1/eventv1connect"
	"github.com/samyakxd/ghl/backend/gen/todo/v1/todov1connect"
	"github.com/samyakxd/ghl/backend/server"
)

func main() {
	mux := http.NewServeMux()

	// Register TodoService
	todoPath, todoHandler := todov1connect.NewTodoServiceHandler(server.NewTodoServer())
	mux.Handle(todoPath, todoHandler)

	// Register EventService
	eventPath, eventHandler := eventv1connect.NewEventServiceHandler(server.NewEventServer())
	mux.Handle(eventPath, eventHandler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Connect-Protocol-Version"},
		ExposedHeaders:   []string{"Content-Type", "Connect-Protocol-Version"},
		AllowCredentials: true,
	})

	addr := ":8080"
	log.Printf("Starting server on %s", addr)
	log.Printf("TodoService available at %s%s", addr, todoPath)
	log.Printf("EventService available at %s%s", addr, eventPath)

	err := http.ListenAndServe(
		addr,
		h2c.NewHandler(corsHandler.Handler(mux), &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
