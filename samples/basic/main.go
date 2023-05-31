package main

import (
	"fmt"
	pandora "github.com/wk-go/pandora_go"
	"log"
	"os"
)

func main() {
	token := os.Getenv("CHAT_TOKEN")
	url := os.Getenv("CHAT_URL")
	if len(token) == 0 || len(url) == 0 {
		panic("请输入必要参数")
	}

	client := pandora.NewClient(url, token)
	//client.SetProxy("http://127.0.0.1:8888")

	data, err := client.ModelsGET()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)

	conversations, err := client.ConversationsGET()
	if err != nil {
		return
	}
	targetConversationTitle := "英语词典"

	var conversation *pandora.Conversation
	for _, v := range conversations {
		if v.Title == targetConversationTitle {
			conversation = v
			break
		}
	}
	if conversation == nil {
		conversation, _, err = client.NewConversation(targetConversationTitle, "我想让你充当英英词典，对于给出的英文单词，你要给出其中文意思以及英文解释，并且给出一个例句，此外不要有其他反馈，第一个单词是\"Hello\"")
		if conversation == nil || err != nil {
			return
		}
		//time.Sleep(20 * time.Second)
	} else {
		conversation, err = client.ConversationGET(conversation.ID)
		if err != nil {
			return
		}
	}

	wordList := []string{"conversation", "message", "result", "hello", "word"}
	for _, word := range wordList {
		result, err := client.ConversationPostFinalResult(conversation.ID, conversation.CurrentNode, word)
		if err != nil {
			return
		}
		log.Println("@@@@@@@@@@@@")
		fmt.Println(result.Message.Content.Parts[0])
		log.Println("############")
		conversation.CurrentNode = result.Message.ID
		//time.Sleep(20 * time.Second)
	}
}
