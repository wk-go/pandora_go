package pandora_go

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

type Client struct {
	UrlPrefix            string // http(s)://ip:port/api
	Token                string
	Proxy                *url.URL // protocol://user:pwd@ip:port
	Model                string
	Headers              map[string]string
	LastConversationTime time.Time //最后一次发送会话时间
	LoginResponse        *LoginResponse
}

func NewClientLogin(urlPrefix, Username, password string) (*Client, error) {
	client := NewClient(urlPrefix, "")
	if err := client.Login(Username, password, ""); err != nil {
		return nil, err
	}
	return client, nil
}

func NewClient(urlPrefix, token string) *Client {
	return &Client{
		UrlPrefix: urlPrefix,
		Token:     token,
		Model:     "text-davinci-002-render-sha",
	}
}

func (c *Client) SetProxy(proxy string) (bool, error) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return false, err
	}
	c.Proxy = proxyUrl
	return true, nil
}

func (c *Client) AddHeader(key, value string) {
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
	c.Headers[key] = value
}

// Login 登录
func (c *Client) Login(username, password, mfaCode string) error {
	_url := c.UrlPrefix + "/auth/login"

	postData := url.Values{
		"username": {username},
		"password": {password},
		"mfa_code": {mfaCode},
	}

	body, err := c.Request("POST", _url, "", []byte(postData.Encode()),
		"Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	if err != nil {
		return err
	}

	loginResp := new(LoginResponse)
	err = json.Unmarshal(body, loginResp)
	if err != nil {
		return err
	}

	if len(loginResp.AccessToken) == 0 {
		if _err, err := NewErrorResponse(body); err != nil {
			return err
		} else {
			return _err
		}
	}

	c.LoginResponse = loginResp
	c.Token = loginResp.AccessToken

	return nil
}

// ModelsGET 模型列表
func (c *Client) ModelsGET() (*ModelList, error) {
	_url := c.UrlPrefix + "/models"
	body, err := c.Request("GET", _url, c.Token, nil)
	if err != nil {
		return nil, err
	}

	modelList := new(ModelList)
	if err = json.Unmarshal(body, modelList); err != nil {
		return nil, err
	}

	return modelList, nil
}

// ConversationsGET 对话列表
func (c *Client) ConversationsGET() ([]*Conversation, error) {
	_url := c.UrlPrefix + "/conversations"
	body, err := c.Request("GET", _url, c.Token, nil)
	if err != nil {
		return nil, err
	}
	conversationListResult := new(ConversationListResult)
	err = json.Unmarshal(body, conversationListResult)
	if err != nil {
		return nil, err
	}
	return conversationListResult.Items, nil
}

// NewConversation 创建新的会话
func (c *Client) NewConversation(title, prompt string) (conversation *Conversation, result *ConversationPostResult, err error) {
	// 创建一个新的并重命名
	parentMessageID := uuid.New().String()
	result, err = c.ConversationPostFinalResult("", parentMessageID, prompt)
	if err != nil {
		return nil, nil, err
	}

	conversations, err := c.ConversationsGET()
	if err != nil {
		return nil, result, err
	}
	for _, v := range conversations {
		if strings.Contains(v.ID, result.ConversationID) {
			conversation = v
			conversation.Title = title
			conversation.CurrentNode = result.Message.ID
			break
		}
	}

	if conversation == nil {
		return nil, result, errors.New("create fail")
	}
	if ok, err := c.ChangeConversationTitle(conversation.ID, title); err != nil || !ok {
		return nil, result, err
	}
	return conversation, result, nil
}

// ConversationGET 对话详情
func (c *Client) ConversationGET(id string) (*Conversation, error) {
	_url := c.UrlPrefix + "/conversation/" + id
	body, err := c.Request("GET", _url, c.Token, nil)
	if err != nil {
		return nil, err
	}
	conversation := new(Conversation)
	err = json.Unmarshal(body, conversation)
	if err != nil {
		return nil, err
	}
	conversation.ID = id
	return conversation, nil
}

// ChangeConversationTitle 修改会话title
func (c *Client) ChangeConversationTitle(id, title string) (bool, error) {
	data := map[string]interface{}{
		"title": title,
	}
	return c.ConversationPATCH(id, data)
}

