package services

type TaxLevel struct {
	Level             string
	MinAmount         float64
	MaxAmount         float64
	TaxRatePercentage float64
}

type TaxCalulator struct {
	TaxLevels            []TaxLevel
	WitholdingTax        float64
	TotalIncome          float64
	AllowancePersonal    float64
	AllowanceDonation    float64
	AllowanceKReceipt    float64
	MaxAllowancePersonal float64
	MaxAllowanceDonation float64
	MaxAllowanceKReceipt float64
}

type TaxAmountLevel struct {
	Level     string
	TaxAmount float64
}

type TaxResult struct {
	TaxAmount       float64
	TaxAmountLevels []TaxAmountLevel
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

func NewTaxCalulator() TaxCalulator {
	return TaxCalulator{
		TaxLevels:            CreateTaxLevels(),
		TotalIncome:          0.00,
		AllowancePersonal:    60000.00,
		AllowanceDonation:    0.00,
		AllowanceKReceipt:    0.00,
		MaxAllowancePersonal: 100000.00,
		MaxAllowanceDonation: 100000.00,
		MaxAllowanceKReceipt: 50000.00,
	}
}

func (t *TaxCalulator) GetAllowancePersonal() float64 {
	if t.AllowancePersonal > t.MaxAllowancePersonal {
		return t.MaxAllowancePersonal
	}
	return t.AllowancePersonal
}

func (t *TaxCalulator) GetAllowanceDonation() float64 {
	if t.AllowanceDonation > t.MaxAllowanceDonation {
		return t.MaxAllowanceDonation
	}
	return t.AllowanceDonation
}

func (t *TaxCalulator) GetAllowanceKReceipt() float64 {
	if t.AllowanceKReceipt > t.MaxAllowanceKReceipt {
		return t.MaxAllowanceKReceipt
	}
	return t.AllowanceKReceipt
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

func (t *TaxCalulator) CalculateIncomeAfterAllowances() float64 {
	return t.TotalIncome - t.GetAllowancePersonal() - t.GetAllowanceDonation() - t.GetAllowanceKReceipt()
}

func (t *TaxCalulator) CalculateTaxResult() TaxResult {

	var totalTaxAmount float64
	var taxAmountLevels []TaxAmountLevel
	remainIncome := t.CalculateIncomeAfterAllowances()

	for _, taxLevel := range t.TaxLevels {
		taxAmountLevel := TaxAmountLevel{Level: taxLevel.Level}
		if remainIncome > 0.0 {
			taxAmountLevel.TaxAmount = t.CalculateTax(remainIncome, taxLevel.MaxAmount, taxLevel.TaxRatePercentage)
			totalTaxAmount = totalTaxAmount + taxAmountLevel.TaxAmount
		}
		taxAmountLevels = append(taxAmountLevels, taxAmountLevel)
		remainIncome = remainIncome - taxLevel.MaxAmount
	}

	totalTaxAmount = totalTaxAmount - t.WitholdingTax

	return TaxResult{
		TaxAmount:       totalTaxAmount,
		TaxAmountLevels: taxAmountLevels,
	}
}
