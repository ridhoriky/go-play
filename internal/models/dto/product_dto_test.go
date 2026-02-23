package dto

import (
	"testing"
)

func TestValidate_Valid(t *testing.T) {
	p := &ProductDTO{Name: "Laptop", Price: 1000, Stock: 5}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected valid, got error: %v", err)
	}
}

func TestValidate_Invalid(t *testing.T) {
	cases := []struct{
		name string
		p ProductDTO
	}{
		{"empty name", ProductDTO{Name: "", Price: 1000, Stock: 1}},
		{"negative price", ProductDTO{Name: "X", Price: -1, Stock: 1}},
		{"negative stock", ProductDTO{Name: "X", Price: 10, Stock: -5}},
	}

	for _, c := range cases {
		if err := c.p.Validate(); err == nil {
			t.Fatalf("case %s: expected error, got nil", c.name)
		}
	}
}
