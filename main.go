package main

import (
	"delayer-client-golang/delayer"
	"fmt"
)

func main() {
	cli := delayer.Client{
		Host:     "127.0.0.1",
		Port:     "6379",
		Database: 0,
		Password: "",
	}
	cli.Init()
	ok, err := cli.Remove("4776797fba0073c6bdd966e99649754c")
	fmt.Println(ok)
	fmt.Println(err)
}
