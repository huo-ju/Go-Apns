package main

import (
	"fmt"
	"github.com/virushuo/Go-Apns"
	"os"
	"time"
)

func main() {
	apn, err := goapns.Connect("apns_dev_cert.pem", "apns_dev_key.pem", "gateway.sandbox.push.apple.com:2195")
	if err != nil {
		fmt.Printf("connect error: %s\n", err.Error())
		os.Exit(1)
	}
	go readError(apn.ErrorChan)

	token := "your device token"

	payload := goapns.Payload{}
	payload.Aps.Alert = "hello world! 0"

	notification := goapns.Notification{}
	notification.DeviceToken = token
	notification.Identifier = 0
	notification.Payload = &payload
	err = apn.SendNotification(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)

	notification.Identifier++
	notification.Payload.Aps.Alert = "hello world! 1"
	err = apn.SendNotification(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)

	notification.Identifier++
	notification.Payload.Aps.Alert = "hello world! 2"
	err = apn.SendNotification(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)

	notification.Identifier++
	notification.DeviceToken = ""
	notification.Payload.Aps.Alert = "hello world! 3"
	err = apn.SendNotification(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)
	time.Sleep(1E9)

	err = apn.Reconnect()
	fmt.Printf("reconnect: %s\n", err)

	notification.Identifier++
	notification.DeviceToken = token
	notification.Payload.Aps.Alert = "re hello world! 0"
	err = apn.SendNotification(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)

	notification.Identifier++
	notification.DeviceToken = ""
	notification.Payload.Aps.Alert = "re hello world! 1"
	err = apn.SendNotification(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)
	time.Sleep(1E9)

	err = apn.Reconnect()
	fmt.Printf("reconnect: %s\n", err)

	notification.Identifier++
	notification.DeviceToken = token
	notification.Payload.Aps.Alert = "rere hello world! 0"
	err = apn.SendNotification(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)

	apn.Close()
	time.Sleep(1E9)
}

func readError(errorChan <-chan goapns.NotificationError) {
	for {
		apnerror := <-errorChan
		fmt.Println(apnerror.Error())
	}
}
