package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"greenlight.honganhpham.net/internal/data"
	"greenlight.honganhpham.net/internal/logger"
	"greenlight.honganhpham.net/internal/mailer"
	"greenlight.honganhpham.net/internal/rate"
)

// TODO: Generate this automatically in build time
const version = "1.0.0"

type config struct {
	port      int
	env       string
	calldepth int
	db        DBConfig

	limiter rate.LimiterConfig
	smtp    mailer.MailerConfig
}

type application struct {
	debug  bool
	config config
	logger *logger.Logger
	models *data.Models
	mailer *mailer.Mailer
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
	flag.IntVar(&cfg.calldepth, "calldepth", 3, "Log level call depth")
	flag.IntVar(&cfg.limiter.RequestsPerSecond, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.BurstSize, "limiter-burst", 4, "Rate limiter maximum burst size")
	flag.IntVar(&cfg.limiter.QueueSize, "limiter-queue", 3, "Rate limiter maximum queue size")
	flag.BoolVar(&cfg.limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.StringVar(&cfg.smtp.Host, "smtp-host", os.Getenv("MAILTRAP_SMTP_HOST"), "SMTP host")
	flag.IntVar(&cfg.smtp.Port, "smtp-port", 1025, "SMTP port")
	flag.StringVar(&cfg.smtp.Username, "smtp-username", os.Getenv("MAILTRAP_SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.Password, "smtp-password", os.Getenv("MAILTRAP_SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.Sender, "smtp-sender", os.Getenv("MAILTRAP_SMTP_SENDER"), "SMTP sender")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	loggerConfig := logger.LoggerConfig{MinLevel: logger.LevelInfo, StackDepth: cfg.calldepth, ShowCaller: true}
	logger := logger.New(os.Stdout, loggerConfig)

	// f, fError := os.OpenFile("./cmd/tmp/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// if fError != nil {
	// 	log.Fatal(fError)
	// }

	// defer f.Close()

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err, nil)
	}

	defer db.Close()

	logger.Info("database connection pool establised", nil)

	app := &application{
		debug:  *debug,
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.Host, cfg.smtp.Port, cfg.smtp.Username, cfg.smtp.Password, cfg.smtp.Sender),
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err, nil)
	}
}
