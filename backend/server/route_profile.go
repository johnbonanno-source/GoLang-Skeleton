package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type NewProfile struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *Server) createProfile(w http.ResponseWriter, r *http.Request) {
	var profile NewProfile
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&profile); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	var id int

	err := s.dbConn.QueryRow(context.Background(), "insert into profiles(email, name) values($1, $2) returning id", profile.Email, profile.Name).Scan(&id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unknown error")
		return
	}

	respondWithJSON(w, http.StatusCreated, Profile{id, profile.Email, profile.Name})
}

type Profile struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *Server) getProfile(w http.ResponseWriter, r *http.Request) {
	authData := getAuthInfo(r)
	log.Printf("authId = %d", authData.userID)
	reqId := r.PathValue("id")
	if reqId == "" {
		respondWithError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	var resp Profile

	err := s.dbConn.QueryRow(context.Background(), "select id, email, name from profiles where id = $1", reqId).Scan(&resp.Id, &resp.Email, &resp.Name)
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
