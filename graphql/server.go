package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	handler "github.com/graphql-go/graphql-go-handler"
)

var (
	SslCrtFile string
	SslKeyFile string
	db         *sql.DB
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Weight      float64 `json:"weight"`
	Price       float64 `json:"price"`
	AssetUrl    string  `json:"asset_url"`
}

type Seller struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	WalletId  int    `json:"walletId"`
}

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"quantity": &graphql.Field{
			Type: graphql.Int,
		},
		"weight": &graphql.Field{
			Type: graphql.Float,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"asset_url": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var sellerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Seller",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"walletId": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"greeting": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "Hello World!", nil
			},
		},
		"product": &graphql.Field{
			Type:        productType,
			Description: "Get product by id or name",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					return getProductById(db, id)
				}
				name, ok := p.Args["name"].(string)
				if ok {
					return getProductByName(db, name)
				}
				return nil, nil
			},
		},
		"products": &graphql.Field{
			Type: graphql.NewList(productType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return getAllProducts(db)
			},
		},
		"sellProduct": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Add a new product to the marketplace catalog",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"quantity": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 1,
				},
				"weight": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"asset_url": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "https://bubbleerp.sysfosolutions.com/img/default-pro.jpg",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var product Product
				var ok bool
				product.Name, ok = p.Args["name"].(string)
				if !ok {
					return nil, errors.New("Missing required argument: name")
				}
				product.Description, ok = p.Args["description"].(string)
				if !ok {
					return nil, errors.New("Missing required argument: description")
				}
				product.Quantity, _ = p.Args["quantity"].(int) // optional with default value
				product.Weight, ok = p.Args["weight"].(float64)
				if !ok {
					return nil, errors.New("Missing required argument: weight")
				}
				product.Price, ok = p.Args["price"].(float64)
				if !ok {
					return nil, errors.New("Missing required argument: price")
				}
				product.AssetUrl, _ = p.Args["asset_url"].(string) // optional with default value

				return insertProduct(db, product)
			},
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

// Disable CORS policy from endpoint
func disableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	// initialize postgre database
	db = connectDatabase()
	defer db.Close()

	// create a graphl-go HTTP handler with our previously defined schema
	h := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: true,
	})

	// static file server to serve Graphiql in-browser editor
	fs := http.FileServer(http.Dir("static"))

	// serve a GraphQL endpoint at `/graphql`
	http.Handle("/graphql", disableCors(h))

	// serve a GraphiQL endpoint at `/`
	http.Handle("/", fs)

	// launch server
	if environment == "prod" {
		log.Fatal(http.ListenAndServeTLS(":8081", SslCrtFile, SslKeyFile, nil)) // https endpoint
	} else {
		log.Fatal(http.ListenAndServe(":8081", nil)) // http endpoint
	}
}
