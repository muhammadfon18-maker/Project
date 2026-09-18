package helper

import (
	"errors"
	"net/http"
)

var UserNotFound = errors.New("The user is not found")

var CarNotFound = errors.New("The car is not found")

var Invalid = errors.New("Invalid operation")

func Error(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
