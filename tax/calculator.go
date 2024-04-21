package tax

import "math"

type Level struct {
	Level             string
	MinAmount         float64
	MaxAmount         float64
	TaxRatePercentage float64
}

type Calulator struct {
	Levels               []Level
	WitholdingTax        float64
	TotalIncome          float64
	AllowancePersonal    float64
	AllowanceDonation    float64
	AllowanceKReceipt    float64
	MaxAllowancePersonal float64
	MaxAllowanceDonation float64
	MaxAllowanceKReceipt float64
}

type LevelAmount struct {
	Level  string
	Amount float64
}

type Result struct {
	Amount       float64
	LevelAmounts []LevelAmount
}

func CreateLevels() []Level {
	return []Level{
		{Level: "0 - 150,000", MinAmount: 0.0, MaxAmount: 150000, TaxRatePercentage: 0},
		{Level: "150,001 - 500,000", MinAmount: 150001.00, MaxAmount: 500000.00, TaxRatePercentage: 10},
		{Level: "500,001 - 1,000,000", MinAmount: 500001.00, MaxAmount: 1000000.00, TaxRatePercentage: 15},
		{Level: "1,000,001 - 2,000,000", MinAmount: 1000001.00, MaxAmount: 2000000.00, TaxRatePercentage: 20},
		{Level: "2,000,001 ขึ้นไป", MinAmount: 2000001.00, MaxAmount: math.MaxFloat64, TaxRatePercentage: 35},
	}
}

func NewTaxCalulator(defaultAllowancePersonal, maxAllowanceKReceipt float64) Calulator {
	return Calulator{
		Levels:               CreateLevels(),
		TotalIncome:          0.00,
		AllowancePersonal:    defaultAllowancePersonal,
		AllowanceDonation:    0.00,
		AllowanceKReceipt:    0.00,
		MaxAllowancePersonal: 100000.00,
		MaxAllowanceDonation: 100000.00,
		MaxAllowanceKReceipt: maxAllowanceKReceipt,
	}
}

func (t *Calulator) GetAllowancePersonal() float64 {
	if t.AllowancePersonal > t.MaxAllowancePersonal {
		return t.MaxAllowancePersonal
	}
	return t.AllowancePersonal
}

func (t *Calulator) GetAllowanceDonation() float64 {
	if t.AllowanceDonation > t.MaxAllowanceDonation {
		return t.MaxAllowanceDonation
	}
	return t.AllowanceDonation
}

func (t *Calulator) GetAllowanceKReceipt() float64 {
	if t.AllowanceKReceipt > t.MaxAllowanceKReceipt {
		return t.MaxAllowanceKReceipt
	}
	return t.AllowanceKReceipt
}

func (t *Calulator) CalculateTax(remainIncome, taxLevelMaxAmount, taxLevelPercentage float64) float64 {
	var incomeForTaxLevel float64
	if remainIncome > taxLevelMaxAmount {
		incomeForTaxLevel = taxLevelMaxAmount
	} else {
		incomeForTaxLevel = remainIncome
	}
	return (incomeForTaxLevel * taxLevelPercentage) / 100.00
}

func (t *Calulator) CalculateIncomeAfterAllowances() float64 {
	return t.TotalIncome - t.GetAllowancePersonal() - t.GetAllowanceDonation() - t.GetAllowanceKReceipt()
}

func (t *Calulator) CalculateTaxResult() Result {

	var totalTaxAmount float64
	var taxAmountLevels []LevelAmount
	remainIncome := t.CalculateIncomeAfterAllowances()

	for _, taxLevel := range t.Levels {
		taxAmountLevel := LevelAmount{Level: taxLevel.Level}
		if remainIncome > 0.0 {
			taxAmountLevel.Amount = t.CalculateTax(remainIncome, taxLevel.MaxAmount, taxLevel.TaxRatePercentage)
			totalTaxAmount = totalTaxAmount + taxAmountLevel.Amount
		}
		taxAmountLevels = append(taxAmountLevels, taxAmountLevel)
		remainIncome = remainIncome - taxLevel.MaxAmount
	}

	totalTaxAmount = totalTaxAmount - t.WitholdingTax

	return Result{
		Amount:       totalTaxAmount,
		LevelAmounts: taxAmountLevels,
	}
}
