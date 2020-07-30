package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var server *http.Server
var stop chan os.Signal


var MicrosoftTeamWebhook string

func init() {
	stop = make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	MicrosoftTeamWebhook = os.Getenv("MICROSOFT_TEAM_WEBHOOK")
	fmt.Println(fmt.Sprintf("Using Microsoft Team Webhook URL: %s", MicrosoftTeamWebhook))
}

func main() {

	fmt.Println("Starting up ...")
	rand.Seed(time.Now().UnixNano())

	server = GetWebServer()
	go func() {
		fmt.Println("Starting xray ms team integration server...")
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(fmt.Sprintf("Http server stopped: %s", err))
			stop <- os.Interrupt
		}
	}()

	<-stop
	fmt.Println("Stop signal received!")
	fmt.Println("Shutting down server...")
	server.Shutdown(context.Background())
	fmt.Println("Server stopped!")
}
