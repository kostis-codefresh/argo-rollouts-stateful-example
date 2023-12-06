package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type InterestApplication struct {
	AppVersion        string
	CurrentRole       string
	RabbitHost        string
	RabbitPort        string
	RabbitReadQueue   string
	RabbitWriteQueue  string
	MessagesProcessed int
}

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	interestApp := InterestApplication{}

	interestApp.AppVersion = os.Getenv("APP_VERSION")
	if len(interestApp.AppVersion) == 0 {
		interestApp.AppVersion = "dev"
	}

	interestApp.CurrentRole = "demoRole"
	interestApp.RabbitHost = "demoHost"
	interestApp.RabbitPort = "demoPort"
	interestApp.RabbitReadQueue = "readExample"
	interestApp.RabbitWriteQueue = "writeExample"
	interestApp.MessagesProcessed = 42

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, interestApp.AppVersion)
	})

	// Kubernetes check if app is ok
	http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "up")
	})

	// Kubernetes check if app can serve requests
	http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yes")
	})

	http.HandleFunc("/api/v1/interest", func(w http.ResponseWriter, r *http.Request) {
		randomSource := rand.NewSource(time.Now().UnixNano())
		calculatedInterest := rand.New(randomSource)
		fmt.Fprint(w, (calculatedInterest.Intn(26) + 10))
	})

	http.HandleFunc("/", interestApp.serveFiles)

	fmt.Printf("Backend version %s is listening now at port %s\n", interestApp.AppVersion, port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatal(err)
}

func (interestApp *InterestApplication) serveFiles(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	p := "." + upath
	if p == "./" {
		interestApp.home(w, r)
		return
	} else {
		p = filepath.Join("./static/", path.Clean(upath))
	}
	http.ServeFile(w, r, p)
}

func (interestApp *InterestApplication) home(w http.ResponseWriter, r *http.Request) {

	// loanApp.findBackendVersion()

	t, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}
	err = t.Execute(w, interestApp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}
