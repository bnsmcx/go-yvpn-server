package main

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "yvpn_server/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	r := chi.NewRouter()

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/docs/doc.json"), //The url pointing to API definition
	))

	http.ListenAndServe(":8000", r)
}

//func main() {
//	client := godo.NewFromToken(os.Getenv("DIGITAL_OCEAN_PAT"))
//	ctx := context.TODO()
//	fmt.Println(client.BillingHistory.List(ctx, nil))
//}
