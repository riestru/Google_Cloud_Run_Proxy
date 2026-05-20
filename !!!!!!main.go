package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
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
		http.Error(w, "Not found", 404)
		return
	}

	// Подключаемся к бэкенду с таймаутом 10 секунд
	dst, err := net.DialTimeout("tcp", target, 10*time.Second)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", target, err)
		http.Error(w, "Backend unavailable", 502)
		return
	}

	// Включаем TCP keepalive на соединении к бэкенду
	if tcpDst, ok := dst.(*net.TCPConn); ok {
		tcpDst.SetKeepAlive(true)
		tcpDst.SetKeepAlivePeriod(30 * time.Second)
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		dst.Close()
		http.Error(w, "Hijacking not supported", 500)
		return
	}

	src, _, err := hijacker.Hijack()
	if err != nil {
		dst.Close()
		log.Printf("Hijack failed: %v", err)
		return
	}

	// Включаем TCP keepalive на соединении клиента
	if tcpSrc, ok := src.(*net.TCPConn); ok {
		tcpSrc.SetKeepAlive(true)
		tcpSrc.SetKeepAlivePeriod(30 * time.Second)
	}

	// Пересылаем HTTP запрос на бэкенд
	r.Write(dst)

	// Двунаправленный туннель
	done := make(chan struct{}, 2)
	go func() {
		defer func() { done <- struct{}{} }()
		io.Copy(dst, src)
		dst.Close()
	}()
	go func() {
		defer func() { done <- struct{}{} }()
		io.Copy(src, dst)
		src.Close()
	}()
	<-done
}
