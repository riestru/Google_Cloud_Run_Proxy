package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listenAddr := ":" + port

	targetIP := os.Getenv("V2RAY_SERVER_IP")
	targetPort := os.Getenv("V2RAY_SERVER_PORT")
	if targetIP == "" || targetPort == "" {
		log.Fatal("ОШИБКА: Переменные окружения V2RAY_SERVER_IP и V2RAY_SERVER_PORT должны быть заданы")
	}
	targetAddr := net.JoinHostPort(targetIP, targetPort)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Не удалось запустить прослушивание на %s: %v", listenAddr, err)
	}
	defer ln.Close() // Хороший тон, хотя в бесконечном цикле не вызовется
	
	log.Printf("Прокси запущен на %s, перенаправление на %s", listenAddr, targetAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Ошибка Accept: %v", err)
			continue
		}
		
		// 🚨 ВАЖНО: УБРАЛИ KeepAlive! 
		// На Cloud Run (за GCLB) TCP KeepAlive часто вызывает разрывы соединений.
		// Пусть соединение живет, пока есть реальный трафик.
		
		go handleConnection(conn, targetAddr)
	}
}

func handleConnection(src net.Conn, targetAddr string) {
	defer src.Close()

	// Используем Dialer, но УБИРАЕМ KeepAlive и УВЕЛИЧИВАЕМ Timeout.
	// 5 секунд — это мало для глобальной сети. 15-30 секунд — в самый раз.
	dialer := net.Dialer{
		Timeout: 15 * time.Second, 
	}

	dst, err := dialer.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Нет связи с VPS %s: %v", targetAddr, err)
		return 
	}
	defer dst.Close()
	
	// Убрали логирование успешных подключений. 
	// В Cloud Run логи стоят денег/имеют лимиты. При активном VPN это засорит Cloud Logging.

	// Безопасный канал для завершения (правильный паттерн, оставляем)
	errChan := make(chan error, 2)

	go func() {
		_, err := io.Copy(dst, src)
		errChan <- err
	}()

	go func() {
		_, err := io.Copy(src, dst)
		errChan <- err
	}()

	// Ждем, пока любая из сторон не закроет соединение
	err = <-errChan
	if err != nil && err != io.EOF {
		// Логируем только реальные ошибки, а не штатное отключение клиента
		log.Printf("Соединение разорвано: %v", err)
	}
}
