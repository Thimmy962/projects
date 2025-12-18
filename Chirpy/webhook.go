package main

import (
	"Chirpy/internal/auth"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)



func (s *Server)webhook(w http.ResponseWriter, req *http.Request) {
	values := auth.GetAPIKey(req.Header)
	if values.Code != http.StatusOK {
		respondWithJSON(w, values.Code, values.Err)
		return
	}

	//compare the apikey sent in the header to the apikey in the .env file
	if Err := strings.Compare(s.apiKey, values.Key); Err != 0 {
		ProcessingError(w, http.StatusUnauthorized, nil)
		return
	}
	var buf bytes.Buffer
	buf.ReadFrom(req.Body)

	data := struct {
		Event string `json:"event"`
		Data struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		ProcessingError(w, http.StatusInternalServerError, err)
		return
	}
	if strings.Compare(data.Event, "user.upgraded") == 0 {
		if err := s.queries.UpdateRed(req.Context(), data.Data.UserId); err != nil {
			ProcessingError(w, http.StatusNotFound, err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}