package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/twotwo/go-blueprint/server/oapi"
)

func main() {
	// create a type that satisfies the `oapi.ServerInterface`, which contains an implementation of every operation from the generated code
	server := oapi.NewServer()

	r := chi.NewMux()

	// add CORS middleware to the router
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "vscode-webview://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // 限制可以跨域的 请求头
		AllowCredentials: true,
		MaxAge:           300, // 缓存预检请求 300 秒
	}))

	// get an `http.Handler` that we can use
	h := oapi.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
