package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"golang-fave/engine/wrapper"
)

const C_Debug = !false

var FParamHost string
var FParamPort int
var FParamWwwDir string
var FVhostHomeDir string

func init() {
	flag.StringVar(&FParamHost, "host", "0.0.0.0", "server host")
	flag.IntVar(&FParamPort, "port", 8080, "server port")
	flag.StringVar(&FParamWwwDir, "dir", "", "virtual hosts directory")
	flag.Parse()
}

func main() {
	if _, err := os.Stat(FParamWwwDir); os.IsNotExist(err) {
		fmt.Println("Virtual hosts directory is not exists")
		fmt.Println("Example: ./fave -host 127.0.0.1 -port 80 -dir ./hosts")
		return
	}
	if FParamWwwDir[len(FParamWwwDir)-1] != '/' {
		FParamWwwDir = FParamWwwDir + "/"
	}

	// Handle
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", FParamHost, FParamPort),
		Handler: mux,
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Starting server at %s:%d", FParamHost, FParamPort)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	// Wait for signal
	<-stop

	log.Printf("Shutting down server...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Server is off!")
	}
}

func vhExists(vhosthome string) bool {
	if st, err := os.Stat(vhosthome); !os.IsNotExist(err) {
		if err == nil {
			fmode := st.Mode()
			if fmode.IsDir() {
				return true
			}
		}
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Build vhost home dir
	host := r.Host
	port := strconv.Itoa(FParamPort)
	index := strings.Index(host, ":")
	if index > -1 {
		port = host[index+1:]
		host = host[0:index]
	}

	// Cut "www" if exists
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}
	FVhostHomeDir = FParamWwwDir + host

	// Check if virtual host exists
	if !vhExists(FVhostHomeDir) {
		host = "localhost"
		FVhostHomeDir = FParamWwwDir + host
	}

	// Set protocol
	r.URL.Scheme = "http"

	// Set server name
	w.Header().Set("Server", "fave.pro")

	// Create and start engine
	wrapper.New(&w, r, host, port, FParamWwwDir, FVhostHomeDir, C_Debug).
		Run(func(e *wrapper.Wrapper) bool {
			return false
		})
}
