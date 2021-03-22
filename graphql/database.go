package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Embedded type alias to allow methods creation while
// keeping access to parent fields/methods
type Database struct {
	*sql.DB
}

var (
	environment string
	db_url      string
)

// init function called by Go before main execution and after variables definition
func init() {
	environment = os.Getenv("ENVIRONMENT")
	db_url = os.Getenv("DATABASE_URL")
}

func connectDatabase() (db Database) {
	// Open db connection
	var err error
	db.DB, err = sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	// Check db connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create table and insert default data if not exist
	err = db.initDb()
	if err != nil {
		log.Fatalf("Error during db initialization: %v\n", err)
	}

	return db
}

func (db *Database) initDb() error {
	var sellers = []Seller{
		{
			FirstName: "Hugo",
			LastName:  "Bollon",
			Email:     "hugo.bollon@gmail.com",
			WalletId:  100218510,
		},
	}

	var products = []Product{
		{
			Name:        "Cyberpunk 2077",
			Description: "Cyberpunk 2077 est un jeu vidéo d''action-RPG en vue à la première personne développé par CD Projekt RED, fondé sur la série de jeu de rôle sur table Cyberpunk 2020 conçue par Mike Pondsmith.",
			Quantity:    100,
			Weight:      10.0,
			Price:       50.0,
			AssetUrl:    "https://images-na.ssl-images-amazon.com/images/I/81%2BpdXH3fjL._AC_SY500_.jpg",
			Seller:      sellers[0],
		},
		{
			Name:        "Assassin''s Creed Valhalla",
			Description: "Assassin''s Creed Valhalla est un jeu vidéo d''action-aventure et de rôle, développé par Ubisoft Montréal et édité par Ubisoft, sorti en novembre 2020 sur Microsoft Windows.",
			Quantity:    100,
			Weight:      8.0,
			Price:       59.99,
			AssetUrl:    "https://image.jeuxvideo.com/medias/158826/1588264397-5261-jaquette-avant.jpg",
			Seller:      sellers[0],
		},
		{
			Name:        "Red Dead Redemption II",
			Description: "Red Dead Redemption II est un jeu vidéo d''action-aventure et de western multiplateforme, développé par Rockstar Studios et édité par Rockstar Games, sorti le 26 octobre 2018 sur PlayStation 4 et Xbox One et le 5 novembre 2019 sur Windows.",
			Quantity:    100,
			Weight:      10.0,
			Price:       50.0,
			AssetUrl:    "https://lh3.googleusercontent.com/HCUkD69MAHEOj84Yi7Kb5vxHpCePTsmQI4g9vYuVPUo-87cWE6ZZIk0tiyYzaiS9zaAFMTXRNYJaaRczRN-yQYw",
			Seller:      sellers[0],
		},
	}

	if *flags.ResetDbFlag {
		if err := db.resetDb(); err != nil {
			return err
		}
	}

	// create `seller` table if not exists
	_, err := db.Query(
		"CREATE TABLE IF NOT EXISTS seller (" +
			"id SERIAL PRIMARY KEY, " +
			"first_name varchar(40) NOT NULL, " +
			"last_name varchar(40) NOT NULL, " +
			"email varchar(255) UNIQUE NOT NULL, " +
			"wallet_id integer UNIQUE NOT NULL)",
	)
	if err != nil {
		return err
	}

	// create `product` table if not exists
	_, err = db.Query(
		"CREATE TABLE IF NOT EXISTS product (" +
			"id SERIAL PRIMARY KEY, " +
			"name varchar(40) UNIQUE NOT NULL, " +
			"description text NOT NULL, " +
			"quantity integer NOT NULL, " +
			"weight real NOT NULL, " +
			"price real NOT NULL, " +
			"asset_url text, " +
			"seller_id integer NOT NULL, " +
			"CONSTRAINT fk_seller " +
			"FOREIGN KEY(seller_id) " +
			"REFERENCES seller(id)" +
			"ON DELETE SET NULL);",
	)
	if err != nil {
		return err
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

	for _, p := range products {
		seller, _ := db.getSellerByEmail(p.Seller.Email)
		_, err = db.Query(
			"INSERT INTO product (name, description, quantity, weight, price, asset_url, seller_id) VALUES ('" +
				p.Name + "', '" +
				p.Description + "', " +
				fmt.Sprintf("%d", p.Quantity) + ", " +
				fmt.Sprintf("%f", p.Weight) + ", " +
				fmt.Sprintf("%f", p.Price) + ", '" +
				p.AssetUrl + "', " +
				fmt.Sprintf("%d", seller.Id) + ") " +
				"ON CONFLICT DO NOTHING",
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Reset whole database if `reset-db` flag is setted
func (db *Database) resetDb() error {
	tables := []string{"product", "seller"}

	for _, table := range tables {
		_, err := db.Query(
			fmt.Sprintf("DROP TABLE IF EXISTS %s", table),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

/* Product interractions */

func (db *Database) getAllProducts() ([]Product, error) {
	var products []Product
	rows, err := db.Query("SELECT id, name, description, quantity, weight, price, asset_url, seller_id FROM product")
	if err != nil {
		return nil, err
	}

	// unmarshall result rows to Product
	for rows.Next() {
		var p Product
		var sellerId int
		err = rows.Scan(
			&p.Id,
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Weight,
			&p.Price,
			&p.AssetUrl,
			&sellerId,
		)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		if p.Id != 0 {
			p.Seller, err = db.getSellerById(sellerId)
			if err != nil {
				log.Fatalf("Failled getting seller from id: %v", err)
			}
		}
		products = append(products, p)
	}

	return products, nil
}

func (db *Database) getProductById(id int) (Product, error) {
	var p Product
	var sellerId int
	// Prepare query, takes a name argument
	query, err := db.Prepare("SELECT id, name, description, quantity, weight, price, asset_url, seller_id FROM product WHERE id=$1")
	if err != nil {
		return p, err
	}

	// Make query with our stmt, passing in name argument
	var rows *sql.Rows
	rows, err = query.Query(id)
	if err != nil {
		return p, err
	}

	// Unmarshal result rows to Product
	if rows.Next() {
		err = rows.Scan(
			&p.Id,
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Weight,
			&p.Price,
			&p.AssetUrl,
			&sellerId,
		)
	}
	if err != nil {
		return p, fmt.Errorf("Scan: %v", err)
	}
	if p.Id != 0 {
		p.Seller, err = db.getSellerById(sellerId)
		if err != nil {
			log.Fatalf("Failled getting seller from id: %v", err)
		}
	}

	return p, nil
}

func (db *Database) getProductByName(name string) (Product, error) {
	var p Product
	var sellerId int
	// Prepare query, takes a name argument
	query, err := db.Prepare("SELECT id, name, description, quantity, weight, price, asset_url, seller_id FROM product WHERE name=$1")
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
			&p.Id,
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Weight,
			&p.Price,
			&p.AssetUrl,
			&sellerId,
		)
	}
	if err != nil {
		return p, fmt.Errorf("Scan: %v", err)
	}
	if p.Id != 0 {
		p.Seller, err = db.getSellerById(sellerId)
		if err != nil {
			log.Fatalf("Failled getting seller from id: %v", err)
		}
	}

	return p, nil
}

func (db *Database) insertProduct(p Product) (bool, error) {
	res, err := db.getProductByName(p.Name)
	if res == (Product{}) && err == nil {
		seller, _ := db.getSellerByEmail(p.Seller.Email)
		_, err = db.Query(
			"INSERT INTO product (name, description, quantity, weight, price, asset_url, seller_id) VALUES ('" +
				p.Name + "', '" +
				p.Description + "', " +
				fmt.Sprintf("%d", p.Quantity) + ", " +
				fmt.Sprintf("%f", p.Weight) + ", " +
				fmt.Sprintf("%f", p.Price) + ", '" +
				p.AssetUrl + "', " +
				fmt.Sprintf("%d", seller.Id) + ") " +
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

func (db *Database) getAllSellers() ([]Seller, error) {
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

func (db *Database) getSellerByEmail(email string) (Seller, error) {
	var s Seller
	// Prepare query, takes a name argument
	query, err := db.Prepare("SELECT id, first_name, last_name, email, wallet_id FROM seller WHERE email=$1")
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
			&s.Id,
			&s.FirstName,
			&s.LastName,
			&s.Email,
			&s.WalletId,
		)
	}
	if err != nil {
		return s, fmt.Errorf("Scan: %v", err)
	}
	if s.Id == 0 {
		return s, fmt.Errorf("User not found")
	}

	return s, nil
}

func (db *Database) getSellerById(id int) (Seller, error) {
	var s Seller
	// Prepare query, takes a name argument
	query, err := db.Prepare("SELECT first_name, last_name, email, wallet_id FROM seller WHERE id=$1")
	if err != nil {
		return s, err
	}

	// Make query with our stmt, passing in name argument
	var rows *sql.Rows
	rows, err = query.Query(id)
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
	if s.Email == "" {
		return s, fmt.Errorf("User not found")
	}

	return s, nil
}

func (db *Database) insertSeller(s Seller) (bool, error) {
	res, err := db.getSellerByEmail(s.Email)
	if res == (Seller{}) && err != nil {
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
