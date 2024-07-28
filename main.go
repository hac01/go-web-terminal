package main

import (
	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	cmd := exec.Command("/bin/bash")
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		log.Println("PTY start error:", err)
		return
	}
	defer ptyFile.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			if _, err := ptyFile.Write(message); err != nil {
				log.Println("PTY write error:", err)
				return
			}
		}
	}()

	buf := make([]byte, 1024*20)
	for {
		n, err := ptyFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("PTY read error:", err)
			}
			return
		}
		if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
			log.Println("WebSocket write error:", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	server := &http.Server{
		Addr:         "0.0.0.0:8088",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("Starting server on :8088")
	log.Fatal(server.ListenAndServe())
}