// DeleteConversation 删除会话
func (c *Client) DeleteConversation(id string) (bool, error) {
	data := map[string]interface{}{
		"is_visible": false,
	}
	return c.ConversationPATCH(id, data)
}

// ConversationPATCH 修改对话信息
func (c *Client) ConversationPATCH(id string, data interface{}) (bool, error) {
	_url := c.UrlPrefix + "/conversation/" + id

	content, _ := json.Marshal(data)

	_, err := c.Request("PATCH", _url, c.Token, content)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConversationPostFinalResult 发起一次聊天并返回最后的结果
func (c *Client) ConversationPostFinalResult(conversationID, parentMessageID, prompt string) (*ConversationPostResult, error) {
	bodySlice, err := c.ConversationPOST(conversationID, parentMessageID, prompt)
	if err != nil {
		return nil, err
	}

	if len(bodySlice) < 3 {
		errResp, err := NewErrorResponse(bodySlice[0])
		if err != nil {
			return nil, err
		}
		return nil, errResp
	}

	resultOffset := 3
	if len(bodySlice[len(bodySlice)-1]) == 0 {
		resultOffset = 4
	}
	resultSlice := bytes.Replace(bodySlice[len(bodySlice)-resultOffset], []byte("data: "), nil, -1)
	conversationPostResult := new(ConversationPostResult)
	err = json.Unmarshal(resultSlice, conversationPostResult)
	if err != nil {
		return nil, err
	}
	return conversationPostResult, nil
}

// ConversationPostListResult 发起一次聊天并返回所有的结果
func (c *Client) ConversationPostListResult(conversationID, parentMessageID, prompt string) ([]*ConversationPostResult, error) {

	bodySlice, err := c.ConversationPOST(conversationID, parentMessageID, prompt)
	if err != nil {
		return nil, err
	}

	if len(bodySlice) < 3 {
		errResp, err := NewErrorResponse(bodySlice[0])
		if err != nil {
			return nil, err
		}
		return nil, errResp
	}

	resultSlice := make([]*ConversationPostResult, 0, len(bodySlice))
	for _, s := range bodySlice {
		if len(s) < 6 {
			continue
		}
		tmp := new(ConversationPostResult)
		err = json.Unmarshal(s[6:], tmp)
		if err != nil {
			continue
		}
		resultSlice = append(resultSlice, tmp)
	}
	return resultSlice, nil
}

// ConversationPOST 发起一次聊天
func (c *Client) ConversationPOST(conversationID, parentMessageID, prompt string) ([][]byte, error) {
	_url := c.UrlPrefix + "/conversation"
	messageID := uuid.New()
	var content []byte
	message := &MessageRequest{
		Action: "next", // variant重新生成
		Messages: []Message{
			{
				ID:     messageID.String(),
				Author: MessageAuthor{Role: "user"},
				Role:   "user",
				Content: MessageContent{
					ContentType: "text",
					Parts:       []string{prompt},
				},
			},
		},
		ConversationID:  nil,
		ParentMessageID: parentMessageID,
		Model:           c.Model,
		TimezoneOffset:  -480,
	}

	if len(conversationID) > 0 {
		message.ConversationID = &conversationID
	}

	content, _ = json.Marshal(message)

	c.LastConversationTime = time.Now()
	body, err := c.Request("POST", _url, c.Token, content)
	if err != nil {
		return nil, err
	}
	bodySlice := bytes.Split(body, []byte("\n\n"))
	return bodySlice, err
}

// Request 发起请求
func (c *Client) Request(method, url, token string, content []byte, headers ...string) ([]byte, error) {
	var req *http.Request
	var err error
	req, err = http.NewRequest(method, url, bytes.NewBuffer(content))

	if len(token) > 0 {
		req.Header.Add("authorization", "Bearer "+token)
	}
	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}
	if len(headers) > 0 && len(headers)%2 == 0 {
		key := ""
		for k, v := range headers {
			if k%2 == 0 {
				key = v
			} else {
				req.Header.Add(key, v)
			}
		}
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := c.RequestDo(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}

func (c *Client) RequestDo(req *http.Request) (*http.Response, error) {

	client := &http.Client{}

	if c.Proxy != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(c.Proxy),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}
