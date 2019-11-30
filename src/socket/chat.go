package socket

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	Username string
	Message  string
	Lasttime int64
}

// 全局信息
var messages []Message
var newmsg Message
var user map[string]string

func Index(w http.ResponseWriter, r *http.Request) {
	homeTemplate, _ := template.ParseFiles("../../web/chat.html")
	homeTemplate.Execute(w, "ws://"+r.Host+"/chatSocket")
}

func ChatSocket(w http.ResponseWriter, r *http.Request) {
	var msg Message

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 进入聊天室把历史聊天记录进行信息分发
	for _, val := range messages {
		s1 := []byte("{\"username\":\"" + val.Username + "\",\"message\":\"" + val.Message + "\",\"lasttime\":\"" + strconv.FormatInt(val.Lasttime, 10) + "\"}")
		err = c.WriteMessage(1, s1)
		if err != nil {
			fmt.Println("发送出错: " + err.Error())
			break
		}
	}
	defer c.Close()

	go func() {
		for {
			log.Println("当前信息",newmsg)
			time.Sleep(time.Second / 2)
			if newmsg.Username != "" {
				log.Println("当前用户",user)
				err :=c.WriteMessage(1, []byte("{\"username\":\""+newmsg.Username+"\",\"message\":\""+newmsg.Message+"\",\"lasttime\":\""+strconv.FormatInt(newmsg.Lasttime, 10)+"\"}"))
				if err != nil {
					fmt.Println("发送出错: " + err.Error())
					c.Close()
					break
				}
				time.Sleep(time.Second)
				newmsg = Message{
					Username: "",
					Message:  "",
					Lasttime: 0,
				}
			}
		}
	}()

	for {
		// 接收数据
		_, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		// 解析信息
		err = json.Unmarshal([]byte(message), &msg)

		// 添加新用户到map中,已经存在的用户不必添加
		if user != nil {
			if _, ok := user[msg.Username]; !ok {
				user[msg.Username] = msg.Username
			}
		}else {
			user = make(map[string]string, 10)
			user[msg.Username] = msg.Username
		}

		// 添加聊天记录到全局信息
		msg.Lasttime = time.Now().Unix()
		newmsg = msg
		messages = append(messages, msg)
		log.Println("聊天记录",messages)
		log.Println("当前聊天记录",newmsg)
	}
}
