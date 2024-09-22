package main

import "C"
import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/patrickmn/go-cache"
	rag_robot "go-robot/rag-robot"
	"net/http"
)

const (
	APP_VERIFICATION_TOKEN = "" // 应用验证token
	APP_ENCRYPT_KEY        = "" // 应用加密密钥
	ROBOT_NAME             = "" // 机器人名称
)

func main() {

	// 注册消息处理器
	handler := dispatcher.NewEventDispatcher(APP_VERIFICATION_TOKEN, APP_ENCRYPT_KEY).
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			// 判断是否已回复过该消息
			if value, ok := rag_robot.Cache.Get(event.EventV2Base.Header.EventID); value != nil && ok {
				return nil
			}
			// 只处理文本消息
			if *event.Event.Message.MessageType != "text" {
				rag_robot.SendErrorMsg("暂不支持非文本消息", event.Event.Message.MessageId)
				rag_robot.Cache.Set(event.EventV2Base.Header.EventID, 1, cache.NoExpiration)
				return nil
			}
			// 只处理 @机器人的消息，并且是@ 本机器人
			if *event.Event.Message.ChatType == "group" && (len(event.Event.Message.Mentions) == 0 || *event.Event.Message.Mentions[0].Name != ROBOT_NAME) {
				return nil
			}
			content := event.Event.Message.Content
			contentStr := rag_robot.ParseContent(*content)
			infer, msg, _ := rag_robot.RobotInfer(contentStr)
			if msg != "" {
				rag_robot.SendErrorMsg(msg, event.Event.Message.MessageId)
				rag_robot.Cache.Set(event.EventV2Base.Header.EventID, 1, cache.NoExpiration)
				return nil
			}
			rag_robot.SendMsg(infer, event.Event.Message.MessageId)
			// 记录已回复的消息ID
			rag_robot.Cache.Set(event.EventV2Base.Header.EventID, 1, cache.NoExpiration)
			return nil
		}).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		return nil
	})

	// 注册 http 路由
	http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler, larkevent.WithLogLevel(larkcore.LogLevelError)))

	// 启动 http 服务
	fmt.Println("http server started", "http://localhost:8080/webhook/event")

	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		panic(err2)
	}
}
