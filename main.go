package main

import (
    "fmt"
    "html/template"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

var Version string

// PageData rappresenta i dati da passare al template
type PageData struct {
    Files []string
}

func main() {

    fmt.Println("Version:", Version)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("templates/home.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })

    http.HandleFunc("/uploadpage", func(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("templates/upload.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })

    http.HandleFunc("/viewpage", func(w http.ResponseWriter, r *http.Request) {
        entries, err := os.ReadDir("./file")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fileNames := make([]string, len(entries))
        for i, entry := range entries {
            fileNames[i] = entry.Name()
        }
        data := PageData{
            Files: fileNames,
        }
        tmpl, err := template.ParseFiles("templates/view.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        err = tmpl.Execute(w, data)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })

    http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
            return
        }
        file, header, err := r.FormFile("fileUpload")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer file.Close()
        dst, err := os.Create("./file/" + filepath.Base(header.Filename))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer dst.Close()
        _, err = io.Copy(dst, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    http.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
        fileName := r.URL.Path[len("/file/"):]
        file, err := os.Open("./file/" + fileName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        defer file.Close()
        w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
        w.Header().Set("Content-Type", "application/octet-stream")
        _, err = io.Copy(w, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })

    cert := "cert/cert.pem"
    key := "cert/key.pem"

    log.Println("Avvio del server sulla porta 443")
    err := http.ListenAndServeTLS(":443", cert, key, nil)
    if err != nil {
        log.Fatal(err)
    }
}