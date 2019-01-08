package main

import (
    "crypto/md5"
    "fmt"
    "github.com/n3integration/epub"
    "io"
    "io/ioutil"
    "log"
    "math/rand"
    "path"
    "regexp"
    "strconv"
    "strings"
    "time"
    "os"
)

type Ebook struct {
	Author      string `json:"author"`
	Title       string `json:"title"`
	Rights      string `json:"rights"`
	Description string `json:"description"`
	Filename    string `json:"filename"`
}

type EbooksFileServer struct {
    config *Config
}

func newEbooksFileServer (config *Config) (newServer *EbooksFileServer) {
    newServer = &EbooksFileServer{}
    newServer.config = config
    return
}

func (server *EbooksFileServer) RemoveTempEbook(filename string) (err error) {
    timeout_ns := time.Duration(60000000000)
    log.Printf("Timing out for %s seconds to allow file download", timeout_ns / 1000000000)
    time.Sleep(timeout_ns)
    filepath := path.Join(server.config.TempDir, filename)
    log.Printf("Will now delete file %s", filepath)
    err = os.Remove(filepath)
    if err != nil {
        log.Printf("Could not delete file at %s!", filepath)
        log.Print(err)
        return err
    }
    log.Print("Done!")
    return nil
}

// hallo ni runer som sitter härutanför
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomFilename() (filename string) {
    rand.Seed(time.Now().UnixNano())

    randomLetters := make([]rune, 7)
    for i := range randomLetters {
        randomLetters[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(randomLetters) + ".epub"
}

func (server *EbooksFileServer) copyToTemporaryFilepath(filename string) (newFilename string, err error) {
    newFilename = generateRandomFilename()
    destFilepath := path.Join(server.config.TempDir, newFilename)
    srcFilepath := path.Join(server.config.StorageDir, filename)
    err = server.Copy(srcFilepath, destFilepath)
    if err != nil {
        log.Printf("Could not copy book from '%s' to temp filepath '%s'!", srcFilepath, destFilepath)
        return
    }
    return
}

func (server *EbooksFileServer) GetTempEbookFilepath(filename string) (filepath string, err error) {
    filepath, err = server.copyToTemporaryFilepath(filename)
    return
}

func (server *EbooksFileServer) GetEbook(ebookFilename string) (ebook []byte, err error) {
    ebookPath := server.config.StorageDir + "/" + ebookFilename
    ebook, err = ioutil.ReadFile(ebookPath)
    return ebook, err
}

func (server *EbooksFileServer) GetEbooks() (ebooks []Ebook) {
    server.checkForNewBooks()
    ebookPaths := server.getEbookPaths([]string{server.config.StorageDir})
    log.Printf("Getting details from %d ebooks", len(ebookPaths))
	for _, ebookPath := range ebookPaths {
        //grab author, title, rights, description, filename (not whole path!)
        log.Printf("Sending book %s details", ebookPath)
        ebookDetails, err := server.getBookDetails(ebookPath)
        if err != nil {
            log.Printf("error trying to open %s : %s 😱", ebookPath)
            log.Print(err)
        }
        ebooks = append(ebooks, ebookDetails)
    }
    return
}

func (server *EbooksFileServer) getBookDetails(ebookPath string) (ebook Ebook, err error) {
    openEpub, err := epub.Open(ebookPath)
    if err != nil {
        return ebook, err
    }
    metadata := openEpub.Opf.Metadata // .opf file contains the ebook's metadata
    if len(metadata.Creator) > 0 {
        ebook.Author = strings.TrimSpace(metadata.Creator[0].Data)
    } else {
        ebook.Author = "unknown author"
    }
    if len(metadata.Title) > 0 {
        ebook.Title = strings.TrimSpace(metadata.Title[0].Data)
    } else {
        ebook.Title = "unknown title"
    }
    if len(metadata.Rights) > 0 {
        ebook.Rights = strings.TrimSpace(metadata.Rights[0])
    } else {
        ebook.Rights = ""
    }
    if len(metadata.Description) > 0 {
        ebook.Description = strings.TrimSpace(metadata.Description[0])
    } else {
        ebook.Description = ""
    }

    splitEbookPath := strings.Split(ebookPath, `/`)
    ebook.Filename = splitEbookPath[len(splitEbookPath)-1]

    defer openEpub.Close()

    return ebook, nil
}

func (server *EbooksFileServer) checkForNewBooks() {
    log.Print("Checking for new ebooks...")
    log.Printf("Searching in: %s", strings.Join(server.config.SearchDirs, "; "))
    ebookSearchPaths := server.getEbookPaths(server.config.SearchDirs)
    log.Printf("Ebooks from search path: %s", strings.Join(ebookSearchPaths, "; "))
    ebookStoragePaths := server.getEbookPaths([]string{server.config.StorageDir})
    log.Print("Hashing ebooks for comparison...")
    hashedSearchEbooks := server.hashEbooks(ebookSearchPaths)
    hashedStorageEbooks := server.hashEbooks(ebookStoragePaths)
    log.Print("Comparing hashes...")
    // find hashes unique in the search path only
    searchEbooksUnique := server.hashedPathsDifference(hashedSearchEbooks, hashedStorageEbooks)
    if len(searchEbooksUnique) > 0 {
        server.copyEbooks(searchEbooksUnique, server.config.StorageDir)
        log.Print("...finished copying.")
    } else {
        log.Print("...no new ebooks found.")
    }
}

func (server *EbooksFileServer) hashEbooks(ebookPaths []string) (hashedPaths map[string]string) {
    // hash all the ebooks at the given paths
    hashedPaths = make(map[string]string)
	for _, path := range ebookPaths {
        ebookData, err := os.Open(path)
        if err != nil {
            log.Print(err)
        }
        hash := md5.New()
        if _, err := io.Copy(hash, ebookData); err != nil {
            log.Print(err)
        }
        hashString := fmt.Sprintf("%x", hash.Sum(nil))
        hashedPaths[hashString] = path
    }
    return
}

func (server *EbooksFileServer) hashedPathsDifference(a map[string]string, b map[string]string) (difference map[string]string) {
    // returns a - b (i.e. what's in a but not in b)
    // O(n^2) but good enough for now
    difference = make(map[string]string)
    for hash_a, path_a := range a {
        matchFound := false
        for hash_b, _ := range b {
            if hash_a == hash_b {
                matchFound = true
                break
            }
        }
        if !matchFound {
            difference[hash_a] = path_a
        }
    }
    return
}

func (server *EbooksFileServer) copyEbooks(searchEbooksUnique map[string]string, StorageDir string) {
    log.Printf("Found %d new ebooks!", len(searchEbooksUnique))
    for _, path := range searchEbooksUnique {
        filename := server.generateNewFilename()
        log.Printf("Copying %s to %s/%s", path, StorageDir, filename)
        server.Copy(path, StorageDir + "/" + filename)
        //err := server.Copy(path, StorageDir + "/" + filename)
        // do something with err I guess
    }
}

func (server *EbooksFileServer) generateNewFilename() string {
    // generate new filename of format `e00001.epub`
    ebookPaths := server.getEbookPaths([]string{server.config.StorageDir})
    highestNum := 0
	ebookRe, _ := regexp.Compile(`^.*\/+e([\d]{6})\.epub$`)
    for _, ebookPath := range ebookPaths {
        reMatch := ebookRe.FindStringSubmatch(ebookPath)
        if len(reMatch) > 0 {
            digit, _ := strconv.Atoi(reMatch[1])
            if digit > highestNum {
                highestNum = digit
            }
        }
    }
    return fmt.Sprintf("e%06d.epub", highestNum + 1)
}

func (server *EbooksFileServer) getEbookPaths(dirs []string) (ebookPaths []string) {
    // grab all the .epub paths in the storage dir
	for _, dir := range dirs {
        files, err := ioutil.ReadDir(dir)
        //log.Printf("Files found in dir %s: %s", dir, strings.Join(files, "; "))
        if err != nil {
            log.Print(err)
        }
        for _, file := range files {
            filename := file.Name()
            // would be better to check type in another way but this will do for now
            // e.g. from wikipedia: "The first file in the archive must be the mimetype file. It must be unencrypted and uncompressed so that non-ZIP utilities can read the mimetype. The mimetype file must be an ASCII file that contains the string "application/epub+zip". This file provides a more reliable way for applications to identify the mimetype of the file than just the .epub extension."
            plopp := filename[len(filename)-5:]
            if plopp == ".epub" {
                ebookPaths = append(ebookPaths, dir + "/" + filename)
            }
        }
    }
    return
}

// https://stackoverflow.com/a/21061062
// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func (server *EbooksFileServer) Copy(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }
    return out.Close()
}
