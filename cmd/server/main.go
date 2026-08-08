// Business logic belongs under internal/, not here
package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"

	"connectrpc.com/connect"

	dbfs "project/db"
	"project/internal/auth"
	"project/internal/gen/auth/authpbconnect"
	"project/internal/gen/db"
	"project/internal/health"
	"project/internal/mail"
	"project/internal/platform/config"
	"project/internal/platform/database"
	"project/internal/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	// Up-only; failure aborts startup, never guesses schema
	migrations, err := fs.Sub(dbfs.Migrations, "migrations")
	if err != nil {
		return err
	}
	if err := database.Migrate(ctx, pool, migrations); err != nil {
		return err
	}

	defer registerBackground(ctx, pool)()

	queries := db.New(pool)

	jwtManager, err := auth.NewJWTManager(cfg.JWTSecret)
	if err != nil {
		return err
	}

	// Unset MAIL_URL discards codes; malformed aborts startup
	mailSender, err := mail.NewSMTPSenderFromURL(cfg.MailURL)
	if err != nil {
		return err
	}
	var loginCodeSender auth.LoginCodeSender = auth.NoopLoginCodeSender{}
	if mailSender != nil {
		loginCodeSender = auth.NewEmailLoginCodeSender(mailSender)
	}

	authHandler := auth.NewHandler(auth.NewService(
		auth.NewRepository(queries),
		loginCodeSender,
		jwtManager,
		cfg.IsGuestRegistration,
	))

	userHandler := user.NewHandler(user.NewService(user.NewRepository(queries)))

	// Unlisted procedures require auth; fail locked, not open
	publicProcedures := map[string]bool{
		authpbconnect.AuthServiceRequestLoginProcedure: true,
		authpbconnect.AuthServiceSubmitLoginProcedure:  true,
	}
	withAuth := connect.WithInterceptors(auth.NewInterceptor(jwtManager, publicProcedures))

	mux := http.NewServeMux()
	mux.Handle("GET /health", health.Handler())
	registerAuth(mux, authHandler, withAuth)
	registerUser(mux, userHandler, withAuth)

	if err := registerFrontend(mux); err != nil {
		return err
	}

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}
