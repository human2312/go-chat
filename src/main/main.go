//main 入口
package main

import (
	"../socket"
	"fmt"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Wrold!<br/><a herf='/chat'>chat room</a>") //这个写入到w的是输出到客户端的
}

func main() {
	http.HandleFunc("/", index) //设置访问的路由
	http.HandleFunc("/echo", socket.Echo)
	http.HandleFunc("/websocket", socket.Home)
	http.HandleFunc("/chat", socket.Index)
	http.HandleFunc("/chatSocket", socket.ChatSocket)
	err := http.ListenAndServe(":8080", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
