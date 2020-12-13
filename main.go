package main

import (
	//"auth/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gitlab.blockrules.com/br/personal/TestRestaurant/routes"
)

func main() {
	e := godotenv.Load()

	if e != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(e)

	port := os.Getenv("PORT")

	// Handle routes
	http.Handle("/", routes.Handlers())

	// serve
	log.Printf("Server listening on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
