package rag_robot

import (
	"context"
	"encoding/json"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go-robot/dao"
	"go-robot/httplib"
	"log"
	"regexp"
)

const (
	APP_ID          = ""
	APP_SECRET      = ""
	WEI_ZAI_API_KEY = ""
	HOST            = ""
)

var Client = lark.NewClient(APP_ID, APP_SECRET)

func RobotInfer(msg string) (dao.ChatResp, string, error) {
	conversation, err := CreatConversation()
	if err != nil {
		log.Default().Println("创建会话出错了", err)
		return dao.ChatResp{}, "出错了，请重试～", err
	}
	chatResp, err := ChatLLm(conversation, msg)
	if err != nil {
		log.Default().Println("发送消息出错了", err)
		return dao.ChatResp{}, "出错了，请重试～", err
	}
	return chatResp, "", nil
}

// 创建会话
func CreatConversation() (string, error) {
	var resp dao.CreatConvResp
	get := httplib.Get(HOST+"/api/new_conversation").Header("Authorization", WEI_ZAI_API_KEY).Retries(3)
	err := get.ToJSON(&resp)
	if err != nil {
		return "", err
	}
	return resp.Data.ID, nil
}

// 发送消息
func ChatLLm(convId, msg string) (dao.ChatResp, error) {
	req := dao.ChatReq{
		ConversationID: convId,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: msg,
			},
		},
		Stream: false,
		Quote:  true,
	}

	var resp dao.ChatResp

	// 创建HTTP请求
	postReq := httplib.Post(HOST+"/api/completion").Header("Authorization", WEI_ZAI_API_KEY).Retries(3)

	// 设置请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return resp, fmt.Errorf("failed to marshal request body: %w", err)
	}
	postReq.Body(reqBody)

	// 发送请求并解析响应
	_ = postReq.ToJSON(&resp)

	return resp, nil
}

type MsgData struct {
	Schema string `json:"schema"`
	Body   struct {
		Elements []MsgElement `json:"elements"`
	} `json:"body"`
}

type MsgElement struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

func SendMsg(chatResp dao.ChatResp, messageId *string) {
	// 构造消息
	data := DealMessageElements(chatResp)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Default().Println("Error marshaling JSON:", err)
		return
	}
	// 创建 ReplyMessageReqBuilder 实例
	replyReq := larkim.NewReplyMessageReqBuilder().
		MessageId(*messageId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			Content(string(jsonData)).
			Build()).
		Build()

	resp, err := Client.Im.Message.Reply(context.Background(), replyReq)

	// 处理错误
	if err != nil {
		log.Default().Println(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		log.Default().Println("send msg error", resp.Code, resp.Msg, resp.RequestId())
		return
	}
}

func SendErrorMsg(msg string, messageId *string) {
	// 构造消息
	data := MsgData{
		Schema: "2.0",
		Body: struct {
			Elements []MsgElement `json:"elements"`
		}{
			Elements: []MsgElement{
				{Tag: "text", Content: msg},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Default().Println("Error marshaling JSON:", err)
		return
	}
	// 创建 ReplyMessageReqBuilder 实例
	replyReq := larkim.NewReplyMessageReqBuilder().
		MessageId(*messageId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			Content(string(jsonData)).
			Build()).
		Build()

	resp, err := Client.Im.Message.Reply(context.Background(), replyReq)

	// 处理错误
	if err != nil {
		log.Default().Println(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		log.Default().Println("send msg error", resp.Code, resp.Msg, resp.RequestId())
		return
	}
}

func msgFilter(msg string) string {
	regex := regexp.MustCompile(`@[^ ]*`)
	return regex.ReplaceAllString(msg, "")

}

func ParseContent(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	text := contentMap["text"].(string)
	return msgFilter(text)
}

func DealMessageElements(resp dao.ChatResp) MsgData {
	data := MsgData{
		Schema: "2.0",
	}

	var elements []MsgElement
	elements = append(elements, MsgElement{Tag: "markdown", Content: resp.Data.Answer})
	elements = append(elements, MsgElement{Tag: "markdown", Content: "\n ---\n"})
	urlData := ""
	for _, item := range resp.Data.Reference.DocAggs {
		urlData += fmt.Sprintf("[%s](%s)\n", item.DocName, HOST+"/document/get/"+item.DocID)
		//elements = append(elements, MsgElement{Tag: "a", Content: "[" + item.DocName + "](" + HOST + "/document/get/" + item.DocID + ")"})
	}
	elements = append(elements, MsgElement{Tag: "markdown", Content: urlData})
	data.Body.Elements = elements
	return data
}
