package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var address string

func main() {

	args := os.Args
	if len(args) == 2 {
		address = args[1]
	}
	fmt.Println("gnr聊天室---")
	fmt.Println("...........................................................")
	test()

}

var client_conns sync.Map

func test() {

	tcpServer()

}

func tcpServer() {

	listen, err := net.Listen("tcp", address)
	if err != nil {

		fmt.Println("tcp 监听出现问题:", err)
		return
	}
	fmt.Println("服务器开启成功，正在进行监听,监听地址：tcp://", address)
	for {
		conn, err := listen.Accept() //建立了一个连接 开启一个gorountine 去处理这个连接的会话
		if err != nil {
			fmt.Println("连接accept 失败：", err)
			continue
		}
		//加入聊天室的消息
		client_conns.Range(func(key any, itemconn any) bool {
			myconn := itemconn.(net.Conn)
			_ = key.(string)
			_, err := myconn.Write([]byte(conn.RemoteAddr().String() + "加入了聊天"))
			if err != nil {
				return false
			}
			return true

		})
		client_conns.Store(conn.RemoteAddr().String(), conn)

		go dealConnect(conn)
	}

}

func dealConnect(conn net.Conn) {

	defer conn.Close()
	for {
		fmt.Println("-------------------------")
		fmt.Printf("服务端：%v 客户端：%v  ", conn.LocalAddr().String(), conn.RemoteAddr().String())
		reader := bufio.NewReader(conn)

		var buffer [128]byte
		n, err := reader.Read(buffer[:])
		if err != nil { //一般就是客户端断开连接的处理
			fmt.Println("收取客户端发来的数据出现错误(默认会断开连接)：", err)

			client_conns.Range(func(key any, itemconn any) bool {
				myconn := itemconn.(net.Conn)
				_ = key.(string)
				_, err := myconn.Write([]byte(conn.RemoteAddr().String() + "退出了聊天"))
				if err != nil {
					return false
				}
				return true

			})

			client_conns.Delete(conn.RemoteAddr().String())
			return
		}
		fmt.Println("收到的字节数：", n)
		if n > 0 {
			if string(buffer[:n]) == "EXIT" { //客户端发来了请求退出的消息
				conn.Write(buffer[:n]) //将EXIT消息 原封不动发给请求退出测客户端 客户正常退出groutine
				continue
			}
			receiveStr := conn.RemoteAddr().String() + "说：" + string(buffer[:n])
			fmt.Println("收到客户端发来的信息是：", receiveStr)
			//广播到所有的链接
			client_conns.Range(func(key any, itemconn any) bool {
				myconn := itemconn.(net.Conn)
				_ = key.(string)
				_, err := myconn.Write([]byte(receiveStr))
				if err != nil {
					return false
				}
				return true

			})

		}

	}
}
