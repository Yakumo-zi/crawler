package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	allocatorContext, _ := chromedp.NewRemoteAllocator(context.Background(), "ws://127.0.0.1:9222")
	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	var example string

	err := chromedp.Run(ctx, chromedp.Navigate("https://pkg.go.dev/time"),
		chromedp.WaitVisible(`body > footer`),
		chromedp.Click(`#example-After`, chromedp.NodeVisible),
		chromedp.Value(`#example-After textarea`, &example))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example :\n%s", example)
}
