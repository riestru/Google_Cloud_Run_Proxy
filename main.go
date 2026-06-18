package main
import (
        "io"
        "net"
        "os"
)
func main() {
        listenAddr := ":" + os.Getenv("PORT")
	    targetIP := os.Getenv("V2RAY_SERVER_IP")
	    targetPort := os.Getenv("V2RAY_SERVER_PORT")
	    targetAddr := net.JoinHostPort(targetIP, targetPort)
        ln, err := net.Listen("tcp", listenAddr)
        if err != nil {
                return
        }
        for {
                conn, err := ln.Accept()
                if err != nil {
                        continue
                }
                go handleConnection(conn, targetAddr)
        }
}
func handleConnection(src net.Conn, targetAddr string) {
        dst, err := net.Dial("tcp", targetAddr)
        if err != nil {
                src.Close()
                return
        }
        go io.Copy(dst, src)
        go io.Copy(src, dst)
}
