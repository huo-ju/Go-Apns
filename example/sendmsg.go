package main

import (
"fmt"
"github.com/virushuo/Go-Apns"
"os"
"math/rand"
"time"
)

func main(){

    client_conn,rchan,wchan,err := goapns.Connect("apns_dev_cert.pem", "apns_dev_key.pem","gateway.sandbox.push.apple.com:2195")
    if err != nil{
        fmt.Printf("connect error: %s\n", err.Error())
        os.Exit(1)
    }

    r := rand.New(rand.NewSource(time.Now().Unix()))

    notification := goapns.Notification{}
    notification.Device_token = "ae5cb3dd7cbffc822050995779cc138cfb70a3c81a36158caf3e8fb71ce7bda1"
    notification.Alert = "hello world!"
    notification.Identifier = r.Uint32()
    fmt.Println(notification)
    wchan <- notification

    apn_error := <-rchan
    fmt.Println(apn_error)
    client_conn.Close()


}
