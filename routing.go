package main

import (
    "fmt"
	"github.com/husobee/vestigo"
    "log"
	"net/http"
    "time"
)

// 🌸 rauting 🌸 //

type EbooksRouter struct {
    config *Config
    fileServer *EbooksFileServer
}

func newEbooksRouter (config *Config) (newRouter *EbooksRouter) {
    newRouter = &EbooksRouter{}
    newRouter.config = config
    newRouter.fileServer = newEbooksFileServer(config)
    return
}

func (router *EbooksRouter) RegisterRoutes(vRouter *vestigo.Router) {
	vRouter.Get("/api/test", router.testHandler)
	vRouter.Post("/api/ebooks",  router.postEbooksHandler)
	vRouter.Post("/api/ebooks/", router.postEbooksHandler)
	vRouter.Post("/api/ebook/:filename",  router.postEbookHandler)
	vRouter.Post("/api/ebook/:filename/", router.postEbookHandler)
}

func (router *EbooksRouter) postEbookHandler(w http.ResponseWriter, r *http.Request) {
	errorResponse := checkPasswordRequest(r, router.config.Password)
	if errorResponse != nil {
        time.Sleep(3000 * time.Millisecond)
		writeErrorResponse(w, errorResponse)
		return
	}
	filename := vestigo.Param(r, "filename")
    ebook, err := router.fileServer.GetEbook(filename)
    if err != nil {
        log.Print(err)
        errorString := fmt.Sprintf("Error trying to open file %s", filename)
		errorResponse := newErrorResponse(500, errorString)
        writeErrorResponse(w, errorResponse)
        return
    }

    attachmentString := "attachment; filename='" + filename + "'"
    w.Header().Set("Content-Disposition", attachmentString)
    w.Header().Set("Content-Type", "application/epub+zip")

    w.Write(ebook)

	writeOKResponse(w, ebook)
}

func (router *EbooksRouter) postEbooksHandler(w http.ResponseWriter, r *http.Request) {
	errorResponse := checkPasswordRequest(r, router.config.Password)
	if errorResponse != nil {
        time.Sleep(3000 * time.Millisecond)
		writeErrorResponse(w, errorResponse)
		return
	}

    ebooks := EbooksResponse{router.fileServer.GetEbooks()}
	writeOKResponse(w, ebooks)
}

func (router *EbooksRouter) testHandler(w http.ResponseWriter, r *http.Request) {
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
