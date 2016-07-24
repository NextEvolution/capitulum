package main

import (
	"github.com/nats-io/nats"
	"fmt"
	"time"
)

func main (){
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to NATS: %s", err))
	}
	defer nc.Close()

	resp, err := nc.Request("dataservice.get.1234.lastscan", []byte{}, time.Second)
	fmt.Printf("Body recieved: %s", string(resp.Data))
}
