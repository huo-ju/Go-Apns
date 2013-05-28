package main

import (
	"fmt"
	"github.com/virushuo/Go-Apns"
	"os"
	"time"
)

func main() {
	apn, err := apns.New("apns_dev_cert.pem", "apns_dev_key.pem", "gateway.sandbox.push.apple.com:2195", 1*time.Second)
	if err != nil {
		fmt.Printf("connect error: %s\n", err.Error())
		os.Exit(1)
	}
	go readError(apn.ErrorChan)

	token := "your device token"

	payload := apns.Payload{}
	payload.Aps.Alert.Body = "hello world! 0"

	notification := apns.Notification{}
	notification.DeviceToken = token
	notification.Identifier = 0
	notification.Payload = &payload
	err = apn.Send(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)

	notification.Identifier++
	notification.Payload.Aps.Alert.Body = "hello world! 1"
	err = apn.Send(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)

	notification.Identifier++
	notification.Payload.Aps.Alert.Body = "hello world! 2"
	err = apn.Send(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)

	notification.Identifier++
	notification.DeviceToken = ""
	notification.Payload.Aps.Alert.Body = "hello world! 3"
	err = apn.Send(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)
	time.Sleep(1E9)

	notification.Identifier++
	notification.DeviceToken = token
	notification.Payload.Aps.Alert.Body = "re hello world! 0"
	err = apn.Send(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)

	notification.Identifier++
	notification.DeviceToken = ""
	notification.Payload.Aps.Alert.Body = "re hello world! 1"
	err = apn.Send(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)
	time.Sleep(1E9)

	notification.Identifier++
	notification.DeviceToken = token
	notification.Payload.Aps.Alert.Body = "rere hello world! 0"
	err = apn.Send(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)
	time.Sleep(2e9)

	notification.Identifier++
	notification.DeviceToken = token
	notification.Payload.Aps.Alert.Body = "rere hello world! 1"
	err = apn.Send(&notification)
	fmt.Printf("send id(%x): %s\n", notification.Identifier, err)
	time.Sleep(2e9)

	apn.Close()
}

func readError(errorChan <-chan error) {
	for {
		apnerror := <-errorChan
		fmt.Println(apnerror.Error())
	}
}
