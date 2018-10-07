package main

import (
    "log"
)

// passwörd stuff //

func checkPasswordRequest(attemptedPassword string, correctPassword string) *errorResponse {
	if attemptedPassword != correctPassword {
        log.Print("Invalid password attempt!")
		return newErrorResponse(401, "invalid password")
	}
    log.Print("Successful password input!")
	return nil
}
