package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type Level int8

const CALL_DEPTH = 3

const MAX_BUFFER_SIZE = 64 << 10

// Cache and reuse byte slices to reduce memory allocations
var bufferPool = sync.Pool{New: func() any { return new([]byte) }}

func getBuffer() *[]byte {
	p := bufferPool.Get().(*[]byte)

	// Reset the buffer while preserving capacity
	*p = (*p)[:0]
	return p
}

// Place the buffer back into the pool
func putBuffer(p *[]byte) {
	// Set a hard-coded limit for buffers returning to the pool
	// If buffer size exceeds the limit, we let the garbage collector reclaim the memory
	if cap(*p) > MAX_BUFFER_SIZE {
		*p = nil
	}

	bufferPool.Put(p)
}

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

// TODO: Upgrade logger with MultiWriter for writing logs to files?
// Check history Level String Methid Implementation
type Logger struct {
	out    io.Writer
	config LoggerConfig
	mu     sync.Mutex // coordinate the writes
}

type LoggerConfig struct {
	MinLevel   Level
	StackDepth int
	ShowCaller bool // Optional to show caller info
}

func New(out io.Writer, cfg LoggerConfig) *Logger {
	return &Logger{
		out:    out,
		config: cfg,
	}
}

type CallerInfo struct {
	File     string `json:"file,omitempty"`
	Function string `json:"function,omitempty"`
	Line     int    `json:"line,omitempty"`
}

func getCaller(calldepth int) *CallerInfo {
	// Return the memory address pointing to the function code located in memory
	pc, file, line, ok := runtime.Caller(calldepth)

	if !ok {
		return nil
	}

	fn := runtime.FuncForPC(pc)

	if fn == nil {
		return nil
	}

	return &CallerInfo{
		File:     filepath.Base(file),
		Function: filepath.Base(fn.Name()),
		Line:     line,
	}
}

// Implicitly implement the Stringer interface here
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

func (l *Logger) Output(msg []byte) (n int, err error) {
	return l.output(LevelError, string(msg), nil)
}

func (l *Logger) Info(msg string, props map[string]string) {
	l.output(LevelInfo, msg, props)
}

func (l *Logger) Error(err error, props map[string]string) {
	l.output(LevelError, err.Error(), props)
}

func (l *Logger) Fatal(err error, props map[string]string) {
	l.output(LevelFatal, err.Error(), props)
	os.Exit(1) // Terminate the app
}

func (l *Logger) output(level Level, msg string, props map[string]string) (int, error) {
	// No need to display level below error
	if level < l.config.MinLevel {
		return 0, nil
	}

	buf := getBuffer()
	defer putBuffer(buf)

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties"`
		Trace      string            `json:"trace,omitempty"`
		// Pointer type is optional here, but it would be useful when caller info is not available or there is an error
		Caller *CallerInfo `json:"caller,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    msg,
		Properties: props,
	}
	// Immediate caller info - Where the log was called from
	if l.config.ShowCaller {
		// Skip runtime.Caller + print + PrintInfo/PrintError
		aux.Caller = getCaller(CALL_DEPTH)
	}

	// Detailed stack trace
	if level >= LevelError {
		if l.config.StackDepth > 0 {
			stack := make([]uintptr, l.config.StackDepth)
			// CALL_DEPTH and l.config.StackDepth serve different purposes
			length := runtime.Callers(CALL_DEPTH, stack[:]) // Createa a slice from an array
			if length > 0 {
				// Get PC values from Callers and return function/file/line information
				frames := runtime.CallersFrames(stack[:length])

				var trace strings.Builder

				for {
					frame, more := frames.Next()
					fmt.Fprintf(&trace, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
					if !more {
						break
					}
				}

				aux.Trace = trace.String()
			}

		} else {
			// If no depth is specified
			aux.Trace = string(debug.Stack())
		}
	}

	jsonData, err := json.Marshal(aux)

	// TODO: Prettify? the jsonData
	*buf = append(*buf, jsonData...)
	*buf = append(*buf, '\n') // Ensure newline

	// Single atomic write
	l.mu.Lock()
	n, err := l.out.Write(*buf)
	l.mu.Unlock()

	return n, err
}
