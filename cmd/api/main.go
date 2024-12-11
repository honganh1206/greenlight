package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"greenlight.honganhpham.net/internal/data"
)

// TODO: Generate this automatically in build time
const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	debug  bool
	config config
	logger *logger
	models *data.Models
}

func main() {
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("GREENLIGHT_DB_DSN")
	// port := os.Getenv("APP_PORT")

	var cfg config

	// Read values from command-line flags to struct
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dsn, "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	f, fError := os.OpenFile("./cmd/tmp/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if fError != nil {
		log.Fatal(fError)
	}

	defer f.Close()

	logger := newLogger(f)

	db, err := openDB(cfg)

	if err != nil {
		logger.errorLog.Fatal(err)
	}

	defer db.Close()

	logger.infoLog.Print("database connection pool establised")

	app := &application{
		debug:  *debug,
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port), // String formatting
		Handler:      http.HandlerFunc(app.ServeHTTP),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second, // TODO: Hardcoded values here
		WriteTimeout: 30 * time.Second,
	}

	logger.infoLog.Printf("Starting %s server on port %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.errorLog.Fatal(err)
}
