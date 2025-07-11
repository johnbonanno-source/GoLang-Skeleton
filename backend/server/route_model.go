package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type NewModel struct {
	Name           string `json:"name"`
	ManufacturerID int    `json:"manufacturer_id"`
}

func (s *Server) createModel(w http.ResponseWriter, r *http.Request) {
	var model NewModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&model); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	var id int

	err := s.dbConn.QueryRow(context.Background(), "insert into models(name, manufacturer_id) values($1, $2) returning id", model.Name, model.ManufacturerID).Scan(&id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unknown error")
		return
	}

	respondWithJSON(w, http.StatusCreated, Model{id, model.Name, model.ManufacturerID})
}

type Model struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ManufacturerID int    `json:"manufacturer_id"`
}

func (s *Server) getModel(w http.ResponseWriter, r *http.Request) {
	reqId := r.PathValue("id")
	if reqId == "" {
		respondWithError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	var resp Model

	err := s.dbConn.QueryRow(context.Background(), "select id, name, manufacturer_id from models where id = $1", reqId).Scan(&resp.Id, &resp.Name, &resp.ManufacturerID)
	if err != nil {
		log.Println(err)
		switch err {
		case pgx.ErrNoRows:
			respondWithError(w, http.StatusBadRequest, "id does not exist")
		default:
			respondWithError(w, http.StatusInternalServerError, "server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}
