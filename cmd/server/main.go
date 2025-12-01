package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/batt0s/micho/internal/api"
)

func main() {
	api := api.API{}
	api.Init()

	shutdown := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)

		<-sigint

		log.Println("Interrupt signal recieved.")
		err := api.Server.Shutdown(context.Background())
		if err != nil {
			log.Println("HTTP Server Shutdown Error: ", err.Error())
		}
		log.Println("Stopped.")
		close(shutdown)
	}()

	api.Server.ListenAndServe()

	<-shutdown

}
