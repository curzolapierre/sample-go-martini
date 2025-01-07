package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func isPrime(value int) bool {
	for i := 2; i <= value/2; i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(
		render.Options{
			Directory: "templates",
		},
	))

	m.Get("/ping", func(r render.Render, _ *http.Request) {
		r.JSON(200, map[string]interface{}{"message": "pong"})
		fmt.Printf("[%v] pong\n", os.Getenv("PORT"))
	})

	m.Get("/", func(r render.Render, req *http.Request) {
		if req.URL.Query().Get("wait") != "" {
			sleep, _ := strconv.Atoi(req.URL.Query().Get("wait"))
			log.Printf("Sleep for %d seconds\n", sleep)
			time.Sleep(time.Duration(sleep) * time.Second)
		}
		if req.URL.Query().Get("prime") != "" {
			val, _ := strconv.Atoi(req.URL.Query().Get("prime"))
			log.Printf("Is %d prime: %t", val, isPrime(val))
		}
		r.HTML(200, "index", nil)
	})

	if os.Getenv("PANIC") == "true" {
		panic("this is crashing")
	}

	if os.Getenv("SLOW_START") != "" {
		startTimeout, _ := strconv.Atoi(os.Getenv("SLOW_START"))

		log.Printf("Sleeping for %v seconds to simulate slow start\n", startTimeout)
		time.Sleep(time.Duration(startTimeout) * time.Second)
	}

	port := "3001"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if len(os.Args) > 1 {
		fmt.Printf("Using port from argument: %s\n", os.Args[1])
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	go http.Serve(listener, m)
	log.Println("Listening on 0.0.0.0:" + port)

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs
	fmt.Println("SIGTERM, time to shutdown")
	listener.Close()
}
