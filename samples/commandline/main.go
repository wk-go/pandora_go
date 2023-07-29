package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	pandora "github.com/wk-go/pandora_go"
	"os"
	"strings"
	"time"
)

func main() {
	//token
	token := os.Getenv("CHAT_TOKEN")
	url := os.Getenv("CHAT_URL")
	if len(token) == 0 || len(url) == 0 {
		panic("请输入必要参数")
	}

	client := pandora.NewClient(url, token)

	//proxy := os.Getenv("CHAT_PROXY")
	//if len(proxy) == 0 {
	//	proxy = "http://127.0.0.1:8888"
	//}
	//_, err := client.SetProxy(proxy)
	//
	//if err != nil {
	//	panic(err)
	//}

	client.AddHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")

	conversations, err := client.ConversationsGET()
	if err != nil {
		_err := new(pandora.ErrorResponse)
		if errors.As(err, &_err) {
			panic(_err)
		}
		panic(err)
	}
	conversationsJson, _ := json.Marshal(conversations)
	fmt.Println("对话列表:", conversationsJson)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入对话名称(默认:互动聊天)：")
	input, _ := reader.ReadString('\n')
	if len(input) == 1 {
		input = "互动聊天"
	} else {
		input = input[:len(input)-1]
	}

	targetConversationTitle := input

	var conversation *pandora.Conversation
	for _, v := range conversations {
		if v.Title == targetConversationTitle {
			conversation = v
			break
		}
	}

	fmt.Print("请输入您的第一句话(默认:你好)：")
	input, _ = reader.ReadString('\n')
	if len(input) == 0 {
		input = "你好"
	}

	var result *pandora.ConversationPostResult

	if conversation == nil {
		conversation, result, err = client.NewConversation(targetConversationTitle, input[:len(input)-1])
		if conversation == nil || err != nil {
			return
		}
		fmt.Println("ChatGPT：", result.Message.Content.Parts[0])
		time.Sleep(20 * time.Second)
	} else {
		conversation, err = client.ConversationGET(conversation.ID)
		if err != nil {
			return
		}
	}

	conversationID := conversation.ID
	parentMessageID := conversation.CurrentNode

	for {
		fmt.Print("你：")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		result, err = client.ConversationPostFinalResult(conversationID, parentMessageID, input)
		if err != nil {
			fmt.Println("ChatGPT调用失败:", err)
			return
		}

		parentMessageID = result.Message.ID
		fmt.Println("ChatGPT：", result.Message.Content.Parts[0])
	}
}
