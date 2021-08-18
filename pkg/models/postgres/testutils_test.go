package postgres

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the `multiStatements=true` parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec() call.
	dbHost := flag.String("host", "localhost", "PostgreSQL Hostname")
	dbPort := flag.Int("port", 5432, "PostgreSQL Port")
	dbUser := flag.String("user", "test_web", "PostgreSQL Username")
	dbPassword := flag.String("password", "pass", "PostgreSQL Password")
	dbName := flag.String("dbname", "test_snippetbox", "PostgreSQL DB Name")
	dbSSLMode := flag.String("sslmode", "disable", "PostgreSQL SSLMode")

	db, err := openDB(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName, *dbSSLMode)

	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from file and execute the statements.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Return the connection pool and an anonymous function which reads and
	// executes the teardown script, and closes the connection pool. We can
	// assign this anonymous function and call it later once our test has
	// completed.
	return db, func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}

func openDB(dbHost string, dbPort int, dbUser string, dbPassword string, dbName string, dbSSLMode string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
