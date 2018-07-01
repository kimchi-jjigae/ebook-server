package main

import (
	"encoding/json"
	"net/http"
    "log"
)

// passwörd stuff //

type passwordRequest struct {
	Password string `json: "password"`
}

func checkPasswordRequest(r *http.Request, correctPassword string) *errorResponse {
	decoder := json.NewDecoder(r.Body)
	var pw passwordRequest
	err := decoder.Decode(&pw)
	if err != nil {
		return newErrorResponse(400, err.Error())
	}
	if pw.Password != correctPassword {
        log.Print("Invalid password attempt!")
		return newErrorResponse(401, "invalid password")
	}
    log.Print("Successful password input!")
	return nil
}
