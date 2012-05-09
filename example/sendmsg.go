package main

import (
"fmt"
github.com/virushuo/Go-Apns
"os"
"math/rand"
"time"
)

func main(){

    apn,err := goapns.Connect("apns_dev_cert.pem", "apns_dev_key.pem","gateway.sandbox.push.apple.com:2195")
    if err != nil{
        fmt.Printf("connect error: %s\n", err.Error())
        os.Exit(1)
    }

    r := rand.New(rand.NewSource(time.Now().Unix()))

    notification := goapns.Notification{}
    notification.Device_token = "YOUR_DEVICE_TOKEN"
    notification.Alert = "hello world! 0"
    notification.Identifier = r.Uint32()
    err =apn.SendNotification(&notification)
    fmt.Println(err)

    notification.Alert = "hello world! 1"
    err =apn.SendNotification(&notification)
    fmt.Println(err)

    notification.Alert = "hello world! 2"
    err =apn.SendNotification(&notification)
    fmt.Println(err)

    notification.Alert = "hello world! 3"
    err =apn.SendNotification(&notification)
    fmt.Println(err)
    apn.Conn.Close()


}
