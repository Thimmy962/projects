package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func (s *Server)createUser(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	buf.ReadFrom(req.Body)
	// an anonymous struct
	email :=  struct {
		Email string `json:"email"`
	}{}
	json.Unmarshal(buf.Bytes(), &email)
	db_user, dbErr := s.queries.CreateUser(req.Context(), sql.NullString{email.Email, true})
	if dbErr != nil {
		ProcessingError(w, http.StatusBadRequest, dbErr)
		return
	}
	respondWithJSON(w, http.StatusCreated, buf, db_user)
}

func (s *Server)deleteUsers(w http.ResponseWriter, req *http.Request) {
	platform := os.Getenv("PLATFORM")
	if strings.Compare(platform, "dev") != 0 {
		ProcessingError(w, http.StatusForbidden, fmt.Errorf("you can only call this api in development mode"))
		return
	}
	dbErr := s.queries.DeleteUsers(req.Context())
	if dbErr != nil {
		ProcessingError(w, http.StatusBadRequest, dbErr)
	}
	respondWithJSONWithoutBuffer(w, http.StatusNoContent, dbErr)
}