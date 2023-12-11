package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"
	"golang.org/x/crypto/bcrypt"
)

type Cookie struct {
	Name     string
	Value    string
	MaxAge   *time.Time
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	SameSite int
}

// CreateToken is function to create a cookie
func CreateToken(id uint, username string) (*Cookie, error) {
	value := map[string]string{
		"id":       fmt.Sprintf("%d", id),
		"username": username,
	}
	// hashedValue with secret
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// max age is 7 days from now
	maxAgeDate := carbon.Now().AddDays(7)

	return &Cookie{
		Name:     "arthur-cookie",
		Value:    string(hashedValue),
		MaxAge:   &maxAgeDate,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: 0,
	}, nil
}

func CheckCookie(cookie *Cookie, id uint, username string) (bool, error) {
	value := map[string]string{
		"id":       fmt.Sprintf("%d", id),
		"username": username,
	}

	// check if cookie is valid
	err := bcrypt.CompareHashAndPassword([]byte(cookie.Value), []byte(value))
	if err != nil {
		return false, errors.New("invalid cookie")
	}
	return true, nil
}
