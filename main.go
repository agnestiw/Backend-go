package main

import (
	"fmt"
	"latihan2/config"
	"latihan2/database"
	"log"
	"os"

	_ "github.com/gofiber/fiber/v2"
)

func main() {
	config.InitLogger()
	config.LoadEnv()

	fmt.Println("DEBUG: JWT_SECRET_KEY yang terbaca adalah ->", os.Getenv("JWT_SECRET_KEY"))

	database.InitDB()
	defer database.DB.Close()

	
	app := config.NewApp()
	log.Fatal(app.Listen(":3000"))
}
