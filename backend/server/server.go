package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

type Server struct {
	port   string
	dbConn *pgxpool.Pool
}

func NewServer(dbConnStr string, port string) (*Server, error) {
	conn, err := pgxpool.New(context.Background(), dbConnStr)
	if err != nil {
		return nil, err
	}

	return &Server{port, conn}, nil
}

func (s *Server) Run() {
	defer s.dbConn.Close()

	mux := http.NewServeMux()

	s.addRoutes(mux)

	handler := cors.Default().Handler(mux)

	log.Printf("Server starting on port %s\n", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.port), addMiddleware(handler)))
}

func (s *Server) addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /profile/{id}", s.getProfile)
	mux.HandleFunc("POST /profile", s.createProfile)
	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("POST /manufacturer", s.createManufacturer)
	mux.HandleFunc("GET /manufacturer/{id}", s.getManufacturer)
	mux.HandleFunc("POST /models", s.createModel)
	mux.HandleFunc("GET /models/{id}", s.getModel)
}
