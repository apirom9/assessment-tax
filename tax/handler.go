package tax

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Store interface {
	UpdateDefaultPersonalDeduction(value float64) error
	GetDefaultPersonalDeduction() (float64, error)
	UpdateMaxKReceipt(value float64) error
	GetMaxKReceipt() (float64, error)
}

type Handler struct {
	Store Store
}

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
	Tax               float64            `json:"tax,omitempty" example:"29000.0"`
	TaxRefund         float64            `json:"taxRefund,omitempty" example:"29000.0"`
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

type UpdatePersonalDeductionRequest struct {
	Amount float64 `json:"amount" example:"29000.0"`
}

type UpdatePersonalDeductionResponse struct {
	Amount float64 `json:"personalDeduction" example:"29000.0"`
}

type UpdateKReceiptRequest struct {
	Amount float64 `json:"amount" example:"29000.0"`
}

type UpdateKReceiptsResponse struct {
	Amount float64 `json:"kReceipt" example:"29000.0"`
}

func (h *Handler) CreateTaxCalculatorFromRequest(request CalculationRequest) (Calulator, error) {

	defaultPersonalAllowance, err := h.Store.GetDefaultPersonalDeduction()
	if err != nil {
		return Calulator{}, err
	}
	maxKreceipt, err := h.Store.GetMaxKReceipt()
	if err != nil {
		return Calulator{}, err
	}
	calculator := NewTaxCalulator(defaultPersonalAllowance, maxKreceipt)

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

func (h *Handler) CreateTaxCalculatorFromCsvRecord(record []string) (Calulator, error) {

	defaultPersonalAllowance, err := h.Store.GetDefaultPersonalDeduction()
	if err != nil {
		return Calulator{}, err
	}
	maxKreceipt, err := h.Store.GetMaxKReceipt()
	if err != nil {
		return Calulator{}, err
	}
	calculator := NewTaxCalulator(defaultPersonalAllowance, maxKreceipt)

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
func (h *Handler) CalculateTax(c echo.Context) error {

	var request CalculationRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	calculator, err := h.CreateTaxCalculatorFromRequest(request)
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

	response := Response{TaxLevelResponses: taxLevelResponses}
	if result.Amount < 0 {
		response.TaxRefund = -result.Amount
	} else {
		response.Tax = result.Amount
	}

	return c.JSON(http.StatusOK, response)
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
func (h *Handler) CalculateTaxCsv(c echo.Context) error {
	file, err := c.FormFile("taxes.csv")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	defer src.Close()
	reader := csv.NewReader(src)
	content, err := reader.ReadAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	var response ResponseForCSV
	for index, record := range content {
		if index == 0 {
			// TODO check column names
			continue
		}
		calculator, err := h.CreateTaxCalculatorFromCsvRecord(record)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
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

// UpdatePersonalDeductionRequest
//
//	@Summary		Update personal deduction
//	@Description	Update personal deduction
//	@Tags			tax
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Router			/admin/deductions/personal [post]
//	@Failure		500	{object}	Err
//	@Failure		400	{object}	Err
//	@Param 			UpdatePersonalDeductionRequest body UpdatePersonalDeductionRequest true "Body for update personal deduction"
func (h *Handler) UpdatePersonalDeduction(c echo.Context) error {
	var request UpdatePersonalDeductionRequest
	if err := c.Bind(&request); err != nil {
		return err
	}
	if request.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Personal deduction must be within 100,000"})
	}
	if request.Amount <= 10000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Personal deduction must be more than 10,000"})
	}
	err := h.Store.UpdateDefaultPersonalDeduction(request.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	personalDeductAmount, err := h.Store.GetDefaultPersonalDeduction()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, UpdatePersonalDeductionResponse{
		Amount: personalDeductAmount,
	})
}

// UpdateKReceipt
//
//	@Summary		Update max k-receipt deduction
//	@Description	Update max k-receipt deduction
//	@Tags			tax
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Router			/admin/deductions/k-receipt [post]
//	@Failure		500	{object}	Err
//	@Failure		400	{object}	Err
//	@Param 			UpdateKReceiptRequest body UpdateKReceiptRequest true "Body for update k-receipt deduction"
func (h *Handler) UpdateKReceipt(c echo.Context) error {
	var request UpdateKReceiptRequest
	if err := c.Bind(&request); err != nil {
		return err
	}
	if request.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "k-receipt deduction must be within 100,000"})
	}
	if request.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "k-receipt deduction must be more than 0"})
	}
	err := h.Store.UpdateMaxKReceipt(request.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	amount, err := h.Store.GetMaxKReceipt()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, UpdateKReceiptsResponse{
		Amount: amount,
	})
}
