package main

import (
	"ecgAgent/utils"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

func main() {
	topic := "topic/ecg"
	broker := "120.55.170.139:1883"
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + broker)
	opts.SetClientID("ecgAgent")
	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Connected to ", broker)
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Println("Connection lost.\n", err.Error())
	}

	client := mqtt.NewClient(opts)
	defer client.Disconnect(100)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	token := client.Subscribe(topic, 1, messageReceiveHandler)
	token.Wait()
	log.Println("Subscribed to ", topic)

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "READY")
	})
	log.Print("running on " + utils.LocalHost() + ":8080")
	r.Run()

}

var messageReceiveHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	log.Println(m.Topic())
}
