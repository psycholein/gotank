//go:generate sh -c "rm -f public/assets/*"
//go:generate coffee -o public/assets/ -c assets/
//go:generate esc -o assets.go -prefix=public/ public/
package main

import (
  "log"
	"net/http"
)

const (
  Addr = ":3000"
)

func main() {
  startServer()
}

func startServer() {
	routes()
	if err := http.ListenAndServe(Addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func routes() {
	//http.Handle("/", http.FileServer(assetFS()))
	http.HandleFunc("/ws", wsHandler)
}

// func assetFS() http.FileSystem {
// 	return FS(*Debug)
// }
