package main

// import "github.com/dfang/qor-demo/config/bindatafs"

import (
  "fmt"
  "net/http"

	"github.com/jinzhu/gorm"
	"github.com/qor/auth"
	"github.com/qor/auth/auth_identity"
	// "github.com/qor/auth/providers/password"
  "github.com/qor/auth_themes/clean"
	"github.com/qor/session/manager"
  "github.com/qor/mailer"
  "github.com/qor/mailer/logger"
  // "github.com/qor/render"
  "github.com/qor/redirect_back"

  _ "github.com/mattn/go-sqlite3"
)

var (
  Mailer         *mailer.Mailer
  RedirectBack   = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
)

var (
	// Initialize gorm DB
	gormDB, _ = gorm.Open("sqlite3", "sample.db")

	// Initialize Auth with configuration
  Auth = clean.New(&auth.Config{
    DB:         gormDB,
    // NO NEED TO CONFIG RENDER, AS IT'S CONFIGED IN CLEAN THEME
    // Render:     render.New(&render.Config{AssetFileSystem: bindatafs.AssetFS.NameSpace("auth")}),
    Mailer:     Mailer,
    UserModel:  User{},
    Redirector: auth.Redirector{RedirectBack: RedirectBack},
  })
)

type User struct { 
  Name string
}

func init() {
	// Migrate AuthIdentity model, AuthIdentity will be used to save auth info, like username/password, oauth token, you could change that.
  gormDB.AutoMigrate(&auth_identity.AuthIdentity{})
  gormDB.AutoMigrate(&User{})

	// Register Auth providers, Allow use username/password
  // NO NEED TO REGISTER, AS IT'S REGISTERED IN CLEAN THEME
  // Auth.RegisterProvider(password.New(&password.Config{}))

  Mailer = mailer.New(&mailer.Config{
		Sender: logger.New(&logger.Config{}),
	})
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


