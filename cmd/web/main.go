package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/shinshARK/snippetbox/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // New import
)

type application struct {
	// errorLog       *log.Logger
	// infoLog        *log.Logger
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

type config struct {
	addr      string
	staticDir string
	dsn       string
}

func main() {
	var cfg config

	flag.StringVar(&cfg.dsn, "dsn", "web:letsgo15942434@/snippetbox?parseTime=true", "MySQL data source name")
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP Network Address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	fmt.Println(cfg)

	// infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	// errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})) // newer normal logs with structured logs
	// errorLog := log.New(logging.NewSlogWrapper(logger), "", 0)                 // wrapper to support errorlogs of http.Serve

	db, err := openDB(cfg.dsn)
	if err != nil {
		// errorLog.Fatal(err)
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		// errorLog.Fatal(err)
		logger.Error(err.Error())
		os.Exit(1)

	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure that the Secure attribute is set on our session cookies.
	// Setting this means that the cookie will only be sent by a user's web
	// browser when a HTTPS connection is being used (and won't be sent over an
	// unsecure HTTP connection).
	sessionManager.Cookie.Secure = true

	app := &application{
		// errorLog:       errorLog,
		// infoLog:        infoLog,
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:      cfg.addr,
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
	}

	// infoLog.Printf("Starting server on %s", cfg.addr)
	logger.Info("Starting server", "port", cfg.addr)

	// err = server.ListenAndServe()
	err = server.ListenAndServeTLS("./tls/localhost+3.pem", "./tls/localhost+3-key.pem")
	// errorLog.Fatal(err)
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
