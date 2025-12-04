package auth

import (
	"runtime"
	"github.com/alexedwards/argon2id"
)



func HashPassword(password string) (string, error) {
	params := &argon2id.Params{Memory: 2, Iterations: 5, Parallelism: uint8(runtime.NumCPU()), SaltLength: 32, KeyLength: 64}
	return argon2id.CreateHash(password, params)
}


func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
