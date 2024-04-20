package tax

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestTaxHandler(t *testing.T) {
	t.Run("given request with total income 500000.0 should return 200 and response with tax 29000.0", func(t *testing.T) {
		body, err := json.Marshal(CalculationRequest{
			TotalIncome:    500000.0,
			WithHoldingTax: 0,
			Allowances: []Allowance{
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

		CalculateTax(c)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := Response{Tax: 29000.0}
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
			Allowances: []Allowance{
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

		CalculateTax(c)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := Response{Tax: 4000.0}
		var got Response
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}
