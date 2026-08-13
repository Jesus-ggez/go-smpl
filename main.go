package main

import (
    "embed"
    "html/template"
    "io/fs"
    "log"
    "net/http"
    "os"
    "sync"
)

//go:embed templates/*.html
//go:embed static/css/*.css
var templateFS embed.FS;


func rootTempl() *template.Template {
    t, err := template.ParseFS(templateFS, "templates/*.html");

    if err != nil {
        panic(err);
    }

    return t
}


type Counter struct {
    mu      sync.RWMutex
    Value int
}

var globCounter = &Counter{ Value: 0 };

func main() {
    log.Print("init app")

    staticFS, err := fs.Sub(templateFS, "static");

    if err != nil {
        log.Fatal("Error loading static files: ", err);
    }

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))));

    http.HandleFunc("/", handleHome);

    http.HandleFunc("/api/counter", handleCounter);


    port := os.Getenv("PORT");

    if port == "" {
        port = "8080";
    }
    log.Fatal(http.ListenAndServe(":" + port, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {

    globCounter.mu.RLock()

    currentValue := globCounter.Value

    globCounter.mu.RUnlock()

    data := struct {
        Title string
        Value int
    }{Title: "Counter App", Value: currentValue}

    if err := rootTempl().ExecuteTemplate(w, "index", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func handleCounter(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    action := r.FormValue("action")

    globCounter.mu.Lock()

    switch action {
    case "sum":
        globCounter.Value++

    case "res":
        globCounter.Value--

    case "reset":
        globCounter.Value = 0

    default:
        globCounter.mu.Unlock()
        http.Error(w, "Invalid action", http.StatusBadRequest)
        return
    }

    newValue := globCounter.Value
    globCounter.mu.Unlock()

    http.Redirect(w, r, "/", http.StatusSeeOther)

    log.Print("Counter updated: action={}, newValue={}", action, newValue)
}
