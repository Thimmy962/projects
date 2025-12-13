package main

import (
	"Chirpy/internal/auth"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

/* 
 * used to get new refresh and access token 
 * once a refresh toekn is used, it becomes blacklisted and unusable
*/
func (s *Server)refresh(w http.ResponseWriter, r *http.Request) {
	// get refresh token
	token, err := auth.GetBearerRefreshToken(r.Header); if err != nil {
		ProcessingError(w, http.StatusBadRequest, err)
		return
	}

	// get user id from refresh token
	id, err := s.queries.GetUserIDFromRefreshToken(r.Context(), token); if err != nil {
		ProcessingError(w, http.StatusInternalServerError, fmt.Errorf("token expired"))
		return
	}

	// create a new access token with the user id and secrect in .env file
	jwtToken, err := auth.MakeJWT(uuid.MustParse(id), s.secret, 3600); if err != nil {
		ProcessingError(w, http.StatusInternalServerError, err)
		return
	}

	// delete the current refresh token
	s.queries.DeleteRefreshToken(r.Context(), token)
	new_token, err := auth.MakeRefreshToken(); if err != nil {
		ProcessingError(w, http.StatusInternalServerError, err)
		return
	}

	tokens := struct {
		Access_Token string `json:"access_token"`
		Refresh_Token string `json:"refresh_token"`
	}{Access_Token: jwtToken, Refresh_Token: new_token}

	respondWithJSON(w, http.StatusOK, tokens)
}