package main

import (
	"JotunBack/server"
	"log"
	"net/http"
)

func main() {
	hub := server.NewHub()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.HandleWebSocket(hub, w, r)
	})

	log.Println("WebSocket server started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
