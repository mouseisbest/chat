package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	fmt.Println("gnr---聊天室")
	fmt.Println("...........................................................")
	ttest()

}

func ttest() {
	client()
}

var wg sync.WaitGroup

func client() {
	address := "127.0.0.1:37001"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("客户端连接服务器出错:", err)
	}
	defer conn.Close()
	wg.Add(1)
	go readMsg(conn) //开启单独的一个gorountine 去处理接收消息
	//主 goroutine 执行发送消息

	inputReader := bufio.NewReader(os.Stdin)

	for {
		//fmt.Print("请输入将要发送的数据(按enter键结束):")
		inputinfos, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("输入信息发生错误:", err)
			continue
		}
		inputinfos = strings.Trim(strings.Trim(inputinfos, "\n"), "\r")
		if strings.ToUpper(inputinfos) == "Q" || strings.ToUpper(inputinfos) == "QUIT" || strings.ToUpper(inputinfos) == "EXIT" {
			conn.Write([]byte("EXIT")) //发送请求断开链接的消息 退出循环
			break
		} else if len(inputinfos) == 0 {
			continue
		}
		//发送数据到服务器
		_, err = conn.Write([]byte(inputinfos))
		if err != nil {
			fmt.Println("数据发生出错 :", err)
			break
		}

	}
	wg.Wait()
	return

}

func readMsg(conn net.Conn) {
	defer wg.Done()

	for {

		//接收服务端返回来的返回值
		var buffer [128]byte

		n, err := conn.Read(buffer[:])
		if err != nil {
			fmt.Println("接收数据发生错误 :", err)
			return //安全断开客户端
		}
		if n > 0 {
			if string(buffer[:n]) == "EXIT" {
				return //安全断开客户端 客户端退出当前groutine
			}
			fmt.Println(string(buffer[:n]))
		}

	}

}
