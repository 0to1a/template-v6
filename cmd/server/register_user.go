package main

import (
	"net/http"

	"connectrpc.com/connect"

	"project/internal/gen/user/userpbconnect"
	"project/internal/user"
)

// Second domain registry; same one-file-per-domain rule
func registerUser(mux *http.ServeMux, handler *user.Handler, opts ...connect.HandlerOption) {
	mux.Handle(userpbconnect.NewUserServiceHandler(handler, opts...))
}
