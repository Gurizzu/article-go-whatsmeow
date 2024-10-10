package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func WAConnect(phone string) (*whatsmeow.Client, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:database/wapp.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}

	devices, err := container.GetAllDevices()
	if err != nil {
		log.Println("Failed to get devices:", err.Error())
		return nil, err
	}

	var deviceStore *store.Device
	for _, device := range devices {
		if device.ID.User == phone {
			deviceStore = device
			break
		}
	}

	if deviceStore == nil {
		log.Println("Device not found for JID:", phone, " - starting new session")
		// if not found, create new session
		deviceStore = container.NewDevice()
	}
	fmt.Println(deviceStore.ID)

	client := whatsmeow.NewClient(deviceStore, waLog.Noop)
	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err := client.Connect()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		msg := v.Message.GetConversation()
		fmt.Printf("Received message from %s: %s\n", v.Info.Sender.User, msg)
	case *events.Receipt:
		fmt.Printf("Message status changed: %v\n", v)
	case *events.Disconnected:
		fmt.Println("Client disconnected")
	case *events.LoggedOut:
		fmt.Println("Client logged out")
	default:
		fmt.Printf("Other event: %T\n", v)
	}
}

func WAConnect1() (*whatsmeow.Client, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:database/wapp2.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Println("Failed to get devices:", err.Error())
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, waLog.Noop)
	if client.Store.ID == nil {

		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		for evt := range qrChan {

			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}

	} else {
		err := client.Connect()
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
