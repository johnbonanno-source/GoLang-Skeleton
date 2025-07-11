package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Manufacturer struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (s *Server) createManufacturer(w http.ResponseWriter, r *http.Request) {
	var manufacturer Manufacturer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&manufacturer); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	var id int

	err := s.dbConn.QueryRow(context.Background(), "insert into manufacturers(name) values($1) returning id", manufacturer.Name).Scan(&id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unknown error")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"id":   id,
		"name": manufacturer.Name,
	})
}

func (s *Server) getManufacturer(w http.ResponseWriter, r *http.Request) {
	reqId := r.PathValue("id")
	if reqId == "" {
		respondWithError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	var manufacturer Manufacturer

	err := s.dbConn.QueryRow(context.Background(), "select id, name from manufacturers where id = $1", reqId).Scan(&manufacturer.Id, &manufacturer.Name)
	if err != nil {
		log.Println(err)
		switch err {
		case pgx.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "manufacturer not found")
		default:
			respondWithError(w, http.StatusInternalServerError, "server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, manufacturer)
}
