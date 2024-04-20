package tax

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AllowanceRequest struct {
	Type   string  `json:"allowanceType" example:"donation"`
	Amount float64 `json:"amount" example:"0.0"`
}

type CalculationRequest struct {
	TotalIncome    float64            `json:"totalIncome" example:"500000.0"`
	WithHoldingTax float64            `json:"wht" example:"0.0"`
	Allowances     []AllowanceRequest `json:"allowances"`
}

type TaxLevelResponse struct {
	Level     string  `json:"level" example:"0-150,000"`
	TaxAmount float64 `json:"tax" example:"0.0"`
}

type Response struct {
	Tax               float64            `json:"tax" example:"29000.0"`
	TaxLevelResponses []TaxLevelResponse `json:"taxLevel"`
}

type ResponseTaxResultForCSV struct {
	TotalIncome float64 `json:"totalIncome" example:"29000.0"`
	Tax         float64 `json:"tax,omitempty" example:"29000.0"`
	TaxRefund   float64 `json:"taxRefund,omitempty" example:"29000.0"`
}

type ResponseForCSV struct {
	Taxes []ResponseTaxResultForCSV `json:"taxes"`
}

type Err struct {
	Message string `json:"message"`
}

func CreateTaxCalculatorFromRequest(request CalculationRequest) (Calulator, error) {

	calculator := NewTaxCalulator()

	calculator.TotalIncome = request.TotalIncome
	calculator.WitholdingTax = request.WithHoldingTax
	for _, allowance := range request.Allowances {
		if allowance.Type == "donation" {
			calculator.AllowanceDonation = allowance.Amount
		} else if allowance.Type == "k-receipt" {
			calculator.AllowanceKReceipt = allowance.Amount
		} else {
			return calculator, errors.New("Unknown allowance type: " + allowance.Type)
		}
	}

	return calculator, nil
}

func CreateTaxCalculatorFromCsvRecord(record []string) (Calulator, error) {

	calculator := NewTaxCalulator()

	totalIncome, err := strconv.ParseFloat(record[0], 64)
	if err != nil {
		return calculator, err
	}
	calculator.TotalIncome = totalIncome

	wht, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return calculator, err
	}
	calculator.WitholdingTax = wht

	donation, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return calculator, err
	}
	calculator.AllowanceDonation = donation

	return calculator, nil
}

// CalculateTax
//
//	@Summary		Calculate Tax
//	@Description	Calculate Tax
//	@Tags			tax
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Router			/tax/calculations [post]
//	@Failure		500	{object}	Err
//	@Failure		400	{object}	Err
//	@Param 			CalculationRequest body CalculationRequest true "Body for calculation request"
func CalculateTax(c echo.Context) error {

	var request CalculationRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	calculator, err := CreateTaxCalculatorFromRequest(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	result := calculator.CalculateTaxResult()
	var taxLevelResponses []TaxLevelResponse
	for _, level := range result.LevelAmounts {
		taxLevelResponses = append(taxLevelResponses, TaxLevelResponse{
			Level:     level.Level,
			TaxAmount: level.Amount,
		})
	}
	return c.JSON(http.StatusOK, Response{Tax: result.Amount, TaxLevelResponses: taxLevelResponses})
}

// CalculateTaxCsv
//
//	@Summary		Calculate Tax for upload CSV file
//	@Description	Calculate Tax for upload CSV file
//	@Tags			tax
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	Response
//	@Router			/tax/calculations/upload-csv [post]
//	@Failure		500	{object}	Err
//	@Failure		400	{object}	Err
//	@Param 			taxes.csv formData file true "Uploaded CSV for tax calculation"
func CalculateTaxCsv(c echo.Context) error {
	file, err := c.FormFile("taxes.csv")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	reader := csv.NewReader(src)
	content, err := reader.ReadAll()
	if err != nil {
		return err
	}
	var response ResponseForCSV
	for index, record := range content {
		if index == 0 {
			// TODO check column names
			continue
		}
		calculator, err := CreateTaxCalculatorFromCsvRecord(record)
		if err != nil {
			return err
		}
		responseTaxResultForCSV := ResponseTaxResultForCSV{TotalIncome: calculator.TotalIncome}
		taxResult := calculator.CalculateTaxResult()
		if taxResult.Amount < 0 {
			responseTaxResultForCSV.TaxRefund = -taxResult.Amount
		} else {
			responseTaxResultForCSV.Tax = taxResult.Amount
		}
		response.Taxes = append(response.Taxes, responseTaxResultForCSV)
	}
	return c.JSON(http.StatusOK, response)
}
