package domain

import "testing"

func TestInventoryAvailableQty(t *testing.T) {
	inv := Inventory{Quantity: 20, ReservedQty: 7}

	if got := inv.AvailableQty(); got != 13 {
		t.Fatalf("expected available qty 13, got %d", got)
	}
}
