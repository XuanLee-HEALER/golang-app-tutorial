// 这个聊天室的设计是每个登录成功的客户端都进入一个大的聊天室
package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type client struct {
	// 每个客户端的socket
	socket *websocket.Conn
	// 客户端发送的数据
	send chan []byte
	// 每个客户端要知道把数据发送到哪个room
	room *room
}

// 将接收到的消息直接传给room
func (c *client) read() {
	// 不确定什么时候函数结束要做的清理工作，可以使用defer
	// defer的运行时性能比直接在函数结束时要差一些，但是代码的整洁性更重要
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		fmt.Printf("receive: %s\n", msg)
		c.room.forward <- msg
	}
}

// 从send接收信息然后向socket写入信息
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		fmt.Printf("write to socket: %s\n", msg)
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
