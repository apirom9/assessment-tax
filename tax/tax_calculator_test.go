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
		if wantAmount != got.TaxAmount {
			t.Errorf("expect tax = %v but got %v", wantAmount, got.TaxAmount)
		}

		wantLevel := []TaxAmountLevel{
			{Level: "0 - 150,000", TaxAmount: 0.00},
			{Level: "150,001 - 500,000", TaxAmount: 29000.00},
			{Level: "500,001 - 1,000,000", TaxAmount: 0.00},
			{Level: "1,000,001 - 2,000,000", TaxAmount: 0.00},
			{Level: "2,000,001 ขึ้นไป", TaxAmount: 0.00},
		}
		if !reflect.DeepEqual(wantLevel, got.TaxAmountLevels) {
			t.Errorf("expected %v but got %v", wantLevel, got.TaxAmountLevels)
		}
	})

	t.Run("given total income 500000.0 and wht 25000.00 should return tax 4000.00", func(t *testing.T) {

		taxCalulator := NewTaxCalulator()
		taxCalulator.TotalIncome = 500000.00
		taxCalulator.WitholdingTax = 25000.00

		want := 4000.0
		got := taxCalulator.CalculateTaxResult()
		if want != got.TaxAmount {
			t.Errorf("expect tax = %v but got %v", want, got.TaxAmount)
		}
		wantLevel := []TaxAmountLevel{
			{Level: "0 - 150,000", TaxAmount: 0.00},
			{Level: "150,001 - 500,000", TaxAmount: 29000.00},
			{Level: "500,001 - 1,000,000", TaxAmount: 0.00},
			{Level: "1,000,001 - 2,000,000", TaxAmount: 0.00},
			{Level: "2,000,001 ขึ้นไป", TaxAmount: 0.00},
		}
		if !reflect.DeepEqual(wantLevel, got.TaxAmountLevels) {
			t.Errorf("expected %v but got %v", wantLevel, got.TaxAmountLevels)
		}
	})

	t.Run("given total income 500000.0 and allowance donate 200000.00 should return tax 19000.00", func(t *testing.T) {

		taxCalulator := NewTaxCalulator()
		taxCalulator.TotalIncome = 500000.00
		taxCalulator.AllowanceDonation = 200000.0

		want := 19000.00
		got := taxCalulator.CalculateTaxResult()
		if want != got.TaxAmount {
			t.Errorf("expect tax = %v but got %v", want, got.TaxAmount)
		}
		wantLevel := []TaxAmountLevel{
			{Level: "0 - 150,000", TaxAmount: 0.00},
			{Level: "150,001 - 500,000", TaxAmount: 19000.00},
			{Level: "500,001 - 1,000,000", TaxAmount: 0.00},
			{Level: "1,000,001 - 2,000,000", TaxAmount: 0.00},
			{Level: "2,000,001 ขึ้นไป", TaxAmount: 0.00},
		}
		if !reflect.DeepEqual(wantLevel, got.TaxAmountLevels) {
			t.Errorf("expected %v but got %v", wantLevel, got.TaxAmountLevels)
		}
	})

	t.Run("given total income 500000.0 and allowance donate 100000.00 k-receipt 200000.0 should return tax 14000.00", func(t *testing.T) {

		taxCalulator := NewTaxCalulator()
		taxCalulator.TotalIncome = 500000.00
		taxCalulator.AllowanceDonation = 100000.00
		taxCalulator.AllowanceKReceipt = 200000.00

		want := 14000.0
		got := taxCalulator.CalculateTaxResult()
		if want != got.TaxAmount {
			t.Errorf("expect tax = %v but got %v", want, got.TaxAmount)
		}
		wantLevel := []TaxAmountLevel{
			{Level: "0 - 150,000", TaxAmount: 0.00},
			{Level: "150,001 - 500,000", TaxAmount: 14000.0},
			{Level: "500,001 - 1,000,000", TaxAmount: 0.00},
			{Level: "1,000,001 - 2,000,000", TaxAmount: 0.00},
			{Level: "2,000,001 ขึ้นไป", TaxAmount: 0.00},
		}
		if !reflect.DeepEqual(wantLevel, got.TaxAmountLevels) {
			t.Errorf("expected %v but got %v", wantLevel, got.TaxAmountLevels)
		}
	})
}
