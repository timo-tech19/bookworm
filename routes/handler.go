package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	Router *mux.Router
	Server *http.Server
}

// Creates a configured http handler
func NewHandler() *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	h.mapRoutes()

	h.Server = &http.Server{
		Addr:         "localhost:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      h.Router,
	}

	return h
}

// Maps api routes to appropriate handler functions
func (h *Handler) mapRoutes() {
	// This route checks if the server is up and running
	h.Router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is still alive")
	})
}

// Starts up server.
// Server will gracefully shutdown on exit (CTRL + C)
func (h *Handler) Serve() {
	var wait time.Duration
	fmt.Println("Starting server on localhost:8080...")
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until receive exit signal
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	h.Server.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
