package main

import (
	"encoding/json"
	"fmt"
	pandora "github.com/wk-go/pandora_go"
	"os"
)

func main() {
	//username
	username := os.Getenv("CHAT_USERNAME")
	password := os.Getenv("CHAT_PASSWORD")
	url := os.Getenv("CHAT_URL")
	if len(username) == 0 || len(url) == 0 {
		panic("请输入必要参数")
	}

	client := pandora.NewClient(url, "")

	client.AddHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")

	proxy := os.Getenv("CHAT_PROXY")
	if len(proxy) == 0 {
		proxy = "http://127.0.0.1:8888"
	}
	_, err := client.SetProxy(proxy)

	if err != nil {
		panic(err)
	}

	client.AddHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")

	if err := client.Login(username, password, ""); err != nil {
		panic(err)
	}

	conversations, err := client.ConversationsGET()
	if err != nil {
		panic(err)
	}

	data, _ := json.Marshal(conversations)
	fmt.Println(string(data))
}
