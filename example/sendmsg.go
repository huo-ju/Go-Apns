package main

import (
	"fmt"
	"github.com/virushuo/Go-Apns"
	"math/rand"
	"os"
	"time"
)

func main() {
	apn, err := goapns.Connect("apns_dev_cert.pem", "apns_dev_key.pem", "gateway.sandbox.push.apple.com:2195")
	if err != nil {
		fmt.Printf("connect error: %s\n", err.Error())
		os.Exit(1)
	}
	go readError(apn.Errorchan)

	r := rand.New(rand.NewSource(time.Now().Unix()))

	context := goapns.Context{}
	context.Aps.Alert = "hello world! 0"

	notification := goapns.Notification{}
	notification.DeviceToken = "YOUR_DEVICE_TOKEN"
	notification.Identifier = r.Uint32()
	notification.Context = &context
	err = apn.SendNotification(&notification)
	fmt.Println(err)

	notification.Context.Aps.Alert = "hello world! 1"
	err = apn.SendNotification(&notification)
	fmt.Println(err)

	notification.Context.Aps.Alert = "hello world! 2"
	err = apn.SendNotification(&notification)
	fmt.Println(err)

	notification.Context.Aps.Alert = "hello world! 3"
	err = apn.SendNotification(&notification)
	fmt.Println(err)
	time.Sleep(5E9)
}

func readError(Errorchan chan goapns.NotificationError) {
	apnerror := <-Errorchan
	fmt.Println(apnerror)
}
