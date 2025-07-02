package main

import (
	"fmt"
	"log"

	"github.com/t1nyb0x/tracktaste/httpclient"
)

func main() {
    accessToken, err := httpclient.GetAccessToken()
    if err != nil {
        log.Fatal(err)
    }
    
    body, err := httpclient.Fetch("https://api.spotify.com/v1/artists/1bY7QMGccPmba1f1frZ8Xb", accessToken)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Response:")
    fmt.Println(body)
}