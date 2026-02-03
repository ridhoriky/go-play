package product

import (
	"ne-project/internal/services/product"
	"net/http"
)


type ProductHandlerItf interface {
	HandleProducts(w http.ResponseWriter, r *http.Request)
	HandleProductByID(w http.ResponseWriter, r *http.Request)
}

type productHandler struct {
	productService product.ProductServiceItf
}

func NewProductHandler(productService product.ProductServiceItf) ProductHandlerItf {
	return &productHandler{
		productService: productService,
	}
}