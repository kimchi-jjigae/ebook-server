# E-book server

I thought it was annoying to manually transfer e-book files from my computer/server to my Kobo eReader, so I wrote this to easily access and download them through the Kobo browser.  
It's also my first project in Go!

When `/api/ebooks/` is queried with the correct password:
- looks for all epub files (this is the only format I've been using but this could easily be expanded upon) in the given "search" directories
- checks if they already exist in the "storage" directory or not
- copies over the non-existent files there 
- returns a json object of all the epub files in the "storage" directory, with various information (author, title, description) as extracted from the epub format

This info can be rendered as wished; I have a simple web page for displaying them in a table with a download link beside them.

When `/api/ebook/:filename` is queried with the correct password:
- the relevant epub file is copied to a temporary file, accessible as a URL for the next 60 seconds
- the temporary URL is returned

This isn't ideal since the temporary URL is publicly accessible. Originally I had the server return the actual binary epub file, downloadable as a blob, but I had a lot of trouble trying to get the Kobo to download blobs with the very limited JS support that its in-built browser has. This is why I've opted for this solution.

I haven't added any front-end code in this repo but I may do at a later date in the future!
The password should be sent to the server through an `x-password` header. I realise now that there are [built in basic HTTP frameworks for this](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication), but this works at least :)

## Some config before building:
- If you don't want to use HTTPS: uncomment `http.ListenAndServe` and comment `http.ListenAndServeTLS` in `main.go`
- Edit the `config.toml` file as necessary (remove `Certificate` and `Key` if not using HTTPS
- Place the `config.toml` file in the dir you run the executable from
