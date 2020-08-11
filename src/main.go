package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strings"
)

var local string
var remote string
var ws *websocket.Conn
var err error

func main() {
	flag.StringVar(&local, "l", "", "Listen local port")
	flag.StringVar(&remote, "c", "", "Connect remote websocket server")
	flag.Usage = func() {
		fmt.Println(`
Usage: 
  -l 8080           Listen local port.		
  -c 127.0.0.1:8080 Connect remote websocket server.
		`)
	}
	flag.Parse()

	go listenInput()

	if local != "" {
		createNewListen(parseAddr(local))
		return
	}

	if remote != "" {
		ws = createNewWs(parseAddr(remote))
		listen()
		return
	}

	flag.Usage()
	return
}

func listenInput() {
	// 监听输入
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _, err := reader.ReadLine()

		ws.WriteMessage(websocket.TextMessage, line)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func parseAddr(addr string) string {
	if strings.Index(addr, ".") > -1 {
		return addr
	}
	if len(addr) < 8 && strings.Index(addr, ":") == -1 {
		return "127.0.0.1:" + addr
	}
	return addr
}

func createNewWs(ip string) *websocket.Conn {
	if strings.Index(ip, "ws") == -1 {
		ip = "ws://" + ip
	}
	ws, _, err = websocket.DefaultDialer.Dial(ip, http.Header{})
	if err != nil {
		fmt.Println(err)
	}
	return ws
}

func createNewListen(ip string) {
	http.HandleFunc("/", handler)
	fmt.Println("server created at " + ip)

	if err = http.ListenAndServe(ip, nil); err != nil {
		fmt.Println(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	ws, err = upgrader.Upgrade(w, r, nil)
	go listen()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func listen() {
	for {
		_, bs, err := ws.ReadMessage()

		fmt.Println(string(bs))

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
