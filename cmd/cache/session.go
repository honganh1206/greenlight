package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"greenlight.honganhpham.net/internal/logger"
)

// Handle the client's session
// Parse and execute commands
// Then write responses back to the client
func startSession(conn net.Conn, logger *logger.Logger) {
	// Ensure the connection will ALWAYS be closed
	defer func() {
		logger.Info("Closing connection", map[string]string{"connection": conn.LocalAddr().String()})
	}()

	// At some point we might be reading from a closed connection
	// And we do not want the server to die in case of an error
	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Errorf("Error: %s", err), nil)
		}
	}()

	p := NewParser(conn, logger)

	for {
		cmd, err := p.command(logger)
		if err != nil {
			logger.Error(fmt.Errorf("Error: %s", err), nil)
			conn.Write([]uint8("-ERR " + err.Error() + "\r\n"))
			break
		}
		// End of a session
		if !cmd.handle(logger) {
			break
		}
	}
}

func (cmd Command) handle(logger *logger.Logger) bool {
	switch strings.ToUpper(cmd.args[0]) {
	case "GET":
		return cmd.get(logger)
	// case "SET":
	// 	return cmd.set(logger)
	case "DEL":
		return cmd.del()
	case "QUIT":
		return cmd.quit(logger)
	default:
		logger.Info("Command not supported", map[string]string{"command": cmd.args[0]})
		cmd.conn.Write([]uint8("-ERR unknown command '" + cmd.args[0] + "'\r\n"))
	}
	return true
}

func (cmd *Command) quit(logger *logger.Logger) bool {
	if len(cmd.args) != 1 {
		cmd.conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.args[0] + "' command\r\n"))
		return true
	}
	logger.Info("Handle QUIT", nil)
	cmd.conn.Write([]uint8("+OK\r\n"))
	return false
}

func (cmd *Command) del() bool {
	count := 0
	for _, key := range cmd.args[1:] {
		if _, ok := cache.LoadAndDelete(key); ok {
			count++
		}
	}
	// Write back to the client the number of keys deleted
	cmd.conn.Write(fmt.Appendf(nil, ":%d\r\n", count))
	return true
}

func (cmd *Command) get(logger *logger.Logger) bool {
	if len(cmd.args) != 2 {
		cmd.conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.args[0] + "' command\r\n"))
		return true
	}
	logger.Info("Handle GET", nil)
	val, _ := cache.Load(cmd.args[1])
	if val != nil {
		res, _ := val.(string)
		if strings.HasPrefix(res, "\"") {
			res, _ = strconv.Unquote(res)
		}
		logger.Info("Response length", map[string]string{"length": strconv.Itoa(len(res))})
		cmd.conn.Write(fmt.Appendf(nil, "$%d\r\n", len(res)))
		cmd.conn.Write(append([]uint8(res), []uint8("\r\n")...)) // Write the key-value
	} else {
		cmd.conn.Write([]uint8("$-1\r\n"))
	}
	return true
}
