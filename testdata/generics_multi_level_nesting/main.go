package main

import (
	"net/http"

	"github.com/nguyennm96/swaggo-v3/testdata/generics_nested_my_version/api"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @host localhost:4000
// @basePath /api
func main() {
	http.HandleFunc("/posts/", api.GetPosts)
	http.ListenAndServe(":8080", nil)
}
