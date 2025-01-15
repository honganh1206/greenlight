package rate

import "time"

// Handle basic and bursty rate limiting
type Limiter struct {
	requests      chan int
	limiter       <-chan time.Time // (Receive-only channel) Process 1 request each time
	burstyLimiter chan time.Time   // For short bursts of requests while preserving the overall rate limit
}

type LimiterConfig struct {
	RequestsPerSecond int
	BurstSize         int
	QueueSize         int
	Enabled           bool
}

func New(cfg LimiterConfig) *Limiter {
	interval := time.Second / time.Duration(cfg.RequestsPerSecond)

	rl := &Limiter{
		requests:      make(chan int, cfg.QueueSize),
		limiter:       time.Tick(interval),
		burstyLimiter: make(chan time.Time, cfg.BurstSize),
	}

	// Init burst capacity
	for i := 0; i < cfg.BurstSize; i++ {
		rl.burstyLimiter <- time.Now()
	}

	// Replenish the tokens
	go func() {
		for t := range time.Tick(interval) {
			rl.burstyLimiter <- t
		}
	}()

	return rl
}

func (rl *Limiter) Allow() bool {
	// Avoid blocking by return immediately when a token is available
	select {
	case <-rl.burstyLimiter:
		return true
	default:
		return false
	}
}
