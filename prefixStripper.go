package main

import (
	"net/http"
)

func prefixStripper(next http.Handler) http.Handler {
	return http.StripPrefix("/app", next)
}
