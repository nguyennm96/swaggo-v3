package main

import (
	"github.com/nguyennm96/swaggo-v3"
	"github.com/nguyennm96/swaggo-v3/testdata/delims/api"
	_ "github.com/nguyennm96/swaggo-v3/testdata/delims/docs"
)

func ReadDoc() string {
	doc, _ := swaggo.ReadDoc("CustomDelims")
	return doc
}

// @title Swagger Example API
// @version 1.0
// @description Testing custom template delimeters
// @termsOfService http://swagger.io/terms/

func main() {
	api.MyFunc()
}
