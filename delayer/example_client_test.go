package delayer_test

import (
	"fmt"
	"delayer-client-golang/delayer"
	"crypto/md5"
	"time"
	"github.com/gomodule/redigo/redis"
)

// 例子
func Example() {
	// 通过连接信息创建客户端
	cli := delayer.Client{
		Host:     "127.0.0.1",
		Port:     "6379",
		Database: 0,
		Password: "",
	}
	cli.Init()

	// 通过已有连接创建客户端
	pool := redis.Pool{}
	conn := pool.Get();
	cli1 := delayer.Client{
		Conn:conn,
	}
	cli1.Init()
}

// Push 例子
func ExampleClient_Push() {
	cli := delayer.Client{
		Host:     "127.0.0.1",
		Port:     "6379",
		Database: 0,
		Password: "",
	}
	cli.Init()
	msg := delayer.Message{
		ID:    fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String()))),
		Topic: "test",
		Body:  "test body",
	}
	reply, err := cli.Push(msg, 10, 600)
	fmt.Println(msg)
	fmt.Println(reply)
	fmt.Println(err)
}

// Pop 例子
func ExampleClient_Pop() {
	cli := delayer.Client{
		Host:     "127.0.0.1",
		Port:     "6379",
		Database: 0,
		Password: "",
	}
	cli.Init()
	msg, err := cli.Pop("test");
	fmt.Println(msg)
	fmt.Println(err)
}

// BPop 例子
func ExampleClient_BPop() {
	cli := delayer.Client{
		Host:     "127.0.0.1",
		Port:     "6379",
		Database: 0,
		Password: "",
	}
	cli.Init()
	msg, err := cli.BPop("test", 10);
	fmt.Println(msg)
	fmt.Println(err)
}

// Remove 例子
func ExampleClient_Remove() {
	cli := delayer.Client{
		Host:     "127.0.0.1",
		Port:     "6379",
		Database: 0,
		Password: "",
	}
	cli.Init()
	ok, err := cli.Remove("9a8482a06630840ce7da9da62d748b8a")
	fmt.Println(ok)
	fmt.Println(err)
}
