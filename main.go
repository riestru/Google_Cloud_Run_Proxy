package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

var routes map[string]string

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	routes = map[string]string{
		"/panelws1": os.Getenv("V2RAY_SERVER_IP1") + ":8183",
		"/panelws2": os.Getenv("V2RAY_SERVER_IP2") + ":8184",
		"/panelws3": os.Getenv("V2RAY_SERVER_IP3") + ":8185",
		"/panelws4": os.Getenv("V2RAY_SERVER_IP4") + ":8186",
		"/panelws5": os.Getenv("V2RAY_SERVER_IP5") + ":8187",
	}

	log.Printf("Starting proxy on port %s", port)
	for path, target := range routes {
		log.Printf("Route: %s -> %s", path, target)
	}

	http.HandleFunc("/", handleHTTP)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Recovered from panic: %v", rec)
		}
	}()

	target, ok := routes[r.URL.Path]
	if !ok {
		log.Printf("No route for path: %s", r.URL.Path)
		http.Error(w, "Not found", 404)
		return
	}

	if target == ":8183" || target == ":8184" || target == ":8185" || target == ":8186" || target == ":8187" {
		log.Printf("Empty IP for path: %s", r.URL.Path)
		http.Error(w, "Backend not configured", 503)
		return
	}

	log.Printf("Connecting to %s for path %s", target, r.URL.Path)

	dst, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", target, err)
		http.Error(w, "Backend unavailable", 502)
		return
	}
	defer dst.Close()

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		log.Printf("Hijacking not supported")
		http.Error(w, "Hijacking not supported", 500)
		return
	}

	src, _, err := hijacker.Hijack()
	if err != nil {
		log.Printf("Hijack failed: %v", err)
		return
	}
	defer src.Close()

	r.Write(dst)

	done := make(chan struct{}, 2)
	go func() {
		io.Copy(dst, src)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(src, dst)
		done <- struct{}{}
	}()
	<-done
}
