package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "test",
		Price: 1.11,
		SKU:   "asd-dda-q",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
