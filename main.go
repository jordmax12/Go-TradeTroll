package main

import (
	"log"
	"time"

	g "github.com/thauanvargas/goearth"
	"github.com/thauanvargas/goearth/shockwave/in"
	"github.com/thauanvargas/goearth/shockwave/out"
)

var (
	target     string
	delay      = 200 * time.Millisecond
	ext        = g.NewExt(g.ExtInfo{
		Title:       "Trade Troll",
		Author:      "TradeTroll",
		Version:     "1.0.0",
		Description: "Trolls people by initiating trades and quickly canceling them",
	})
)

func main() {
	ext.Initialized(func(e g.InitArgs) {
		log.Println("Trade Troll Extension initialized")
	})

	ext.Connected(func(e g.ConnectArgs) {
		log.Println("Connected to " + e.Host)
	})

	ext.Disconnected(func() {
		log.Println("Disconnected")
	})

	ext.Intercept(in.TRADE_OPEN).With(handleTradeOpen)

	ext.Run()
}

func handleTradeOpen(e *g.Intercept) {
	log.Println("Trade opened! Canceling in", delay, "seconds...")

	go func() {
		time.Sleep(delay)

		log.Println("Canceling trade...")
		ext.Send(out.TRADE_CLOSE)
		target = ""
	}()
}
