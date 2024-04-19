package services

import (
	"testing"
)

func TestTaxCalculation(t *testing.T) {
	t.Run("given total income 500000.0 should return tax 29000.0", func(t *testing.T) {

		taxCalulator := TaxCalulator{
			TaxLevels:         CreateTaxLevels(),
			TaxAllowances:     CreateTaxAllowances(),
			WitholdingTax:     0.00,
			TotalIncome:       500000.00,
			PersonalDeduction: 60000.00,
		}

		want := 29000.0
		got := taxCalulator.CalculateTaxResult()
		if want != got.TaxAmount {
			t.Errorf("expect tax = %v but got %v", want, got.TaxAmount)
		}
	})

	t.Run("given total income 500000.0 and wht 25000.00 should return tax 4000.00", func(t *testing.T) {

		taxCalulator := TaxCalulator{
			TaxLevels:         CreateTaxLevels(),
			TaxAllowances:     CreateTaxAllowances(),
			WitholdingTax:     25000.00,
			TotalIncome:       500000.00,
			PersonalDeduction: 60000.00,
		}

		want := 4000.0
		got := taxCalulator.CalculateTaxResult()
		if want != got.TaxAmount {
			t.Errorf("expect tax = %v but got %v", want, got.TaxAmount)
		}
	})
}

func CreateTaxLevels() []TaxLevel {
	return []TaxLevel{
		{Level: "0 - 150,000", MinAmount: 0.0, MaxAmount: 150000, TaxRatePercentage: 0},
		{Level: "150,001 - 500,000", MinAmount: 150001.00, MaxAmount: 500000.00, TaxRatePercentage: 10},
		{Level: "500,001 - 1,000,000", MinAmount: 500001.00, MaxAmount: 1000000.00, TaxRatePercentage: 15},
		{Level: "1,000,001 - 2,000,000", MinAmount: 1000001.00, MaxAmount: 2000000.00, TaxRatePercentage: 20},
		{Level: "2,000,001 ขึ้นไป", MinAmount: 2000001.00, MaxAmount: -1.00, TaxRatePercentage: 35},
	}
}

func CreateTaxAllowances() []TaxAllowance {
	return []TaxAllowance{
		{Type: "donation", Amount: 0.0},
		{Type: "k-receipt", Amount: 0.0},
	}
}
