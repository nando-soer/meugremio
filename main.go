package main

import (
	"context"
	"embed"
	"log"
	"os"
	"tcc/controllers"
	"tcc/models"
	"tcc/routes"
	"tcc/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}

	studentModel := models.CreateStudentModel(pool)
	studentService := services.CreateStudentService(studentModel)
	studentController := controllers.CreateStudentController(studentService)

	router := gin.Default()
	routes.Init(router, studentController)

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
