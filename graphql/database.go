package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	connection string
	host       string
	port       int
	dbname     string
	user       string
	password   string
)

// init function called by Go before main execution and after variables definition
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connection = os.Getenv("DB_CONNECTION")
	host = os.Getenv("DB_HOST")
	port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Printf("Error: %v\nUsing default port\n", err)
		port = 5432
	}
	dbname = os.Getenv("DB_DATABASE")
	user = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
}

func connectDatabase() *sql.DB {
	// Define connection string for lib/pq
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)

	// Open db connection
	db, err := sql.Open(connection, psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Check db connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create table and insert default data if not exist
	err = initDb(db)
	if err != nil {
		log.Fatalf("Error during db initialization: %v\n", err)
	}

	return db
}

func initDb(db *sql.DB) error {
	var products = []Product{
		{
			Id:          1,
			Name:        "Test1",
			Description: "Un truc super lourd",
			Quantity:    100,
			Weight:      10.0,
			Price:       20.0,
		},
		{
			Id:          2,
			Name:        "Test2",
			Description: "Un truc super lourd",
			Quantity:    200,
			Weight:      40.0,
			Price:       10.0,
		},
		{
			Id:          3,
			Name:        "Test3",
			Description: "Un truc super lourd",
			Quantity:    500,
			Weight:      80.0,
			Price:       38.0,
		},
	}

	_, err := db.Query(
		"CREATE TABLE IF NOT EXISTS products (" +
			"id SERIAL," +
			"name varchar(40) NOT NULL PRIMARY KEY," +
			"description text NOT NULL," +
			"quantity integer NOT NULL," +
			"weight real NOT NULL," +
			"price real NOT NULL)",
	)
	if err != nil {
		return err
	}

	for _, p := range products {
		_, err = db.Query(
			"INSERT INTO products (name, description, quantity, weight, price) VALUES ('" +
				p.Name + "', '" +
				p.Description + "', " +
				fmt.Sprintf("%d", p.Quantity) + ", " +
				fmt.Sprintf("%f", p.Weight) + ", " +
				fmt.Sprintf("%f", p.Price) + ") " +
				"ON CONFLICT DO NOTHING",
		)
		if err != nil {
			return err
		}
	}

	return nil
}
