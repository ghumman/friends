package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Message struct {
	Color  string
	Time   time.Time
	Signal int
}

type ColorStats struct {
	Color         string
	TotalCount    int
	TotalSignal   int
	AverageSignal float64
}

type Aggregator struct {
	mu            sync.Mutex
	stats         map[string]*ColorStats
	windowDuration time.Duration
}

func NewAggregator(windowDuration time.Duration) *Aggregator {
	return &Aggregator{
		stats:         make(map[string]*ColorStats),
		windowDuration: windowDuration,
	}
}

func generateMessage() Message {
	colors := []string{"Red", "Yellow", "Green"}
	
	// Random signal between -27000 and 27000
	signal := rand.Intn(54001) - 27000 // 0 to 54000, then subtract 27000
	
	// Random interval between 0 and 1000ms
	randomInterval := time.Duration(rand.Intn(1000)) * time.Millisecond
	
	return Message{
		Color:  colors[rand.Intn(len(colors))],
		Time:   time.Now().Add(-randomInterval),
		Signal: signal,
	}
}

func (a *Aggregator) ProcessMessage(msg Message) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Initialize stats for this color if not exists
	if _, exists := a.stats[msg.Color]; !exists {
		a.stats[msg.Color] = &ColorStats{
			Color:       msg.Color,
			TotalCount:  0,
			TotalSignal: 0,
		}
	}
	
	// Update stats
	stats := a.stats[msg.Color]
	stats.TotalCount++
	stats.TotalSignal += msg.Signal
	stats.AverageSignal = float64(stats.TotalSignal) / float64(stats.TotalCount)
}

func (a *Aggregator) GetStats() []ColorStats {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	result := make([]ColorStats, 0, len(a.stats))
	for _, stats := range a.stats {
		result = append(result, *stats)
	}
	return result
}

func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.stats = make(map[string]*ColorStats)
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Create aggregator with 5-second window
	aggregator := NewAggregator(5 * time.Second)
	
	// Create ticker for message generation (every 10ms)
	messageTicker := time.NewTicker(10 * time.Millisecond)
	defer messageTicker.Stop()
	
	// Create ticker for aggregation reporting (every 5 seconds)
	reportTicker := time.NewTicker(5 * time.Second)
	defer reportTicker.Stop()
	
	// Channel to signal when to stop
	done := make(chan bool)
	
	// Message generation goroutine
	go func() {
		for {
			select {
			case <-messageTicker.C:
				msg := generateMessage()
				aggregator.ProcessMessage(msg)
			case <-done:
				return
			}
		}
	}()
	
	// Reporting goroutine
	go func() {
		for {
			select {
			case <-reportTicker.C:
				fmt.Println("\n" + strings.Repeat("=", 60))
				fmt.Printf("📊 AGGREGATION REPORT (Last 5 seconds)\n")
				fmt.Println(strings.Repeat("=", 60))
				
				stats := aggregator.GetStats()
				totalMessages := 0
				
				for _, stat := range stats {
					fmt.Printf("\n🎨 Color: %s\n", stat.Color)
					fmt.Printf("   📝 Total Messages: %d\n", stat.TotalCount)
					fmt.Printf("   📊 Total Signal Sum: %d\n", stat.TotalSignal)
					fmt.Printf("   📈 Average Signal Value: %.2f\n", stat.AverageSignal)
					totalMessages += stat.TotalCount
				}
				
				fmt.Printf("\n📦 Total Messages in Window: %d\n", totalMessages)
				fmt.Println(strings.Repeat("=", 60))
				
				// Reset for next window
				aggregator.Reset()
				
			case <-done:
				return
			}
		}
	}()
	
	// Run for 15 seconds to show multiple reports
	fmt.Println("Starting message aggregation system...")
	fmt.Println("Generating messages every 10ms")
	fmt.Println("Reporting every 5 seconds")
	fmt.Println("Press Enter to stop...\n")
	
	// Wait for user input or run for fixed time
	// For demo, let's run for 15 seconds
	time.Sleep(15 * time.Second)
	close(done)
	
	fmt.Println("\n✅ System stopped.")
}