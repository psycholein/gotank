//go:generate coffee -o public/ -c assets/
//go:generate sass --update assets:public
//go:generate haml assets/index.haml public/index.html
//go:generate sh -c "cd assets && find . -name '*.ect' -exec cp --parents \\{\\} ../public \\;"
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
	addr  = ":3000"
	debug = false
)

func main() {
	fmt.Println("Start...")
	embd.InitGPIO()

	handleCtrlC()
	event.InitEvents()
	registerModules()
	initModules()
	startServer()
}

func handleCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		end()
		os.Exit(0)
	}()
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
	http.Handle("/", http.FileServer(assetFS()))
	http.HandleFunc("/ws", wsHandler)
}

func assetFS() http.FileSystem {
	return FS(debug)
}
