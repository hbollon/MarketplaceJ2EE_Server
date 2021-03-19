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
	environment string
	connection  string
	host        string
	port        int
	dbname      string
	user        string
	password    string
)

// init function called by Go before main execution and after variables definition
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	environment = os.Getenv("ENVIRONMENT")
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
	SslCrtFile = os.Getenv("SSL_CRT_FILE")
	SslKeyFile = os.Getenv("SSL_KEY_FILE")
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
			Name:        "Cyberpunk 2077",
			Description: "Cyberpunk 2077 est un jeu vidéo d''action-RPG en vue à la première personne développé par CD Projekt RED, fondé sur la série de jeu de rôle sur table Cyberpunk 2020 conçue par Mike Pondsmith.",
			Quantity:    100,
			Weight:      10.0,
			Price:       50.0,
			AssetUrl:    "https://images-na.ssl-images-amazon.com/images/I/81%2BpdXH3fjL._AC_SY500_.jpg",
		},
		{
			Name:        "Assassin''s Creed Valhalla",
			Description: "Assassin''s Creed Valhalla est un jeu vidéo d''action-aventure et de rôle, développé par Ubisoft Montréal et édité par Ubisoft, sorti en novembre 2020 sur Microsoft Windows.",
			Quantity:    100,
			Weight:      8.0,
			Price:       59.99,
			AssetUrl:    "https://image.jeuxvideo.com/medias/158826/1588264397-5261-jaquette-avant.jpg",
		},
		{
			Name:        "Red Dead Redemption II",
			Description: "Red Dead Redemption II est un jeu vidéo d''action-aventure et de western multiplateforme, développé par Rockstar Studios et édité par Rockstar Games, sorti le 26 octobre 2018 sur PlayStation 4 et Xbox One et le 5 novembre 2019 sur Windows.",
			Quantity:    100,
			Weight:      10.0,
			Price:       50.0,
			AssetUrl:    "https://lh3.googleusercontent.com/HCUkD69MAHEOj84Yi7Kb5vxHpCePTsmQI4g9vYuVPUo-87cWE6ZZIk0tiyYzaiS9zaAFMTXRNYJaaRczRN-yQYw",
		},
	}

	var sellers = []Seller{
		{
			FirstName: "Hugo",
			LastName:  "Bollon",
			Email:     "hugo.bollon@gmail.com",
			WalletId:  100218510,
		},
	}

	// create `product` table if not exists
	_, err := db.Query(
		"CREATE TABLE IF NOT EXISTS product (" +
			"id SERIAL PRIMARY KEY," +
			"name varchar(40) UNIQUE NOT NULL," +
			"description text NOT NULL," +
			"quantity integer NOT NULL," +
			"weight real NOT NULL," +
			"price real NOT NULL," +
			"asset_url text)",
	)
	if err != nil {
		return err
	}

	// create `seller` table if not exists
	_, err = db.Query(
		"CREATE TABLE IF NOT EXISTS seller (" +
			"id SERIAL PRIMARY KEY," +
			"first_name varchar(40) NOT NULL," +
			"last_name varchar(40) NOT NULL," +
			"email varchar(255) UNIQUE NOT NULL," +
			"wallet_id integer UNIQUE NOT NULL)",
	)
	if err != nil {
		return err
	}

	for _, p := range products {
		_, err = db.Query(
			"INSERT INTO product (name, description, quantity, weight, price, asset_url) VALUES ('" +
				p.Name + "', '" +
				p.Description + "', " +
				fmt.Sprintf("%d", p.Quantity) + ", " +
				fmt.Sprintf("%f", p.Weight) + ", " +
				fmt.Sprintf("%f", p.Price) + ", '" +
				p.AssetUrl + "') " +
				"ON CONFLICT DO NOTHING",
		)
		if err != nil {
			return err
		}
	}

	for _, s := range sellers {
		_, err = db.Query(
			"INSERT INTO seller (first_name, last_name, email, wallet_id) VALUES ('" +
				s.FirstName + "', '" +
				s.LastName + "', '" +
				s.Email + "', " +
				fmt.Sprintf("%d", s.WalletId) + ") " +
				"ON CONFLICT DO NOTHING",
		)
		if err != nil {
			return err
		}
	}

	return nil
}

/* Product interractions */

