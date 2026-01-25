package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)



func HashPassword(password string) (string, error) {
	params := &argon2id.Params{Memory: 2, Iterations: 5, Parallelism: uint8(runtime.NumCPU()), SaltLength: 32, KeyLength: 64}
	return argon2id.CreateHash(password, params)
}


func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

//Makes web tokens
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	regClaims := jwt.RegisteredClaims{Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject: userID.String(),
	}
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, regClaims)
	signed, err := tokens.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}


//Validates JWT
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.MapClaims{}


	// Parse the token using your secret
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	// Extract user ID (stored in the "sub" claim)
	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid subject claim")
	}

	// Convert string → UUID
	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// Extracts the userId from authentication token
func GetBearerToken(headers http.Header, secret string) (string, error) {
	// extract the authorization header
	rawHeader := headers.Get("Authorization")
	if rawHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// split on spacec
	parts := strings.Fields(rawHeader)
	if len(parts) < 2 {
		return "", fmt.Errorf("authorization header is malformed (expected: Bearer <token>)")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("authorization type must be Bearer")
	}

	tokenString := parts[1]
	if tokenString == "" {
		return "", fmt.Errorf("bearer token is empty")
	}


	id, err := ValidateJWT(tokenString, secret)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return id.String(), nil
}


// for extracting refresh token from the headers
func GetBearerRefreshToken(headers http.Header) (string, error) {
	// extract the authorization header
	rawHeader := headers.Get("Authorization")
	if rawHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// split on spacec
	parts := strings.Fields(rawHeader)
	if len(parts) < 2 {
		return "", fmt.Errorf("authorization header is malformed (expected: Bearer <token>)")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("authorization type must be Bearer")
	}

	tokenString := parts[1]
	if tokenString == "" {
		return "", fmt.Errorf("bearer token is empty")
	}
	return tokenString, nil
}

func MakeRefreshToken() (string, error) {
	buf := make([]byte, 32)
	rand.Read(buf)
	refresh_token := fmt.Sprintf("%x%s", buf, time.Now().String())
	hashed := fmt.Sprintf("%x", sha256.Sum256([]byte(refresh_token)))
	return hashed, nil
}

type Values struct{
			Key string
			Err error
			Code int
		}

func GetAPIKey(headers http.Header) Values {
	apiKey := headers.Get("Apikey")
	if apiKey == "" {
		return struct{
			Key string
			Err error
			Code int
		}{Key: "", Err: errors.New("apikey header is missing"), Code: http.StatusBadRequest}
	}

	arr := strings.Fields(apiKey)
	if len(arr) != 2 {
		return struct {
			Key string
			Err error
			Code int}{Key: "", Err: errors.New("apikey header is malformed"), Code: http.StatusBadRequest}
	}

	return struct {
			Key string
			Err error
			Code int}{Key: "", Err: nil, Code: http.StatusOK}
}