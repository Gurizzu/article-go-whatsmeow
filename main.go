package main

import (
	"context"
	"fmt"
	"go-wa/config"
	"log"

	"github.com/subosito/gotenv"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func init() {
	env := ".env"
	if e := gotenv.Load(env); e != nil {
		gotenv.Load()
	}
	config.Init()
	if config.DEBUG_MODE {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
	}
}

func main() {
	wac, err := WAConnect("6281234567890")
	if err != nil {
		fmt.Println(err)
		return
	}

	messageText := "Hi again! Just checking in to see if you received my last message. Let me know if you have any questions!"
	msg := &messageText
	wac.SendMessage(context.Background(), types.JID{
		User:   "628111111111111",
		Server: types.DefaultUserServer,
	}, &waE2E.Message{
		Conversation: msg,
	})

	defer wac.Disconnect()
}
