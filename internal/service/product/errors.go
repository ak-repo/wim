package product

import "errors"

var ErrSKUExists = errors.New("product with this SKU already exists")
