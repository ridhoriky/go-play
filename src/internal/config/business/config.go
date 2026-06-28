package business

type BusinessConfig struct {
	MinStockCheckout int `yaml:"min_stock_checkout" env:"BUSINESS_MIN_STOCK_CHECKOUT" env-default:"1"`
	MaxCartItems     int `yaml:"max_cart_items" env:"BUSINESS_MAX_CART_ITEMS" env-default:"50"`
	OrderExpiryHours int `yaml:"order_expiry_hours" env:"BUSINESS_ORDER_EXPIRY_HOURS" env-default:"24"`
}
