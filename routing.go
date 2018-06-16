package main

import (
	//"encoding/json"
	"github.com/husobee/vestigo"
	"net/http"
)

// 🌸 rauting 🌸 //

func RegisterRoutes(router *vestigo.Router) {
	router.Get("/api/test", testHandler)
	router.Post("/api/ebooks",  postEbooksHandler)
	router.Post("/api/ebooks/", postEbooksHandler)
	router.Post("/api/ebook/:id",  postEbookHandler)
	router.Post("/api/ebook/:id/", postEbookHandler)
}

func postEbookHandler(w http.ResponseWriter, r *http.Request) {
	errorResponse := checkPasswordRequest(r)
	if errorResponse != nil {
		writeErrorResponse(w, errorResponse)
		return
	}

    ebook := EbooksResponse{getEbook()}
	writeOKResponse(w, ebook)
}

func postEbooksHandler(w http.ResponseWriter, r *http.Request) {
	errorResponse := checkPasswordRequest(r)
	if errorResponse != nil {
		writeErrorResponse(w, errorResponse)
		return
	}

    ebooks := EbooksResponse{getEbooks()}
	writeOKResponse(w, ebooks)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	testEbook := Ebook{
        Author: "Kim Plong",
        Title: "Good Book",
        Rights: "Copyright 2018",
        Description: "Good book with good content and good cover and good characters and good plot",
        Filename: "goodbook.epub",
	}
    testEbooks := []Ebook{
        testEbook,
    }
    testResponse := EbooksResponse{
        Ebooks: testEbooks,
    }

	writeOKResponse(w, testResponse)
}

// ebook stuff //

type EbooksResponse struct {
	Ebooks []Ebook `json:"ebooks"`
}

type EbookResponse struct {
	Id         int64 `json:"id"`
	Downloaded int64 `json:"downloaded"`
}
