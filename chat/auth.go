package main

import "net/http"

type authHandler struct {
	next http.Handler
}

func (a *authHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, err := request.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		writer.Header().Set("Location", "/login")
		writer.WriteHeader(http.StatusTemporaryRedirect) // send response
		return
	}

	if err != nil {
		// other error
		http.Error(writer, err.Error(), http.StatusInternalServerError) // send error response
		return
	}

	// success call next handler
	a.next.ServeHTTP(writer, request)
}

func MustAuth(handler http.Handler) *authHandler {
	return &authHandler{next: handler}
}
