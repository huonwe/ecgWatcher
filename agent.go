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

	router := gin.Default()
	// 网页静态资源
	router.LoadHTMLGlob("views/*")
	router.Static("/static", "./static")
	// everyone can access
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	log.Print("running on " + utils.LocalHost() + ":8080")
	router.Run("0.0.0.0:8080")

}

var messageReceiveHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	log.Println(m.Topic())
}
