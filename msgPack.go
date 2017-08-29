package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/vmihailenco/msgpack.v2"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:9066", "http service address")

var messages = make(chan string)

func ByteSlice(b []byte) []byte { return b }

func signHold() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("ctrl+c关闭:", s)
	messages <- "end"
}

func toClient(uid int) {
	fmt.Println("uid:" + strconv.Itoa(uid))

	u := url.URL{Scheme: "ws", Host: *addr, Path: ""}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("dial:", err)
	}
	defer c.Close()

	go func() {
		defer c.Close()

		for {
			_, message, _ := c.ReadMessage()
			var o interface{}
			msgpack.Unmarshal(message, &o)
			fmt.Println("=============================")
			fmt.Println(o)
			fmt.Println("=============================")
		}
	}()
	
	//模拟登陆
	in := []string{"10100", strconv.Itoa(uid), strconv.Itoa(1), strconv.Itoa(2), strconv.Itoa(1), strconv.Itoa(1), strconv.Itoa(1)}
	bb, err := msgpack.Marshal(in)
	fmt.Println(string(bb))
	err = c.WriteMessage(websocket.TextMessage, bb)
	if err != nil {
		log.Println("write:", err)
	}

	log.Println("write OK  1")

	//rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(time.Duration(RandInt64(1, 3)) * time.Second)
	defer ticker.Stop()

	ticker4 := time.NewTicker(time.Duration(RandInt64(4, 8)) * time.Second)
	defer ticker4.Stop()

	ticker2 := time.NewTicker(time.Duration(RandInt64(8, 13)) * time.Second)
	defer ticker2.Stop()

	ticker6 := time.NewTicker(time.Duration(RandInt64(15, 20)) * time.Second)
	defer ticker6.Stop()

	for {
		select {
		case t := <-ticker.C:
			t.String()
			enterBase := []string{"10101"} //前后端拟定的协议ID
			cc, _ := msgpack.Marshal(enterBase)
			err = c.WriteMessage(websocket.TextMessage, cc)
			log.Println("write OK  2")
		case t4 := <-ticker4.C:
			t4.String()
			enterStep := []string{"10103", strconv.Itoa(999)}
			dd, _ := msgpack.Marshal(enterStep)
			err = c.WriteMessage(websocket.TextMessage, dd)
			log.Println("write OK  5")
			ticker4.Stop()
		case t2 := <-ticker2.C:
			t2.String()
			monsterInfo := []string{"30086"}
			dd, _ := msgpack.Marshal(monsterInfo)
			err = c.WriteMessage(websocket.TextMessage, dd)
			log.Println("write OK  3")
			ticker2.Stop()
		case t6 := <-ticker6.C:
			t6.String()
			fightInfo := []string{"10107"}
			dd, _ := msgpack.Marshal(fightInfo)
			err = c.WriteMessage(websocket.TextMessage, dd)
			log.Println("write OK  4")
			ticker6.Stop()
		}
	}
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func main() {

	i := 6000

	go signHold()
	
	//3000个连接发送数据
	for i < 9000 {
		go toClient(i)
		i = i + 1
	}

	<-messages

	fmt.Println("Test stop")

}
