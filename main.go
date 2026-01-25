package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Indomie", Price: 2000, Stock: 20},
	{ID: 2, Name: "Bakso", Price: 10000, Stock: 10},
	{ID: 3, Name: "Mie Ayam", Price: 10000, Stock: 10},
	{ID: 4, Name: "Kelapa", Price: 10000, Stock: 10},
	{ID: 5, Name: "Jambu", Price: 10000, Stock: 10},
	{ID: 6, Name: "Manngga", Price: 10000, Stock: 10},
}

func main() {
	// GET localhost:8080/health

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API running",
		})

	})

	// GET localhost:8080/products
	// POST localhost:8080/products
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "POST" {
			// validation payload
			var newProduct Product
			err := json.NewDecoder(r.Body).Decode(&newProduct)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			newProduct.ID = len(products) + 1
			products = append(products, newProduct)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newProduct)

		} else if r.Method == "GET" {
			json.NewEncoder(w).Encode(products)
		}
	})

	// GET localhost:8080/products/{:id}
	// PUT localhost:8080/products/{:id}
	// DELETE localhost:8080/products/{:id}
	http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "PUT" {
			updateProduct(w, r)
		} else if r.Method == "GET" {
			getProductById(w, r)
		} else if r.Method == "DELETE" {
			deleteProduct(w, r)
		}
	})

	fmt.Println("server running di port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	//validation id
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product id", http.StatusBadRequest)
		return
	}
	for _, p := range products {
		if p.ID == id {
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Product Not found", http.StatusNotFound)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	//validation id
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product id", http.StatusBadRequest)
		return
	}

	// validation payload
	var updateProduct Product
	err = json.NewDecoder(r.Body).Decode(&updateProduct)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i := range products {
		if products[i].ID == id {
			updateProduct.ID = id
			products[i] = updateProduct
			json.NewEncoder(w).Encode(updateProduct)
			return
		}
	}

	http.Error(w, "Product Not found", http.StatusNotFound)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	//validation id
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product id", http.StatusBadRequest)
		return
	}
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "delete product success",
			})
			return
		}
	}

	http.Error(w, "Product Not found", http.StatusNotFound)
}
