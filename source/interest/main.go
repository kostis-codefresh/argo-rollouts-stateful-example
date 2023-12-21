package main

import (
	"container/list"
	"context"
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

type InterestApplication struct {
	AppVersion      string
	CurrentRole     string
	RabbitHost      string
	RabbitPort      string
	RabbitReadQueue string

	mu                sync.RWMutex
	MessagesProcessed int
	LastMessages      *list.List //Assume that last 5 are enough
	dummyCounter      int

	stopNow context.CancelFunc //Used to cancel message reading when conf changes
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

	interestApp.readCurrentConfiguration()

	interestApp.MessagesProcessed = 0
	interestApp.LastMessages = list.New()
	interestApp.startReadingMessages()

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, interestApp.AppVersion)
	})

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, interestApp.MessagesProcessed)
	})

	http.HandleFunc("/role", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, interestApp.CurrentRole)
	})

	// Kubernetes check if app is ok
	http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "up")
	})

	// Kubernetes check if app can serve requests
	http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yes")
	})

	//Sends a dummy message to our own queue (just for testing purposes)
	http.HandleFunc("/dummy", func(w http.ResponseWriter, r *http.Request) {
		interestApp.publishMessage()
		fmt.Fprintf(w, "Sent %d", interestApp.dummyCounter)
	})

	http.HandleFunc("/list", interestApp.listNotifications)

	http.HandleFunc("/api/v1/interest", func(w http.ResponseWriter, r *http.Request) {
		randomSource := rand.NewSource(time.Now().UnixNano())
		calculatedInterest := rand.New(randomSource)
		fmt.Fprint(w, (calculatedInterest.Intn(26) + 10))
	})

	http.HandleFunc("/", interestApp.serveFiles)

	fmt.Printf("Backend version %s is listening now at port %s\n", interestApp.AppVersion, port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
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

func (interestApp *InterestApplication) listNotifications(w http.ResponseWriter, req *http.Request) {
	interestApp.mu.RLock()
	defer interestApp.mu.RUnlock()
	for m := interestApp.LastMessages.Front(); m != nil; m = m.Next() {
		fmt.Fprintf(w, "<div class=\"entry\"><span>%s</span></div>", m.Value)
	}
	fmt.Fprintf(w, "<strong id=\"count\" hx-swap-oob=\"true\">%d</strong>", interestApp.MessagesProcessed)
}
