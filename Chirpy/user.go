package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/google/uuid"
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
		ProcessingError(w, 400, fmt.Errorf("password length is less than 8"))
		return
	}

	hash, err := auth.HashPassword(user.PWord); if err != nil {
		ProcessingError(w, 500, err)
		return
	}
	user.PWord=hash
	
	params := database.CreateUserParams{Email: user.Email, HashedPassword: user.PWord}
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

func (s *Server)getUserToken(w http.ResponseWriter, req *http.Request) {
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
	data, err := s.queries.GetUserPassword(req.Context(), user.Email)
	if err != nil {
		ProcessingError(w, 404, errors.New("username or password is wrong"))
		return
	}

	// compare a string to hash
	// could have hitten the DB once but the task requires presenting data without hashed password
	// so I had to choose between more processing of data or hitting the DB twice, so I choose the later
	correct_password, err := auth.CheckPasswordHash(user.PWord, data)
	if err != nil {
		ProcessingError(w, 500, err)
		return
	}

	if !correct_password {
		ProcessingError(w, 401, errors.New("email or password incorrect"))
		return
	}
	userRowData := database.GetUserParams{Email: user.Email, HashedPassword: data}
	userRow, err := s.queries.GetUser(req.Context(), userRowData)
	if err != nil {
		ProcessingError(w, 500, err)
		return
	}
	
	token, err := auth.MakeJWT(uuid.MustParse(userRow.ID), s.secret, 3600 * time.Second)
	if err != nil {
		ProcessingError(w, 500, err)
		return
	}

	refresh_token, _ := auth.MakeRefreshToken();
	params := database.CreateRefreshTokenParams{Token: refresh_token,
	UserID: userRow.ID}
	err = s.queries.CreateRefreshToken(req.Context(), params)
	if err != nil {
		ProcessingError(w, 500, err)
		return
	}

	det := struct {
		ID        string `json:"id`
		CreatedAt time.Time `json:"created"`
		UpdatedAt time.Time `json:"updated"`
		Email     string	`json:"email"`
		Token	  string `json:"token"`
		Refresh_Token string `json:"refresh_token"`
	}{ID: userRow.ID,
		CreatedAt: userRow.CreatedAt,
		UpdatedAt: userRow.UpdatedAt,
		Email: user.Email,
		Token: token,
		Refresh_Token: refresh_token,
	}

	respondWithJSON(w, 200, det)
}


func (s *Server)getUserDet(w http.ResponseWriter, req *http.Request) {
	id, err := auth.GetBearerToken(req.Header, s.secret); if err != nil {
		ProcessingError(w, 400, err)
		return
	}

	user, err := s.queries.GetUserByID(req.Context(), id)
	if err != nil {
		ProcessingError(w, 500, err)
		return
	}
	det := struct {
		Email, Id, Token string
	}{Email: user.Email, Id: user.ID, Token: id}
	respondWithJSON(w, 200, det)
}


func (s *Server)editUserDetail(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetBearerToken(r.Header, s.secret); if err != nil {
		ProcessingError(w, http.StatusUnauthorized, err)
		return
	}
	var buf bytes.Buffer
	buf.ReadFrom(r.Body)
	// an anonymous struct
	user :=  struct {
		Email string `json:"email"`
		Pword string `json:"password"`
	}{}
	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		ProcessingError(w, 500, err)
	}

	if len(user.Pword) < 8 {
		ProcessingError(w, 400, fmt.Errorf("password length is less than 8"))
		return
	}

	
	if len(user.Pword) < 8 {
		ProcessingError(w, http.StatusBadRequest, errors.New("password less than 8 digits"))
		return
	}

	hash, err := auth.HashPassword(user.Pword); if err != nil {
		ProcessingError(w, http.StatusInternalServerError, err)
		return
	}

	user.Pword = hash
	err = s.queries.ChangeDetail(r.Context(), database.ChangeDetailParams{Email: user.Email,
		HashedPassword: user.Pword,
		ID: userId}); if err != nil {
			ProcessingError(w, http.StatusBadRequest, err)
			return
		}
	
	row, err := s.queries.GetUserByID(r.Context(), userId); if err != nil {
		ProcessingError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, row)
}