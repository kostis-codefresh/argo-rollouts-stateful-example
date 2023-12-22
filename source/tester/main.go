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
	"sync"
	"time"
)

type TesterApplication struct {
	RabbitHost         string
	RabbitPort         string
	RabbitQueue        string
	RabbitPreviewHost  string
	RabbitPreviewPort  string
	RabbitPreviewQueue string

	mu                          sync.RWMutex
	ProductionMessagesProcessed int
	PreviewMessagesProcessed    int
	ProductionMessagesSent      int
	PreviewMessagesSent         int
}

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "7000"
	}

	testerApp := TesterApplication{}
	testerApp.RabbitHost = "localhost"
	testerApp.RabbitPort = "5672"
	testerApp.RabbitQueue = "myProductionQueue"
	testerApp.RabbitPreviewHost = "localhost"
	testerApp.RabbitPreviewPort = "5672"
	testerApp.RabbitPreviewQueue = "myPreviewQueue"

	// Kubernetes check if app is ok
	http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "up")
	})

	// Kubernetes check if app can serve requests
	http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yes")
	})

	//Sends a message to production
	http.HandleFunc("/production", func(w http.ResponseWriter, r *http.Request) {
		testerApp.publishProductionMessage()
		fmt.Fprintf(w, "Sent %d", testerApp.ProductionMessagesSent)
	})

	http.HandleFunc("/preview", func(w http.ResponseWriter, r *http.Request) {
		testerApp.publishPreviewMessage()
		fmt.Fprintf(w, "Sent %d", testerApp.PreviewMessagesSent)
	})

	http.HandleFunc("/api/v1/interest", func(w http.ResponseWriter, r *http.Request) {
		randomSource := rand.NewSource(time.Now().UnixNano())
		calculatedInterest := rand.New(randomSource)
		fmt.Fprint(w, (calculatedInterest.Intn(26) + 10))
	})

	http.HandleFunc("/", testerApp.serveFiles)

	fmt.Printf("Tester is listening now at port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func (testerApp *TesterApplication) serveFiles(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	p := "." + upath
	if p == "./" {
		testerApp.home(w, r)
		return
	} else {
		p = filepath.Join("./static/", path.Clean(upath))
	}
	http.ServeFile(w, r, p)
}

func (testerApp *TesterApplication) home(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}
	err = t.Execute(w, testerApp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}
