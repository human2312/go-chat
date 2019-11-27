package socket

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Message struct {
	Username string
	Message  string
}

type User struct {
	Username string
}

type Datas struct {
	Messages []Message
	Users    []User
}

// 全局信息
var datas Datas
var users map[http.ResponseWriter]string

func Index(w http.ResponseWriter, r *http.Request) {
	homeTemplate, _ := template.ParseFiles("../../web/chat.html")
	homeTemplate.Execute(w, "ws://"+r.Host+"/chatSocket")
}

func ChatSocket(w http.ResponseWriter, r *http.Request) {
	var msg Message

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	//defer c.Close()
	for {
		// 接收数据
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		// 解析信息
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			fmt.Println("解析数据异常")
		}

		log.Printf("recv: %s", msg.Username)
		// 通过webSocket将当前信息分发
		err = c.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("发送出错: " + err.Error())
			break
		}

	}
}
