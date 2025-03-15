package main

import (
	"net"
	"os"
	"sync"

	"greenlight.honganhpham.net/internal/logger"
)

var cache sync.Map // Store and retrieve values here

type Cache struct {
	listener net.Listener
	logger   *logger.Logger
	done     chan os.Signal
}

func main() {
	listener, err := net.Listen("tcp", ":6380")

	loggerConfig := logger.LoggerConfig{MinLevel: logger.LevelInfo, StackDepth: 3, ShowCaller: true}
	logger := logger.New(os.Stdout, loggerConfig)

	if err != nil {
		logger.Fatal(err, nil)
	}

	logger.Info("Listening on tcp://0.0.0.0:6380", nil)

	cache := &Cache{listener: listener, logger: logger, done: make(chan os.Signal, 1)}

	cache.listen(logger)

	// return cache
}

func (c *Cache) listen(logger *logger.Logger) {

	for {

		conn, err := c.listener.Accept()
		logger.Info("New connection", map[string]string{"connection": conn.LocalAddr().String()})
		if err != nil {
			logger.Fatal(err, nil)
			// return
		}
		go startSession(conn, logger)
	}

}
