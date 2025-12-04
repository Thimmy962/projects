package main

import (
	"Chirpy/internal/database"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (s *Server)CreateChirp(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			ProcessingError(w, http.StatusBadRequest, err)
			return
		}
	}
	chirp := Chirp{}
	unMashallError := json.Unmarshal(buf.Bytes(), &chirp); if unMashallError != nil {
		ProcessingError(w, http.StatusInternalServerError, unMashallError)
		return
	}

	switch chirp.validBody() {
	case -1:
		ProcessingError(w, http.StatusBadRequest, fmt.Errorf("length of the body is below lower limit(1)"))
		return
	case 1:
		ProcessingError(w, http.StatusBadRequest, fmt.Errorf("length of body is above upper limit(140)"))
		return
	}

	db_chirp, err := s.queries.CreateChirp(req.Context(), database.CreateChirpParams{Body: chirp.Body, UserID:chirp.UserId})
	if err != nil {
		log.Println(err.Error())
		ProcessingError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w,http.StatusCreated, db_chirp)
}


func (s *Server)ListChirps(w http.ResponseWriter, req *http.Request) {
	db_chirps, err := s.queries.ListChirps(req.Context())
	if err != nil {
		log.Println(err.Error())
		ProcessingError(w, http.StatusBadRequest, err)
		return
	}
	
	for _, chirp := range db_chirps {
		chirp.Body = profaneFUnc(chirp.Body)
	}

	respondWithJSON(w, http.StatusOK, db_chirps)
}


func (s *Server)GetChirp(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	data, err := s.queries.GetChirp(req.Context(), id)
	if err != nil {
		ProcessingError(w, http.StatusNotFound, err)
	}
	respondWithJSON(w, 200, data)
}