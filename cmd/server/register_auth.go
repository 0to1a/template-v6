package main

import (
	"net/http"

	"connectrpc.com/connect"

	"project/internal/auth"
	"project/internal/gen/auth/authpbconnect"
)

// One file per domain; never a second registry
func registerAuth(mux *http.ServeMux, handler *auth.Handler, opts ...connect.HandlerOption) {
	mux.Handle(authpbconnect.NewAuthServiceHandler(handler, opts...))
}
