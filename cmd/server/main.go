package main

import (
	"github.com/SchunckLeonardo/go-expert-api/configs"
	"github.com/SchunckLeonardo/go-expert-api/internal/entity"
	"github.com/SchunckLeonardo/go-expert-api/internal/infra/database"
	"github.com/SchunckLeonardo/go-expert-api/internal/infra/webserver/handlers"
	"github.com/SchunckLeonardo/go-expert-api/internal/infra/webserver/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"

	_ "github.com/SchunckLeonardo/go-expert-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Go Expert API Example
//	@version		1.0
//	@description	Product API with authentication
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Leonardo Schunck Rainha

//	@host						localhost:8080
//	@BasePath					/
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading configs: %v", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error to loading database connection: %v", err)
	}

	_ = db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDb := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDb)

	userDb := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDb, config.TokenAuth, config.JWTExpiresIn)

	r := chi.NewRouter()

	r.Use(middlewares.LogRequest)
	r.Use(middleware.Recoverer)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.FetchProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	// User
	r.Post("/users", userHandler.Create)
	r.Post("/sessions", userHandler.GetJWT)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))

	httpAddress := config.WebServerHost + ":" + config.WebServerPort
	log.Println("Server running on " + httpAddress)
	_ = http.ListenAndServe(httpAddress, r)
}