func getAllProducts(db *sql.DB) ([]Product, error) {
	var products []Product
	rows, err := db.Query("SELECT name, description, quantity, weight, price, asset_url FROM product")
	if err != nil {
		return nil, err
	}

	// unmarshall result rows to Product
	for rows.Next() {
		var p Product
		err = rows.Scan(
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Weight,
			&p.Price,
			&p.AssetUrl,
		)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		products = append(products, p)
	}

	return products, nil
}

func getProductById(db *sql.DB, id int) (Product, error) {
	var p Product
	// Prepare query, takes a name argument
	query, err := db.Prepare("SELECT name, description, quantity, weight, price, asset_url FROM product WHERE id=$1")
	if err != nil {
		return p, err
	}

	// Make query with our stmt, passing in name argument
	var rows *sql.Rows
	rows, err = query.Query(id)
	if err != nil {
		return p, err
	}

	// Unmarshall result rows to Product
	if rows.Next() {
		err = rows.Scan(
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Weight,
			&p.Price,
			&p.AssetUrl,
		)
	}
	if err != nil {
		return p, fmt.Errorf("Scan: %v", err)
	}

	return p, nil
}

func getProductByName(db *sql.DB, name string) (Product, error) {
	var p Product
	// Prepare query, takes a name argument
	query, err := db.Prepare("SELECT name, description, quantity, weight, price, asset_url FROM product WHERE name=$1")
	if err != nil {
		return p, err
	}

	// Make query with our stmt, passing in name argument
	var rows *sql.Rows
	rows, err = query.Query(name)
	if err != nil {
		return p, err
	}

	// Unmarshal result rows to Product instances
	if rows.Next() {
		err = rows.Scan(
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Weight,
			&p.Price,
			&p.AssetUrl,
		)
	}
	if err != nil {
		return p, fmt.Errorf("Scan: %v", err)
	}

	return p, nil
}

func insertProduct(db *sql.DB, p Product) (bool, error) {
	res, err := getProductByName(db, p.Name)
	if res == (Product{}) && err == nil {
		_, err = db.Query(
			"INSERT INTO product (name, description, quantity, weight, price, asset_url) VALUES ('" +
				p.Name + "', '" +
				p.Description + "', " +
				fmt.Sprintf("%d", p.Quantity) + ", " +
				fmt.Sprintf("%f", p.Weight) + ", " +
				fmt.Sprintf("%f", p.Price) + ", '" +
				p.AssetUrl + "') " +
				"ON CONFLICT DO NOTHING",
		)
	} else {
		err = fmt.Errorf("The product '%s' is already registered.", p.Name)
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

/* Seller interractions */

func getAllSellers(db *sql.DB) ([]Seller, error) {
	var sellers []Seller
	rows, err := db.Query("SELECT first_name, last_name, email, wallet_id FROM seller")
	if err != nil {
		return nil, err
	}

	// unmarshal result rows to Seller instances
	for rows.Next() {
		var s Seller
		err = rows.Scan(
			&s.FirstName,
			&s.LastName,
			&s.Email,
			&s.WalletId,
		)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		sellers = append(sellers, s)
	}

	return sellers, nil
}

func getSellerByEmail(db *sql.DB, email string) (Seller, error) {
	var s Seller
	// Prepare query, takes a name argument
	query, err := db.Prepare("SELECT first_name, last_name, email, wallet_id FROM seller WHERE email=$1")
	if err != nil {
		return s, err
	}

	// Make query with our stmt, passing in name argument
	var rows *sql.Rows
	rows, err = query.Query(email)
	if err != nil {
		return s, err
	}

	// Unmarshal result rows to Seller instance
	if rows.Next() {
		err = rows.Scan(
			&s.FirstName,
			&s.LastName,
			&s.Email,
			&s.WalletId,
		)
	}
	if err != nil {
		return s, fmt.Errorf("Scan: %v", err)
	}

	return s, nil
}

func insertSeller(db *sql.DB, s Seller) (bool, error) {
	res, err := getSellerByEmail(db, s.Email)
	if res == (Seller{}) && err == nil {
		_, err = db.Query(
			"INSERT INTO seller (first_name, last_name, email, wallet_id) VALUES ('" +
				s.FirstName + "', '" +
				s.LastName + "', '" +
				s.Email + "', " +
				fmt.Sprintf("%d", s.WalletId) + ") " +
				"ON CONFLICT DO NOTHING",
		)
	} else {
		err = fmt.Errorf("The email '%s' is already registered as seller.", s.Email)
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
