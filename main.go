package main

import (
  "fmt"
  "net/http"

	"github.com/jinzhu/gorm"
	"github.com/qor/auth"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/providers/password"
	"github.com/qor/session/manager"

  _ "github.com/mattn/go-sqlite3"
)

var (
	// Initialize gorm DB
	gormDB, _ = gorm.Open("sqlite3", "sample.db")

	// Initialize Auth with configuration
	Auth = auth.New(&auth.Config{
		DB: gormDB,
  })
)

type User struct { 
  Name string
}

func init() {
	// Migrate AuthIdentity model, AuthIdentity will be used to save auth info, like username/password, oauth token, you could change that.
  gormDB.AutoMigrate(&auth_identity.AuthIdentity{})
  gormDB.AutoMigrate(&User{})

	// Register Auth providers
	// Allow use username/password
  Auth.RegisterProvider(password.New(&password.Config{}))
}

func main() {
	mux := http.NewServeMux()

  // Mount Auth to Router
  mux.Handle("/auth/", Auth.NewServeMux())
  mux.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":9000", manager.SessionManager.Middleware(mux))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "/auth/login")
}


