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
		r.Header().Add("strict-transport-security", "max-age=31536000")
		r.Header().Add("includeSubDomains", "true")
		r.Header().Add("preload", "true")
		r.HTML(200, "index", nil)
	})

	if os.Getenv("PANIC") == "true" {
		panic("this is crashing")
	}

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	go http.Serve(listener, m)
	log.Println("Listening on 0.0.0.0:" + port)

	go func() {
		// for {
		// 	log.Println("TEST_VARIABLE:", os.Getenv("TEST_VARIABLE"))
		// 	time.Sleep(10 * time.Second)
		// }
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs
	fmt.Println("SIGTERM, time to shutdown")
	listener.Close()
}
