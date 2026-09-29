package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

const maxBodyBytes = 4096

type calculationRequest struct {
	Expression *string `json:"expression"`
	Operation  string  `json:"operation"`
	A          *string `json:"a"`
	B          *string `json:"b"`
}

func validateEnvelope(r *http.Request) *requestError {
	if r.URL.Path != calculatePath {
		return &requestError{http.StatusNotFound, codeNotFound, "Endpoint not found."}
	}
	if r.Method != http.MethodPost {
		return &requestError{http.StatusMethodNotAllowed, codeMethodNotAllowed, "Use POST for calculations."}
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != jsonMediaType {
		return &requestError{http.StatusUnsupportedMediaType, codeUnsupportedMediaType, "Send Content-Type: application/json."}
	}
	return nil
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (*calculationRequest, *requestError) {
	body := http.MaxBytesReader(w, r.Body, maxBodyBytes)
	defer body.Close()
	input, err := decodeSingleObject(body)
	if err == nil {
		return input, nil
	}
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		return nil, &requestError{http.StatusRequestEntityTooLarge, codeRequestTooLarge, "Request must be at most 4 KB."}
	}
	return nil, &requestError{http.StatusBadRequest, codeInvalidRequest, "Send one JSON object containing an expression, or operation and string operands a and b."}
}

func decodeSingleObject(body io.Reader) (*calculationRequest, error) {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	var input *calculationRequest
	if err := decoder.Decode(&input); err != nil {
		return nil, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			err = errors.New("multiple JSON values")
		}
		return nil, err
	}
	return input, nil
}

func (input *calculationRequest) validate() *requestError {
	if input == nil {
		return missingOperation()
	}
	if input.Expression != nil {
		if input.Operation != "" || input.A != nil || input.B != nil {
			return &requestError{http.StatusBadRequest, codeInvalidRequest, "Send either an expression or separate operands, not both."}
		}
		return nil
	}
	if input.Operation == "" {
		return missingOperation()
	}
	return nil
}

func missingOperation() *requestError {
	return &requestError{http.StatusBadRequest, codeInvalidRequest, "An operation is required."}
}
