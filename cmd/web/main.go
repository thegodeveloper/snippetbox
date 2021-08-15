package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"snippetbox.hachiko.app/pkg/models/postgres"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *postgres.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	db_host := flag.String("host", "localhost", "PostgreSQL Hostname")
	db_port := flag.Int("port", 5432, "PostgreSQL Port")
	db_user := flag.String("user", "hachiko", "PostgreSQL Username")
	db_password := flag.String("password", "nirvana", "PostgreSQL Password")
	db_name := flag.String("dbname", "snippetbox", "PostgreSQL DB Name")
	db_sslmode := flag.String("sslmode", "disable", "PostgreSQL SSLMode")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*db_host, *db_port, *db_user, *db_password, *db_name, *db_sslmode)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &postgres.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(db_host string, db_port int, db_user string, db_password string, db_name string, db_sslmode string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db_host, db_port, db_user, db_password, db_name, db_sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
