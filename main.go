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


var MicrosoftTeamId string
var MicrosoftTeamChannelId string
var MicrosoftAccessToken string

func init() {
	stop = make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	MicrosoftTeamId = os.Getenv("MICROSOFT_TEAM_ID")
	MicrosoftTeamChannelId = os.Getenv("MICROSOFT_TEAM_CHANNEL_ID")
	MicrosoftAccessToken = os.Getenv("MICROSOFT_ACCESS_TOKEN")

	fmt.Println(fmt.Sprintf("Using Microsoft Team Id: %s", MicrosoftTeamId))
	fmt.Println(fmt.Sprintf("Using Microsoft Channel Team Id: %s", MicrosoftTeamChannelId))
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
