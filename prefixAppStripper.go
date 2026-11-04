package main

import (
	"net/http"
)

func prefixAppStripper(next http.Handler) http.Handler {
	return http.StripPrefix("/app", next)
}
