package utils

import (
	//"auth/models"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //Gorm postgres dialect interface
	"github.com/joho/godotenv"
	"gitlab.blockrules.com/br/personal/TestRestaurant/models"
)

//ConnectDB function: Make database connection
func ConnectDB() *gorm.DB {

	fmt.Println("IN connect DB!!!")

	//Load environmenatal variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	databaseName := os.Getenv("databaseName")
	databaseHost := os.Getenv("databaseHost")

	//Define DB connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)

	//connect to db URI
	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}
	// close db when not in use
	// defer db.Close()

	// Migrate the schema
	fmt.Println("automigrate!!!")

	db.AutoMigrate(
		&models.Order{}, &models.MenuItem{}, &models.User{})

	user := models.User{Username: "bob", Password: "bobsecret", Address: "4 Memory Lane", Distance: 1.2}
	result := db.FirstOrCreate(&user) // pass pointer of data to Create
	if result.Error != nil {
		fmt.Println("error creating user")
	}

	menuitem := models.MenuItem{ID: 1, Description: "Big Salad", Price: 7.50, PrepTime: 5.0}
	result = db.FirstOrCreate(&menuitem) // pass pointer of data to Create
	if result.Error != nil {
		fmt.Println("error creating user")
	}
	menuitem = models.MenuItem{ID: 2, Description: "Family Meal", Price: 35.00, PrepTime: 10.0}
	result = db.FirstOrCreate(&menuitem)
	menuitem = models.MenuItem{ID: 3, Description: "Dessert Platter", Price: 15.00, PrepTime: 8.0}
	result = db.FirstOrCreate(&menuitem)

	fmt.Println("Successfully connected", db)
	return db
}
