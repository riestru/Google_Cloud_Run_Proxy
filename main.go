package main

import (
	"io"
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
	}

	http.HandleFunc("/", handleHTTP)
	http.ListenAndServe(":"+port, nil)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	target, ok := routes[r.URL.Path]
	if !ok {
		http.Error(w, "Not found", 404)
		return
	}

	// Подключаемся к бэкенду
	dst, err := net.Dial("tcp", target)
	if err != nil {
		http.Error(w, "Backend unavailable", 502)
		return
	}
	defer dst.Close()

	// Захватываем соединение клиента
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", 500)
		return
	}
	src, _, err := hijacker.Hijack()
	if err != nil {
		dst.Close()
		return
	}
	defer src.Close()

	// Пересылаем HTTP запрос на бэкенд
	r.Write(dst)

	// Двунаправленный туннель
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
