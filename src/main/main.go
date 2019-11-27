//main 入口
package main

import (
	"flag"
	"log"
	"net/http"
	"../socket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", socket.Echo)
	http.HandleFunc("/websocket", socket.Home)
	http.HandleFunc("/chat", socket.Index)
	http.HandleFunc("/chatSocket", socket.ChatSocket)
	log.Fatal(http.ListenAndServe(*addr, nil))
}