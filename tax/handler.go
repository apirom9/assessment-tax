package tax

import (
	"errors"
	"net/http"

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

type Err struct {
	Message string `json:"message"`
}

func CreateTaxCalculator(request CalculationRequest) (Calulator, error) {

	calculator := NewTaxCalulator()

	calculator.TotalIncome = request.TotalIncome
	calculator.WitholdingTax = request.WithHoldingTax
	for _, allowance := range request.Allowances {
		if allowance.Type == "donation" {
			calculator.AllowanceDonation = allowance.Amount
		} else {
			return calculator, errors.New("Unknown allowance type: " + allowance.Type)
		}
	}

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

	calculator, err := CreateTaxCalculator(request)
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
