package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sikozonpc/notebase/data"
	t "github.com/sikozonpc/notebase/types"
)

type APIServer struct {
	addr  string
	store Storage
}

type EndpointHandler func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/highlight", makeHTTPHandler(s.handleHighlights))

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}

func (s *APIServer) handleHighlights(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetHighlights(w, r)
	}

	if r.Method == http.MethodPost {
		return s.handleCreateHighlight(w, r)
	}

	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *APIServer) handleGetHighlights(w http.ResponseWriter, r *http.Request) error {
	highlights, err := s.store.GetHighlights()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, highlights)

}

func (s *APIServer) handleCreateHighlight(w http.ResponseWriter, r *http.Request) error {
	req := new(t.CreateHighlightRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	highlight := data.NewHighlight(req.Text, req.Location, req.Note, req.UserId, req.BookId)

	if err := s.store.CreateHighlight(*highlight); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, highlight)

}
