package main

import (
	"context"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
	"database/sql"
	_ "github.com/lib/pq"
	"sanriohub.pavelkan.net/internal/data" 

)

const version = "1.0.0"

type config struct {
    port int
    env string
	db struct {
		dsn string
	}
}

type application struct {
    config config
    logger *log.Logger
	models data.Models
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
    flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
    flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:1234@localhost:5432/sanrio?sslmode=disable", "PostgreSQL DSN") 
	flag.Parse()


	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

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
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
	return nil, err
	}
	
	return db, nil
	}