package main

import (
	//"encoding/json"
    "fmt"
	"github.com/husobee/vestigo"
    "log"
	"net/http"
    "time"
)

// 🌸 rauting 🌸 //

func RegisterRoutes(router *vestigo.Router) {
	router.Get("/api/test", testHandler)
	router.Post("/api/ebooks",  postEbooksHandler)
	router.Post("/api/ebooks/", postEbooksHandler)
	router.Post("/api/ebook/:filename",  postEbookHandler)
	router.Post("/api/ebook/:filename/", postEbookHandler)
}

func postEbookHandler(w http.ResponseWriter, r *http.Request) {
	errorResponse := checkPasswordRequest(r)
	if errorResponse != nil {
        time.Sleep(3000 * time.Millisecond)
		writeErrorResponse(w, errorResponse)
		return
	}
	filename := vestigo.Param(r, "filename")
    ebook, err := getEbook(filename)
    if err != nil {
        log.Print(err)
        errorString := fmt.Sprintf("Error trying to open file %s", filename)
		errorResponse := newErrorResponse(500, errorString)
        writeErrorResponse(w, errorResponse)
        return
    }

    w.Write(ebook)

    attachmentString := "attachment; filename='" + filename + "'"
    w.Header().Set("Content-Disposition", attachmentString)
    w.Header().Set("Content-Type", "application/epub+zip")

	writeOKResponse(w, ebook)
}

func postEbooksHandler(w http.ResponseWriter, r *http.Request) {
	errorResponse := checkPasswordRequest(r)
	if errorResponse != nil {
        time.Sleep(3000 * time.Millisecond)
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
