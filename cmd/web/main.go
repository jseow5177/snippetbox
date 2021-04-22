package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	"github.com/jseow5177/snippetbox/pkg/models/mysql"
)

type contextKey string
const contextKeyIsAuthenticated = contextKey("isAuthenticated")

// Application-wide configuration
type config struct {
	Addr string
	StaticDir string
}

// Define an application struct to hold application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	config *config
	snippets *mysql.SnippetModel
	users *mysql.UserModel
	templateCache map[string]*template.Template
	session *sessions.Session
}

func main() {

	// ========== Create custom info and error loggers ========== //

	// Use log.New() to create a custom logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string prefix for message
	// (INFO followed by a tab), and flags to indicate what additional information to include (local date and time).
	// Note that the flags are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as the destination and use
	// the log.Lshortfile flag to include the relevant file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// ========== Load env variables ========== //
	err := godotenv.Load()
	if err != nil {
		errorLog.Fatal("Error loading .env file")
	}

	pwd := os.Getenv("MYSQL_PASSWORD")
	secret := os.Getenv("SECRET")

	// ========== Parse runtime configuration settings for the application ========== //

	// Initialize application wide configuration
	cfg := new(config)

	// Define a new command-line flag with the name "addr", a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the flag
	// will be stored in the addr variable at runtime.
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")

	// Define a command-line flag for path to static directory.
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")

	// Define a command-line flag for MySQL DSN string.
	// DSN string for the driver has the format of username:password@protocol(address)/dbname?param=value
	// Default value of protocol is 'tcp'.
	// Default value of address is 'localhost:3306'.
	// parseTime param changes the output type of DATE and DATETIME values to Go's time.Time
	dsn := flag.String("dsn", fmt.Sprintf("web:%s@tcp(localhost:3306)/snippetbox?parseTime=true", pwd), "MySQL data source name")

	// We use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr variable.
	// You need to call this *before* you use the addr variable otherwise it will always
	// contain the default value of ":4000". If any errors are encountered during parsing
	// the application will be terminated.
	flag.Parse()

	// ========== Connect to MySQL DB ========== //
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Defer a call to db.Close(), so the connection pool is closed before the main() function exits.
	defer db.Close();

	// ========== Create template cache ========== //
	tc, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// ========== Initialize a new session and save into app dependency ========== //
	// Initialize a new session manager, pass in the secret key as the parameter.
	// sessions.New() returns a pointer to a Session struct.
	// Configure the session such that it expires after 12 hours.
	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour

	// ========== Establish app dependencies for routes and handlers ========== //

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		config: cfg, // Pointer to app config
		snippets: &mysql.SnippetModel{DB: db}, // Pointer to SnippetModel
		users: &mysql.UserModel{DB: db}, // Pointer to UserModel
		templateCache: tc,
		session: session, // Add session manager to application dependencies
	}

	// ========== Create and run HTTP server ========== //

	// Initialize a new http.Server struct. We set the Addr and Handler fields so that the server
	// uses the same network address and routes as before, and set the ErrorLog field so that the server
	// now uses the custom errorLog logger in the event of any problems.
	srv := &http.Server {
		Addr: cfg.Addr,
		ErrorLog: errorLog, // Use custom error logger in the HTTP server
		Handler: app.routes(), // Return ServeMux with application routes
	}

	// The value returned from the flag.String() function is a pointer to the flag value,
	// not the value itself. So we need to dereference the pointer (i.e. prefix it with the * symbol)
	// before using it.
	infoLog.Printf("Starting server on %s", cfg.Addr) // Use custom info logger
	// Use the ListenAndServeTLS() method to start the HTTPS server. We
	// pass in the paths to the TLS certificatr and corresponding private key
	// as the two parameters.
	//err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	err = srv.ListenAndServe()
	errorLog.Fatal(err) // Use custom error logger
}

// openDB() wraps sql.Open and returns a sql.DB connection pool
// for a given DSN
func openDB(dsn string) (*sql.DB, error) {
	// sql.Open() doesn't create any connections, all it does is initialize the pool of connections for future use.
	// Actual connections to the database are established lazily, as and when needed for the first time.
	// To verify that everything is set up correct, we need to use db.Ping() to create a connection and check for errors.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping(); // Use db.Ping() to test if connection is successful
	if err != nil {
		return nil, err
	}

	// sql.DB is not a connection.
	// It is a pool of connections in the background managed by Go's database/sql package.
	return db, nil
}