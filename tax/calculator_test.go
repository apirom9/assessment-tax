package tax

import (
	"reflect"
	"testing"
)

func TestTaxCalculation(t *testing.T) {
	t.Run("given total income 500000.0 should return tax 29000.0", func(t *testing.T) {

		taxCalulator := NewTaxCalulator()
		taxCalulator.TotalIncome = 500000.00

		wantAmount := 29000.0
		got := taxCalulator.CalculateTaxResult()
		if wantAmount != got.Amount {
			t.Errorf("expect tax = %v but got %v", wantAmount, got.Amount)
		}

		wantLevel := []LevelAmount{
			{Level: "0 - 150,000", Amount: 0.00},
			{Level: "150,001 - 500,000", Amount: 29000.00},
			{Level: "500,001 - 1,000,000", Amount: 0.00},
			{Level: "1,000,001 - 2,000,000", Amount: 0.00},
			{Level: "2,000,001 ขึ้นไป", Amount: 0.00},
		}
		if !reflect.DeepEqual(wantLevel, got.LevelAmounts) {
			t.Errorf("expected %v but got %v", wantLevel, got.LevelAmounts)
		}
	})

	t.Run("given total income 500000.0 and wht 25000.00 should return tax 4000.00", func(t *testing.T) {

		taxCalulator := NewTaxCalulator()
		taxCalulator.TotalIncome = 500000.00
		taxCalulator.WitholdingTax = 25000.00

		want := 4000.0
		got := taxCalulator.CalculateTaxResult()
		if want != got.Amount {
			t.Errorf("expect tax = %v but got %v", want, got.Amount)
		}
		wantLevel := []LevelAmount{
			{Level: "0 - 150,000", Amount: 0.00},
			{Level: "150,001 - 500,000", Amount: 29000.00},
			{Level: "500,001 - 1,000,000", Amount: 0.00},
			{Level: "1,000,001 - 2,000,000", Amount: 0.00},
			{Level: "2,000,001 ขึ้นไป", Amount: 0.00},
		}
		if !reflect.DeepEqual(wantLevel, got.LevelAmounts) {
			t.Errorf("expected %v but got %v", wantLevel, got.LevelAmounts)
		}
	})

	t.Run("given total income 500000.0 and allowance donate 200000.00 should return tax 19000.00", func(t *testing.T) {

		taxCalulator := NewTaxCalulator()
		taxCalulator.TotalIncome = 500000.00
		taxCalulator.AllowanceDonation = 200000.0

		want := 19000.00
		got := taxCalulator.CalculateTaxResult()
		if want != got.Amount {
			t.Errorf("expect tax = %v but got %v", want, got.Amount)
		}
		wantLevel := []LevelAmount{
			{Level: "0 - 150,000", Amount: 0.00},
			{Level: "150,001 - 500,000", Amount: 19000.00},
			{Level: "500,001 - 1,000,000", Amount: 0.00},
			{Level: "1,000,001 - 2,000,000", Amount: 0.00},
			{Level: "2,000,001 ขึ้นไป", Amount: 0.00},
		}
		if !reflect.DeepEqual(wantLevel, got.LevelAmounts) {
			t.Errorf("expected %v but got %v", wantLevel, got.LevelAmounts)
		}
	})

	t.Run("given total income 500000.0 and allowance donate 100000.00 k-receipt 200000.0 should return tax 14000.00", func(t *testing.T) {

		taxCalulator := NewTaxCalulator()
		taxCalulator.TotalIncome = 500000.00
		taxCalulator.AllowanceDonation = 100000.00
		taxCalulator.AllowanceKReceipt = 200000.00

		want := 14000.0
		got := taxCalulator.CalculateTaxResult()
		if want != got.Amount {
			t.Errorf("expect tax = %v but got %v", want, got.Amount)
		}
		wantLevel := []LevelAmount{
			{Level: "0 - 150,000", Amount: 0.00},
			{Level: "150,001 - 500,000", Amount: 14000.0},
			{Level: "500,001 - 1,000,000", Amount: 0.00},
			{Level: "1,000,001 - 2,000,000", Amount: 0.00},
			{Level: "2,000,001 ขึ้นไป", Amount: 0.00},
		}
		if !reflect.DeepEqual(wantLevel, got.LevelAmounts) {
			t.Errorf("expected %v but got %v", wantLevel, got.LevelAmounts)
		}
	})

	t.Run("given total income 5000000.0 should return tax 1051500.00", func(t *testing.T) {

		taxCalulator := NewTaxCalulator()
		taxCalulator.TotalIncome = 5000000.00

		wantAmount := 1051500.00
		got := taxCalulator.CalculateTaxResult()
		if wantAmount != got.Amount {
			t.Errorf("expect tax = %v but got %v", wantAmount, got.Amount)
		}

		wantLevel := []LevelAmount{
			{Level: "0 - 150,000", Amount: 0.00},
			{Level: "150,001 - 500,000", Amount: 50000.00},
			{Level: "500,001 - 1,000,000", Amount: 150000.00},
			{Level: "1,000,001 - 2,000,000", Amount: 400000.00},
			{Level: "2,000,001 ขึ้นไป", Amount: 451500.00},
		}
		if !reflect.DeepEqual(wantLevel, got.LevelAmounts) {
			t.Errorf("expected %v but got %v", wantLevel, got.LevelAmounts)
		}
	})

	t.Run("given total income 100000.0 and wht 1000.0 should return tax -1000.0", func(t *testing.T) {

		taxCalulator := NewTaxCalulator()
		taxCalulator.TotalIncome = 100000.00
		taxCalulator.WitholdingTax = 1000.00

		wantAmount := -1000.0
		got := taxCalulator.CalculateTaxResult()
		if wantAmount != got.Amount {
			t.Errorf("expect tax = %v but got %v", wantAmount, got.Amount)
		}

		wantLevel := []LevelAmount{
			{Level: "0 - 150,000", Amount: 0.00},
			{Level: "150,001 - 500,000", Amount: 0.00},
			{Level: "500,001 - 1,000,000", Amount: 0.00},
			{Level: "1,000,001 - 2,000,000", Amount: 0.00},
			{Level: "2,000,001 ขึ้นไป", Amount: 0.00},
		}
		if !reflect.DeepEqual(wantLevel, got.LevelAmounts) {
			t.Errorf("expected %v but got %v", wantLevel, got.LevelAmounts)
		}
	})
}
