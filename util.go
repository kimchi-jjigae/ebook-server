package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"
)

func argToInt(arg string) (int, *errorResponse) {
	argInt, err := strconv.Atoi(arg)

	if err != nil {
		return argInt, newErrorResponse(400, "arg '"+arg+"' not an int")
	}

	return argInt, nil
}

func argToInt64(arg string) (int64, *errorResponse) {
	n, err := argToInt(arg)
	return int64(n), err
}

type errorResponse struct {
	code          int
	errorString   string
	MessageString string `json:"error"`
}

func (e *errorResponse) Error() string {
	return e.errorString
}

func newErrorResponse(code int, errorString string) (response *errorResponse) {
	response = &errorResponse{}
	response.code = code
	response.MessageString = errorString

	if response.code >= 500 && response.code < 600 {
		response.MessageString = strconv.Itoa(code) + " - internal server error"
		response.errorString = errorString
	}

	return
}

func writeOKResponse(w http.ResponseWriter, responseData interface{}) {
	w.WriteHeader(200)

	response, err := json.Marshal(responseData)

	if err != nil {
		log.Print("marshal error: " + err.Error())
		writeErrorResponse(w, newErrorResponse(500, "could not marshal object of type '"+reflect.TypeOf(responseData).Name()+"'"))
		return
	}

	w.Write(response)
}

func writeErrorResponse(w http.ResponseWriter, response *errorResponse) {
	w.WriteHeader(response.code)

	if len(response.errorString) != 0 {
		log.Print("message: " + response.MessageString)
		log.Print("error: " + strconv.Itoa(response.code) + " - " + response.errorString)
		debug.PrintStack()
	}

	bytes, err := json.Marshal(response)

	if err != nil {
		log.Panic(err)
	}

	w.Write(bytes)
}


