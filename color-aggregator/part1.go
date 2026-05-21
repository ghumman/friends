package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	Color string
	Time  time.Time
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Colors to choose from
	colors := []string{"Red", "Yellow", "Green"}
	
	// Create a ticker that fires every 10 milliseconds
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	
	// Counter to stop after 10 messages
	messagesPrinted := 0
	
	fmt.Println("Starting message printer...\n")
	
	for range ticker.C {
		// Randomly select a color
		color := colors[rand.Intn(len(colors))]
		
		// Generate random interval between 0 and 1 second
		randomInterval := time.Duration(rand.Intn(1000)) * time.Millisecond
		
		// Current time minus random interval
		messageTime := time.Now().Add(-randomInterval)
		
		// Create message
		msg := Message{
			Color: color,
			Time:  messageTime,
		}
		
		// Print in human-readable format
		fmt.Printf("[%s] %s - Time: %s (%.0f ms ago)\n", 
			msg.Color,
			msg.Time.Format("15:04:05.000"),
			msg.Color,
			randomInterval.Seconds()*1000,
		)
		
		messagesPrinted++
		
		// Stop after 10 messages
		if messagesPrinted >= 10 {
			break
		}
	}
	
	fmt.Println("\nDone! Printed 10 messages.")
}