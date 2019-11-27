package socket

import (
	"encoding/json"
	"fmt"
	"html/template"
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

	//defer c.Close()
	for {
		// 接收数据
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		// 解析信息
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			fmt.Println("解析数据异常")
		}
		fmt.Println(message)
		//获取到前端最后阅读时间
		lasttime := msg.Lasttime

		// 添加聊天记录到全局信息
		msg.Lasttime = time.Now().Unix()
		messages = append(messages, msg)

		// 通过webSocket将当前信息分发
		for _, val := range messages {
			s1 := []byte("{\"username\":\"" + val.Username + "\",\"message\":\"" + val.Message + "\",\"lasttime\":\"" + strconv.FormatInt(val.Lasttime, 10) + "\"}")
			fmt.Println(val)
			fmt.Println(lasttime)
			if val.Lasttime > lasttime {
				err = c.WriteMessage(mt, s1)
				if err != nil {
					fmt.Println("发送出错: " + err.Error())
					break
				}
			}
		}

	}
}
