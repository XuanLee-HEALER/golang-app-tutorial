package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte      // 持有客户端发送过来的信息，然后转给其它客户端
	join    chan *client     // 加入房间的客户端
	leave   chan *client     // 离开房间的客户端
	clients map[*client]bool // 所有的客户端
}

// 创建实例的helper函数
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

// 在需要同步或者修改共享内存时，可以使用select语句
func (r *room) run() {
	for {
		// 任何channel接收到值，select都会执行对应的case语句块
		// 同一时间只能有一个case被执行，所以可以保证多线程环境下的同步
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			fmt.Printf("room receive: %s\n", msg)
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

// 声明常量可以减少硬编码的内容，随着代码量增长可以放在专门的文件或者集中在文件头部
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// 可被复用，所以只需要创建一次
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin:     nil,
}

func (r *room) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// http请求要升级成为socket请求
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Fatal("ServeHttp:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
