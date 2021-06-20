// 这个聊天室的设计是每个登录成功的客户端都进入一个大的聊天室
package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn        // 每个客户端的socket
	send     chan *message          // 客户端发送的数据
	room     *room                  // 每个客户端要知道把数据发送到哪个room
	userData map[string]interface{} // 每个客户的个人信息
}

// 将接收到的消息直接传给room
func (c *client) read() {
	// 不确定什么时候函数结束要做的清理工作，可以使用defer
	// defer的运行时性能比直接在函数结束时要差一些，但是代码的整洁性更重要
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		fmt.Printf("receive: %s\n", msg)
		c.room.forward <- msg
	}
}

// 从send接收信息然后向socket写入信息
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
