package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	pass := []byte(password)
	hash, e := bcrypt.GenerateFromPassword(pass, 10)
	if e != nil {
		return nil, e
	}

	return hash, nil
}

func CompareHash(password, hash string) bool {
	pass := []byte(password)
	hsh := []byte(hash)

	if e := bcrypt.CompareHashAndPassword(hsh, pass); e != nil {
		return false
	}

	return true
}
