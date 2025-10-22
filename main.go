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
)

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

	log.Fatal(app.Listen(":3000"))
}
