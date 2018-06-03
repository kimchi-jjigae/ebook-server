package main

import (
	//"encoding/json"
    "crypto/md5"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
)

type Ebook struct {
	Author      string `json:"author"`
	Title       string `json:"title"`
	Rights      string `json:"rights"`
	Description string `json:"description"`
	Filename    string `json:"filename"`
}

var searchDirs []string = []string{
    "/home/kim/dcc",
}

var storageDir string = "/home/kim/bookz"

func getEbooks() (ebooks []Ebook) {
    checkForNewBooks()
    ebookPaths := getEbookPaths([]string{storageDir})
	for _, ebookPath := range ebookPaths {
        //grab author, title, rights, description, filename (not whole path!)
        log.Print(ebookPath)
    }
    return
}

func checkForNewBooks() {
    ebookSearchPaths := getEbookPaths(searchDirs)
    ebookStoragePaths := getEbookPaths([]string{storageDir})
    hashedSearchEbooks := hashEbooks(ebookSearchPaths)
    hashedStorageEbooks := hashEbooks(ebookStoragePaths)
    // find hashes unique in the search path only
    searchEbooksUnique := hashedPathsDifference(hashedSearchEbooks, hashedStorageEbooks)
    copyEbooks(searchEbooksUnique, storageDir)
}

func hashEbooks(ebookPaths []string) (hashedPaths map[string]string) {
    // hash all the ebooks at the given paths
    hashedPaths = make(map[string]string)
	for _, path := range ebookPaths {
        ebookData, err := os.Open(path)
        if err != nil {
            log.Fatal(err)
        }
        hash := md5.New()
        if _, err := io.Copy(hash, ebookData); err != nil {
            log.Fatal(err)
        }
        hashString := fmt.Sprintf("%x", hash.Sum(nil))
        hashedPaths[hashString] = path
    }
    return
}

func hashedPathsDifference(hashedSearchEbooks map[string]string, hashedStorageEbooks map[string]string) (difference map[string]string) {
    difference = make(map[string]string)
    return
}

func copyEbooks(searchEbooksUnique map[string]string, storageDir string) {
}

func generateNewFilename() {
    // maybe like sort filenames and get last one, maybe should be b00032.epub (prefix + number) and increase the number or something
}

func getEbookPaths(dirs []string) (ebookPaths []string) {
    // grab all the .epub paths in the storage dir
	for _, dir := range dirs {
        files, err := ioutil.ReadDir(dir)
        if err != nil {
            log.Fatal(err)
        }
        for _, file := range files {
            filename := file.Name()
            // would be better to check type in another way but this will do for now
            if filename[len(filename)-4:] == ".epub" {
                ebookPaths = append(ebookPaths, dir + filename)
            }
        }
    }
    return
}
