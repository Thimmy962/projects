package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
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
	user :=  struct {
		Email string `json:"email"`
		PWord string `json:"password"`
	}{}
	err := json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		ProcessingError(w, 500, err)
	}

	if len(user.PWord) < 8 {
		ProcessingError(w, 400, fmt.Errorf("password length is less than 8\n"))
		return
	}

	hash, err := auth.HashPassword(user.PWord); if err != nil {
		ProcessingError(w, 500, err)
		return
	}
	user.PWord=hash
	
	params := database.CreateUserParams{Email: sql.NullString{user.Email, true}, HashedPassword: user.PWord}
	db_user, dbErr := s.queries.CreateUser(req.Context(), params)
	if dbErr != nil {
		ProcessingError(w, http.StatusBadRequest, dbErr)
		return
	}
	respondWithJSON(w, http.StatusCreated, db_user)
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
	respondWithJSON(w, http.StatusNoContent, dbErr)
}



func (s *Server)GetUser(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	buf.ReadFrom(req.Body)
	// an anonymous struct
	user :=  struct {
		Email string `json:"email"`
		PWord string `json:"password"`
	}{}

	err := json.Unmarshal(buf.Bytes(), &user); if err != nil {
		ProcessingError(w, 500, err)
		return
	}

	// this get the hashed passwd
	data, err := s.queries.GetUserPassword(req.Context(), sql.NullString{user.Email, true})
	if err != nil {
		ProcessingError(w, 404, err)
		return
	}

	// compare a string to hash
	// could have hitten the DB once but the task requires presenting data without hashed password
	// so I had to choose between more process of data or hitting the DB twice, so I choose the latter
	correct_password, err := auth.CheckPasswordHash(user.PWord, data)
	if err != nil {
		ProcessingError(w, 500, err)
		return
	}

	if !correct_password {
		ProcessingError(w, 401, fmt.Errorf("email or password incorrect"))
		return
	}
	userRowData := database.GetUserParams{Email: user.Email, HashedPassword: data}
	userRow, err := s.queries.GetUser(req.Context(), userRowData)
	if err != nil {
		ProcessingError(w, 500, err)
		return
	}
	respondWithJSON(w, 200, userRow)
}