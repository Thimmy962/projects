package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func (s *Server)createChirp(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			ProcessingError(w, http.StatusBadRequest, err)
			return
		}
	}
	// get the user id from the auth token
	userID, err := auth.GetBearerToken(req.Header, s.secret)
	if err != nil {
		ProcessingError(w, 401, err)
		return
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

	db_chirp, err := s.queries.CreateChirp(req.Context(), database.CreateChirpParams{Body: chirp.Body, UserID: userID})
	if err != nil {
		log.Println(err.Error())
		ProcessingError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w,http.StatusCreated, db_chirp)
}


func (s *Server)listChirps(w http.ResponseWriter, req *http.Request) {
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


func (s *Server)getChirp(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	data, err := s.queries.GetChirp(req.Context(), id)
	if err != nil {
		ProcessingError(w, http.StatusNotFound, err)
	}
	respondWithJSON(w, 200, data)
}



func (s *Server)deleteChirp(w http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("id")
	userId, err := auth.GetBearerToken(req.Header, s.secret); if err != nil {
		ProcessingError(w, http.StatusBadRequest, err)
		return
	}

	//returns the chirp is and the userId of the user that created the chirp
	chirpRow, err := s.queries.GetChirpByChirpID(req.Context(), chirpId)
	if err != nil {
			ProcessingError(w, http.StatusNotFound, errors.New("not found"))
			return
	}

	// compares the id of the user that created the chirp and the Id of the user extracted from the access token
	if strings.Compare(userId, chirpRow.UserID) != 0 {
		ProcessingError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	if s.queries.DeleteChirpById(req.Context(), chirpId) != nil {
		ProcessingError(w, http.StatusInternalServerError, errors.New("an error occured"))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}