//go:generate sh -c "rm -f public/assets/*"
//go:generate coffee -o public/assets/ -c assets/
//go:generate esc -o assets.go -prefix=public/ public/
package main

import (
	"fmt"
	"gotank/event"
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	addr = ":3000"
)

func main() {
	fmt.Println("Start...")
	embd.InitGPIO()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		end()
		os.Exit(0)
	}()

	event.InitEvents()
	registerModules()
	initModules()
	startServer()
}

func end() {
	embd.CloseGPIO()
	stopModules()
	event.Stop()
}

func startServer() {
	routes()
	if err := http.ListenAndServe(addr, nil); err != nil {
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
