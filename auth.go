package main

import (
	"encoding/json"
	"net/http"
)

// passwörd stuff //

type passwordRequest struct {
	Password string `json: "password"`
}

func checkPasswordRequest(r *http.Request) *errorResponse {
	decoder := json.NewDecoder(r.Body)
	var pw passwordRequest
	err := decoder.Decode(&pw)
	if err != nil {
		return newErrorResponse(400, err.Error())
	}
	if !checkPassword(pw.Password) {
		return newErrorResponse(401, "invalid password")
	}
	return nil
}

func checkPassword(password string) bool {
	return password == "shittypassword"
}