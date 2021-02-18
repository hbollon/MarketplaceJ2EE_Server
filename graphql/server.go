package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	handler "github.com/graphql-go/graphql-go-handler"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"firstName"`
	Description string  `json:"lastName"`
	Quantity    int     `json:"quantity"`
	Weight      float32 `json:"weight"`
	Price       float32 `json:"price"`
	AssetUrl    string  `json:"asset_url"`
}

var products []Product

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
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, product := range products {
						if int(product.Id) == id {
							return product, nil
						}
					}
				}
				name, ok := p.Args["name"].(string)
				if ok {
					for _, product := range products {
						if string(product.Name) == name {
							return product, nil
						}
					}
				}
				return nil, nil
			},
		},
		"products": &graphql.Field{
			Type: graphql.NewList(productType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return products, nil
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
	// fetch all products from db
	db := connectDatabase()
	defer db.Close()
	rows, err := db.Query("SELECT name, description, quantity, weight, price, asset_url FROM products")
	if err != nil {
		log.Fatal(err)
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

	http.ListenAndServe(":8081", nil)
}
