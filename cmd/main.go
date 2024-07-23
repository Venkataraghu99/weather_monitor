package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"weather_monitor/clients"
	"weather_monitor/handlers"
	"weather_monitor/utils"
)

func main() {

	logrus.Infof("Weather monitor started")

	// initialize environment variables and http client
	env := utils.NewEnv()
	httpClient := &http.Client{}

	// initialize service dependency
	openWeatherClient := clients.NewDefaultOpenWeatherClient(httpClient, env)
	handler := handlers.NewDefaultWeatherMonitor(env, openWeatherClient)

	// create http router
	r := mux.NewRouter()
	r.HandleFunc("/weather/status", handler.GetWeatherStatus)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", env.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)

}
