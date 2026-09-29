package httpapi

import (
	"net/http"

	"calculator/internal/calculator"
)

func NewHandler() http.Handler { return http.HandlerFunc(handleCalculation) }

func handleCalculation(w http.ResponseWriter, r *http.Request) {
	if err := validateEnvelope(r); err != nil {
		err.write(w)
		return
	}
	input, err := decodeRequest(w, r)
	if err != nil {
		err.write(w)
		return
	}
	if err := input.validate(); err != nil {
		err.write(w)
		return
	}
	result, calculationErr := evaluateRequest(input)
	if calculationErr != nil {
		fail(w, http.StatusBadRequest, calculationErr.Code, calculationErr.Message)
		return
	}
	writeJSON(w, http.StatusOK, calculationResponse{Result: result})
}

func evaluateRequest(input *calculationRequest) (string, *calculator.Error) {
	if input.Expression != nil {
		return calculator.Evaluate(*input.Expression)
	}
	return calculator.Calculate(input.Operation, input.A, input.B)
}
