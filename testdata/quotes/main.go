package main

import (
	"github.com/nguyennm96/swaggo-v3"
	"github.com/nguyennm96/swaggo-v3/testdata/quotes/api"
	_ "github.com/nguyennm96/swaggo-v3/testdata/quotes/docs"
)

func ReadDoc() string {
	doc, _ := swaggo.ReadDoc()
	return doc
}

// @title Swagger Example API
// @version 1.0
// @description.markdown
// @tag.name api
// @tag.description.markdown
// @termsOfService http://swagger.io/terms/

func main() {
	api.RandomFunc()
}
