package tax

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

type MockStore struct {
	PersonalDeductionAmount float64
}

func (m *MockStore) UpdateDefaultPersonalDeduction(value float64) error {
	m.PersonalDeductionAmount = value
	return nil
}

func (m *MockStore) GetDefaultPersonalDeduction() (float64, error) {
	return m.PersonalDeductionAmount, nil
}

func NewMockStore() *MockStore {
	return &MockStore{PersonalDeductionAmount: 60000.00}
}

func TestTaxHandler(t *testing.T) {
	t.Run("given request with total income 500000.0 should return 200 and response with tax 29000.0", func(t *testing.T) {
		body, err := json.Marshal(CalculationRequest{
			TotalIncome:    500000.0,
			WithHoldingTax: 0,
			Allowances: []AllowanceRequest{
				{Type: "donation", Amount: 0.0},
			},
		})
		if err != nil {
			t.Errorf("Unable to create body request, error: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)

		handler := Handler{Store: NewMockStore()}
		handler.CalculateTax(c)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := Response{
			Tax: 29000.0,
			TaxLevelResponses: []TaxLevelResponse{
				{"0 - 150,000", 0.00},
				{"150,001 - 500,000", 29000.00},
				{"500,001 - 1,000,000", 0.00},
				{"1,000,001 - 2,000,000", 0.00},
				{"2,000,001 ขึ้นไป", 0.00},
			},
		}
		var got Response
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given request with total income 500000.0 and wht 25000.0 should return 200 and response with tax 4000.0", func(t *testing.T) {
		body, err := json.Marshal(CalculationRequest{
			TotalIncome:    500000.0,
			WithHoldingTax: 25000.0,
			Allowances: []AllowanceRequest{
				{Type: "donation", Amount: 0.0},
			},
		})
		if err != nil {
			t.Errorf("Unable to create body request, error: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)

		handler := Handler{Store: NewMockStore()}
		handler.CalculateTax(c)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := Response{
			Tax: 4000.0,
			TaxLevelResponses: []TaxLevelResponse{
				{"0 - 150,000", 0.00},
				{"150,001 - 500,000", 29000.00},
				{"500,001 - 1,000,000", 0.00},
				{"1,000,001 - 2,000,000", 0.00},
				{"2,000,001 ขึ้นไป", 0.00},
			},
		}
		var got Response
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given request with total income 500000.0 donation 200000.0 should return 200 and response with tax 19000.0", func(t *testing.T) {
		body, err := json.Marshal(CalculationRequest{
			TotalIncome:    500000.0,
			WithHoldingTax: 0,
			Allowances: []AllowanceRequest{
				{Type: "donation", Amount: 200000.0},
			},
		})
		if err != nil {
			t.Errorf("Unable to create body request, error: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)

		handler := Handler{Store: NewMockStore()}
		handler.CalculateTax(c)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := Response{
			Tax: 19000.0,
			TaxLevelResponses: []TaxLevelResponse{
				{"0 - 150,000", 0.00},
				{"150,001 - 500,000", 19000.00},
				{"500,001 - 1,000,000", 0.00},
				{"1,000,001 - 2,000,000", 0.00},
				{"2,000,001 ขึ้นไป", 0.00},
			},
		}
		var got Response
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given request with total income 500000.0 k-receipt 200000.0 donation 100000.0 should return 200 and response with tax 14000.0", func(t *testing.T) {
		body, err := json.Marshal(CalculationRequest{
			TotalIncome:    500000.0,
			WithHoldingTax: 0,
			Allowances: []AllowanceRequest{
				{"k-receipt", 200000.0},
				{"donation", 100000.0},
			},
		})
		if err != nil {
			t.Errorf("Unable to create body request, error: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)

		handler := Handler{Store: NewMockStore()}
		handler.CalculateTax(c)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := Response{
			Tax: 14000.0,
			TaxLevelResponses: []TaxLevelResponse{
				{"0 - 150,000", 0.00},
				{"150,001 - 500,000", 14000.0},
				{"500,001 - 1,000,000", 0.00},
				{"1,000,001 - 2,000,000", 0.00},
				{"2,000,001 ขึ้นไป", 0.00},
			},
		}
		var got Response
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given request with total income 100000.0 wht 1000.0 should return 200 and response with tax refund 1000.0", func(t *testing.T) {
		body, err := json.Marshal(CalculationRequest{
			TotalIncome:    100000.0,
			WithHoldingTax: 1000,
			Allowances:     []AllowanceRequest{},
		})
		if err != nil {
			t.Errorf("Unable to create body request, error: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)

		handler := Handler{Store: NewMockStore()}
		handler.CalculateTax(c)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := Response{
			TaxRefund: 1000.0,
			TaxLevelResponses: []TaxLevelResponse{
				{"0 - 150,000", 0.00},
				{"150,001 - 500,000", 0.0},
				{"500,001 - 1,000,000", 0.00},
				{"1,000,001 - 2,000,000", 0.00},
				{"2,000,001 ขึ้นไป", 0.00},
			},
		}
		var got Response
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given request with CSV file should return 200 and response with tax info", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		dataPart, err := writer.CreateFormFile("taxes.csv", "taxes.csv")
		if err != nil {
			t.Errorf("Unable to create form file with error: %v", err)
		}
		f, err := os.Open("../test.csv")
		if err != nil {
			t.Errorf("Unable to open file with error: %v", err)
		}
		_, err = io.Copy(dataPart, f)
		if err != nil {
			t.Errorf("Unable to copy file with error: %v", err)
		}
		err = writer.Close()
		if err != nil {
			t.Errorf("Unable to close writer after write body request, error : %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)

		handler := Handler{Store: NewMockStore()}
		handler.CalculateTaxCsv(c)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := ResponseForCSV{
			Taxes: []ResponseTaxResultForCSV{
				{TotalIncome: 500000.00, Tax: 29000},
				{TotalIncome: 600000.00, TaxRefund: 3000},
				{TotalIncome: 750000.00, Tax: 3750},
			},
		}
		var got ResponseForCSV
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given request update personal deduction 29000.0 should return 200 and response with personal deduction amount 29000.0", func(t *testing.T) {
		body, err := json.Marshal(UpdatePersonalDeductionRequest{29000.00})
		if err != nil {
			t.Errorf("Unable to create body request, error: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)

		handler := Handler{Store: NewMockStore()}
		handler.UpdatePersonalDeduction(c)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := UpdatePersonalDeductionResponse{29000.00}
		var got UpdatePersonalDeductionResponse
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}
