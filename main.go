package main

import (
	"log"
	"strings"
	"time"

	g "github.com/thauanvargas/goearth"
	"github.com/thauanvargas/goearth/shockwave/in"
	"github.com/thauanvargas/goearth/shockwave/out"
)

// Global variables for the extension
var (
	ext        *g.Ext
	target     string
)

func main() {
	// Create extension info
	info := g.ExtInfo{
		Title:       "Trade Troll",
		Author:      "TradeTroll",
		Version:     "1.0.0",
		Description: "Trolls people by initiating trades and quickly canceling them",
	}

	// Create the extension
	ext = g.NewExt(info)

	// Handle initialization
	ext.Initialized(func(info g.InitArgs) {
		log.Println("Trade Troll Extension initialized")
	})

	// Handle when connected to Habbo
	ext.Connected(func(e g.ConnectArgs) {
		log.Println("Connected to " + e.Host)
	})

	// Handle chat messages for commands
	ext.Intercept(in.CHAT).With(handleChat)

	// Handle trade responses
	ext.Intercept(in.TRADE_OPEN).With(handleTradeOpen)
	ext.Intercept(in.TRADE_ACCEPT).With(handleTradeAccept)

	// Run the extension
	ext.Run()
}

// handleChat processes chat messages for commands
func handleChat(e *g.Intercept) {
	message := e.Packet.ReadString()
	log.Println("Message:", message)

	// Look for ":tradet " pattern anywhere in the message
	tradeCmdIndex := strings.Index(strings.ToLower(message), ":tradet ")
	if tradeCmdIndex == -1 {
		return // No :tradet command found
	}

	// Extract the part starting from ":tradet "
	cmdPart := message[tradeCmdIndex:]
	log.Println("Command part:", cmdPart)

	// Split the command part to get the username
	parts := strings.Fields(cmdPart)
	if len(parts) < 2 {
		log.Println("Usage: :tradet <username>")
		return
	}

	username := parts[1]
	log.Println("Target username:", username)
	startTrade(username)

	e.Block()
}

// startTrade starts a trade with a user
func startTrade(username string) {
	log.Println("Starting trade with:", username)
	target = username

	// Send trade request
	ext.Send(out.TRADE_OPEN, username)
}

// handleTradeOpen handles when a trade is opened
func handleTradeOpen(e *g.Intercept) {
	log.Println("Trade opened! Canceling in 0.5 seconds...")

	// Start a timer to cancel the trade after 0.5 seconds
	go func() {
		time.Sleep(200 * time.Millisecond)

		log.Println("Canceling trade...")
		ext.Send(out.TRADE_CLOSE)
		target = ""
	}()
}

// handleTradeAccept handles when the other person accepts the trade
func handleTradeAccept(e *g.Intercept) {
	log.Println("They accepted the trade! Canceling immediately...")

	// Cancel the trade immediately if they accept
	ext.Send(out.TRADE_CLOSE)
}
