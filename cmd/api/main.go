package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
	"database/sql"
	_ "github.com/lib/pq"
	"sanriohub.pavelkan.net/internal/data" 
    "sanriohub.pavelkan.net/internal/mailer"
		"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const version = "1.0.0"

type config struct {
    port int
    env string
	fill bool
	migrations string
	db struct {
		dsn string
	}
	smtp struct {
		host string
		port int
		username string
		password string
		sender string
	}
}

type application struct {
    config config
    logger *log.Logger
	models data.Models
	mailer mailer.Mailer
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server host")
	flag.StringVar(&cfg.migrations, "migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db", "postgresql://doadmin:show-password@db-postgresql-fra1-12636-do-user-16680281-0.c.db.ondigitalocean.com:25060/defaultdb?sslmode=require", "PostgreSQL DSN")
    flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
    flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
    flag.StringVar(&cfg.smtp.username, "smtp-username", "7590dd4f0d4c7f", "SMTP username")
    flag.StringVar(&cfg.smtp.password, "smtp-password", "756f8e866ede5b", "SMTP password")
    flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.alexedwards.net>", "SMTP sender")
    flag.Parse()


	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	logger.Printf(cfg.migrations)
	db, err := openDB(cfg)
	if err != nil {
	    logger.Fatal(err)
	}

	defer db.Close()
	
	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	if cfg.fill {
		err = data.PopulateDatabase(app.models)
		if err != nil {
			logger.Fatal(err, nil)
			return
		}
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
    err = srv.ListenAndServe()
    logger.Fatal(err)

}


func openDB(cfg config) (*sql.DB, error) {
	
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	if cfg.migrations != "" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}
		m, err := migrate.NewWithDatabaseInstance(
			cfg.migrations,
			"postgres", driver)
		if err != nil {
			return nil, err
		}
		if err := m.Up(); err != nil {
			// Log or return the migration error
			return nil, fmt.Errorf("migration failed: %v", err)
		}
	}

	return db, nil
}