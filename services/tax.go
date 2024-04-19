package services

type TaxLevel struct {
	Level             string
	MinAmount         float64
	MaxAmount         float64
	TaxRatePercentage float64
}

type TaxAllowance struct {
	Type   string
	Amount float64
}

type TaxCalulator struct {
	TaxLevels         []TaxLevel
	TaxAllowances     []TaxAllowance
	WitholdingTax     float64
	TotalIncome       float64
	PersonalDeduction float64
}

type TaxAmountLevel struct {
	Level     string
	TaxAmount float64
}

type TaxResult struct {
	TaxAmount       float64
	TaxAmountLevels []TaxAmountLevel
}

func (t *TaxCalulator) GetIncomeAfterPersonalDeduction() float64 {
	return t.TotalIncome - t.PersonalDeduction
}

func (t *TaxCalulator) CalculateTax(remainIncome, taxLevelMaxAmount, taxLevelPercentage float64) float64 {
	var incomeForTaxLevel float64
	if remainIncome > taxLevelMaxAmount {
		incomeForTaxLevel = taxLevelMaxAmount
	} else {
		incomeForTaxLevel = remainIncome
	}
	return (incomeForTaxLevel * taxLevelPercentage) / 100.00
}

func (t *TaxCalulator) CalculateTaxResult() TaxResult {

	var totalTaxAmount float64
	var taxAmountLevels []TaxAmountLevel
	remainIncome := t.GetIncomeAfterPersonalDeduction()

	for _, taxLevel := range t.TaxLevels {
		taxAmountLevel := TaxAmountLevel{Level: taxLevel.Level}
		if remainIncome > 0.0 {
			taxAmountLevel.TaxAmount = t.CalculateTax(remainIncome, taxLevel.MaxAmount, taxLevel.TaxRatePercentage)
			totalTaxAmount = totalTaxAmount + taxAmountLevel.TaxAmount
		}
		taxAmountLevels = append(taxAmountLevels, taxAmountLevel)
		remainIncome = remainIncome - taxLevel.MaxAmount
	}

	return TaxResult{
		TaxAmount:       totalTaxAmount,
		TaxAmountLevels: taxAmountLevels,
	}
}
