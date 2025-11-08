package main

import (
	"context"
	"fmt"
	"latihan2/config"
	"latihan2/database"
	"log"
	"os"

	mongoRepo "latihan2/app/repository/mongo"

	_ "github.com/gofiber/fiber/v2"

	_ "latihan2/docs"

	"github.com/gofiber/swagger"
	// "github.com/swaggo/files"
	// "github.com/swaggo/gin-swagger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/


func main() {
	config.InitLogger()
	config.LoadEnv()

	fmt.Println("DEBUG: JWT_SECRET_KEY yang terbaca adalah ->", os.Getenv("JWT_SECRET_KEY"))

	database.InitPostgresDB()
	database.InitMongoDB()

	mongoRepo.InitPekerjaanCollection(database.MongoDB)

	defer database.DB.Close()
	defer database.MongoClient.Disconnect(context.Background())

	app := config.NewApp()

	// swagger gin
	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":3000"))
}
