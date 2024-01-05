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
type Data struct {
    SpaceUsed float64
    SpaceTotal float64
    Files []string
}

func DirSize(path string) (int64, error) {
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            size += info.Size()
        }
        return err
    })
    return size, err
}

func logRequests(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Received %s request for %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
        next.ServeHTTP(w, r)
    })
}

func main() {

    fmt.Println("Version:", Version)

    mux := http.NewServeMux()

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("templates/index.html")
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

    mux.HandleFunc("/uploadpage", func(w http.ResponseWriter, r *http.Request) {
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

    mux.HandleFunc("/viewpage", func(w http.ResponseWriter, r *http.Request) {
        entries, err := os.ReadDir("./file")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fileNames := make([]string, len(entries))
        for i, entry := range entries {
            fileNames[i] = entry.Name()
        }
        size, err := DirSize("./file")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        totalSize := int64(8589934592) // Sostituisci con il valore reale
        data := Data{
            Files: fileNames,
            SpaceUsed: float64(size) / 1024 / 1024 / 1024,
            SpaceTotal: float64(totalSize) / 1024 / 1024 / 1024,
        }
        tmpl := template.New("view.html").Funcs(template.FuncMap{
            "format": func(v float64) string {
                return fmt.Sprintf("%.2f", v)
            },
            "div": func(a, b float64) float64 {
                return a / b
            },
            "mul": func(a, b float64) float64 {
                return a * b
            },
        })
        tmpl, err = tmpl.ParseFiles("templates/view.html")
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

    mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
        file := r.URL.Query().Get("file")
        err := os.Remove("./file/" + file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/viewpage", http.StatusSeeOther)
    })

    mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
            return
        }
        err := r.ParseMultipartForm(10 << 20) // 10 MB
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        files := r.MultipartForm.File["fileUpload"]
        for _, fileHeader := range files {
            file, err := fileHeader.Open()
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            defer file.Close()
            dst, err := os.Create("./file/" + filepath.Base(fileHeader.Filename))
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
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    mux.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
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

    mux.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("media"))))

    cert := "cert/cert.pem"
    key := "cert/key.pem"

    log.Println("Avvio del server sulla porta 443")
    err := http.ListenAndServeTLS("0.0.0.0:443", cert, key, logRequests(mux))
    if err != nil {
        log.Fatal(err)
    }
}