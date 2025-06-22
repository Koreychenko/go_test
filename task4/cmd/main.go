package main

import (
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"main/internal/user"
	"main/internal/webserver"
	"os"
)

const address = ":9993"

// @title           User CRUD API
// @version         1.0
// @description     This is the example of CRUD API using go
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	srv := webserver.NewWebserver(address)

	user.RegisterHandlers(srv)

	srv.RegisterHandler("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	err := srv.Start()

	if err != nil {
		fmt.Printf("Error %v", err)

		os.Exit(1)
	}
}
