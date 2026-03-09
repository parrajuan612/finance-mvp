package main

import (
	"log"
	"os"

	"finanzas-mvp/cmd/api"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró el archivo .env, usando variables de entorno del sistema")
	}
	file, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("cannot open error.log:", err)
	}
	log.SetOutput(file)
	r := gin.Default()
	api.InitServer(r)
}
