package main

import (
	"app/cloudflare"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func main() {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := new(http.Client)
	client.Jar = cookieJar
	client.Timeout = 10 * time.Second

	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36" // user ua
	cfClient := cloudflare.CreateClient(client)
	cfClient.CreatePayload(userAgent)
	cfClient.Solve()
}
