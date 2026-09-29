package httpapi

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	valid := `{"operation":"add","a":"0.1","b":"0.2"}`
	tests := []struct {
		name, method, path, contentType, body string
		status                                int
		code, result                          string
	}{
		{"success", "POST", "/api/calculate", "application/json", valid, 200, "", "0.3"},
		{"charset", "POST", "/api/calculate", "application/json; charset=utf-8", valid, 200, "", "0.3"},
		{"subtract", "POST", "/api/calculate", "application/json", `{"operation":"subtract","a":"2","b":"3"}`, 200, "", "-1"},
		{"multiply", "POST", "/api/calculate", "application/json", `{"operation":"multiply","a":"2","b":"3"}`, 200, "", "6"},
		{"divide", "POST", "/api/calculate", "application/json", `{"operation":"divide","a":"2","b":"4"}`, 200, "", "0.5"},
		{"zero divisor", "POST", "/api/calculate", "application/json", `{"operation":"divide","a":"1","b":"-0.000"}`, 400, "DIVISION_BY_ZERO", ""},
		{"unsupported", "POST", "/api/calculate", "application/json", `{"operation":"sqrt","a":"9"}`, 400, "UNSUPPORTED_OPERATION", ""},
		{"missing b", "POST", "/api/calculate", "application/json", `{"operation":"add","a":"1"}`, 400, "INVALID_OPERAND", ""},
		{"null a", "POST", "/api/calculate", "application/json", `{"operation":"add","a":null,"b":"1"}`, 400, "INVALID_OPERAND", ""},
		{"invalid decimal", "POST", "/api/calculate", "application/json", `{"operation":"add","a":"NaN","b":"1"}`, 400, "INVALID_OPERAND", ""},
		{"numeric operand", "POST", "/api/calculate", "application/json", `{"operation":"add","a":1,"b":"1"}`, 400, "INVALID_REQUEST", ""},
		{"unknown field", "POST", "/api/calculate", "application/json", `{"operation":"add","a":"1","b":"1","extra":1}`, 400, "INVALID_REQUEST", ""},
		{"missing operation", "POST", "/api/calculate", "application/json", `{"a":"1","b":"1"}`, 400, "INVALID_REQUEST", ""},
		{"null body", "POST", "/api/calculate", "application/json", `null`, 400, "INVALID_REQUEST", ""},
		{"array body", "POST", "/api/calculate", "application/json", `[]`, 400, "INVALID_REQUEST", ""},
		{"empty body", "POST", "/api/calculate", "application/json", ``, 400, "INVALID_REQUEST", ""},
		{"malformed", "POST", "/api/calculate", "application/json", `{`, 400, "INVALID_REQUEST", ""},
		{"trailing value", "POST", "/api/calculate", "application/json", valid + ` {}`, 400, "INVALID_REQUEST", ""},
		{"trailing garbage", "POST", "/api/calculate", "application/json", valid + ` !`, 400, "INVALID_REQUEST", ""},
		{"trailing whitespace", "POST", "/api/calculate", "application/json", valid + " \n", 200, "", "0.3"},
		{"large body", "POST", "/api/calculate", "application/json", `{"a":"` + strings.Repeat("1", 4096) + `"}`, 413, "REQUEST_TOO_LARGE", ""},
		{"large trailing whitespace", "POST", "/api/calculate", "application/json", valid + strings.Repeat(" ", 4096), 413, "REQUEST_TOO_LARGE", ""},
		{"wrong method", "GET", "/api/calculate", "application/json", "", 405, "METHOD_NOT_ALLOWED", ""},
		{"wrong media", "POST", "/api/calculate", "text/plain", valid, 415, "UNSUPPORTED_MEDIA_TYPE", ""},
		{"missing media", "POST", "/api/calculate", "", valid, 415, "UNSUPPORTED_MEDIA_TYPE", ""},
		{"unknown route", "POST", "/api/other", "application/json", valid, 404, "NOT_FOUND", ""},
	}
	handler := NewHandler()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, req)
			if response.Code != tt.status {
				t.Fatalf("got %d: %s", response.Code, response.Body)
			}
			if response.Header().Get("Content-Type") != "application/json" {
				t.Fatal("expected JSON")
			}
			if tt.status == 405 && response.Header().Get("Allow") != "POST" {
				t.Fatal("missing Allow header")
			}
			var body struct {
				Result string `json:"result"`
				Error  struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Result != tt.result || body.Error.Code != tt.code {
				t.Fatalf("unexpected response: %s", response.Body)
			}
			if tt.status != 200 && body.Error.Message == "" {
				t.Fatal("missing error message")
			}
		})
	}
}

func TestExpressionRequests(t *testing.T) {
	tests := []struct {
		body         string
		status       int
		result, code string
	}{
		{`{"expression":"2+3*4"}`, 200, "14", ""},
		{`{"expression":"(2+3)*4"}`, 200, "20", ""},
		{`{"expression":"1/3*3"}`, 200, "1", ""},
		{`{"expression":"1/(2-2)"}`, 400, "", "DIVISION_BY_ZERO"},
		{`{"expression":"1+"}`, 400, "", "INVALID_EXPRESSION"},
		{`{"expression":""}`, 400, "", "INVALID_EXPRESSION"},
		{`{"expression":123}`, 400, "", "INVALID_REQUEST"},
		{`{"expression":null}`, 400, "", "INVALID_REQUEST"},
		{`{"expression":"1+2","operation":"add","a":"1","b":"2"}`, 400, "", "INVALID_REQUEST"},
		{`{"expression":"` + strings.Repeat("1", 257) + `"}`, 400, "", "EXPRESSION_TOO_LONG"},
	}
	for _, tt := range tests {
		t.Run(tt.body, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			NewHandler().ServeHTTP(response, req)
			var body struct {
				Result string
				Error  struct{ Code string }
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if response.Code != tt.status || body.Result != tt.result || body.Error.Code != tt.code {
				t.Fatalf("got %d: %s", response.Code, response.Body)
			}
		})
	}
}

func TestAdvancedRequests(t *testing.T) {
	for _, tt := range []struct{ expression, result, code string }{
		{"2^3^2", "512", ""}, {"sqrt(2)", "1.414213562373", ""},
		{"200*10%", "20", ""}, {"sqrt(-1)", "", "NEGATIVE_SQUARE_ROOT"},
		{"0^0", "", "UNDEFINED_POWER"}, {"2^0.5", "", "INVALID_EXPONENT"},
		{"2^101", "", "EXPONENT_OUT_OF_RANGE"}, {"(2^100)^100", "", "CALCULATION_TOO_LARGE"},
	} {
		t.Run(tt.expression, func(t *testing.T) {
			payload, err := json.Marshal(map[string]string{"expression": tt.expression})
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(string(payload)))
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			NewHandler().ServeHTTP(response, req)
			var body struct {
				Result string
				Error  struct{ Code string }
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			status := 200
			if tt.code != "" {
				status = 400
			}
			if response.Code != status || body.Result != tt.result || body.Error.Code != tt.code {
				t.Fatalf("got %d: %s", response.Code, response.Body)
			}
		})
	}
}
