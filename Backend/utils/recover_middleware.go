package utils

import (
	"fmt"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defer the recover function to catch any panics
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace if there's a panic
				fmt.Printf("Recovered from panic: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Call the next handler (request processing)
		next.ServeHTTP(w, r)
	})
}
