# Calculator

A full-stack calculator built with Go and React + TypeScript. Supports chained addition, subtraction, multiplication, division, integer powers, square roots, and percentages, with parentheses, negative numbers, and standard operator precedence. All calculations run on the backend.

## Requirements

- Go 1.24 or newer
- Node.js 22.12 or newer (a supported LTS release is recommended)
- npm

## Run locally

In one terminal:

```sh
cd backend
go mod download
go run ./cmd/server
```

In another:

```sh
cd frontend
npm ci
npm run dev
```

Open http://127.0.0.1:5173. Vite proxies `/api` to Go at `127.0.0.1:8080`. Both processes must run. If 5173 is occupied, use the URL Vite prints.

The API address can be changed using `ADDR`, for example `ADDR=127.0.0.1:8081 go run ./cmd/server`; update `frontend/vite.config.ts` to match. No CORS configuration is needed with the development proxy.

## Run with Docker

From the repository root:

```sh
docker build -t calculator:local .
docker run --rm -p 127.0.0.1:8081:8080 calculator:local
```

Open http://127.0.0.1:8081. This runs both the production frontend and API on one port; Go and Node do not need to be installed on the host. Stop with Ctrl+C. Choose another host port if 8081 is occupied.

The multi-stage Dockerfile builds React with Node and the API with Go, then copies only the static files and binary into a non-root, shell-free image. No database, volume, or separate frontend server is required. See [Docker's multi-stage build documentation](https://docs.docker.com/build/building/multi-stage/).

To serve a production build without Docker:

```sh
cd frontend
npm ci
npm run build
cd ../backend
STATIC_DIR=../frontend/dist go run ./cmd/server
```

Open http://127.0.0.1:8080. Without `STATIC_DIR`, Go remains API-only for the Vite development workflow.

## Using the calculator

- Enter an expression using the keypad or type directly into the display. While editing, the expression occupies the main display. After equals, it moves to the smaller top line with the result below.
- `×` and `÷` take precedence over `+` and `−`: `2+3×4` is `14`.
- Parentheses override precedence: `(2+3)×4` is `20`.
- Equal-precedence addition/subtraction and multiplication/division evaluate left to right: `20÷5÷2` is `2`.
- Unary signs work with numbers and groups: `−(2+3)` is `−5` and `2×−3` is `−6`.
- Keyboard input works immediately and after keypad clicks: use digits, `.`, `+`, `-`, `*`, `/`, `^`, `%`, `sqrt(...)`, and parentheses. Press `=` or Enter to evaluate; AC or Escape clears; Backspace removes the last character. Clicking the display also allows normal caret and selection editing.
- After a result, an operator starts the next expression from the displayed result, while a digit, decimal point, or parenthesis begins a new expression. To edit an earlier expression, use the display's normal text editing controls.
- The top display shows the current operation: `3+5 =` displays `8`; pressing `+ 3 =` then shows `8 + 3` above `11`. Continuation uses the displayed decimal, including rounding: `1/3 =`, then `× 3 =` gives `0.999999999999`. Entering `1/3*3` as one expression still gives exactly `1`.

The keypad appends at the end. Keyboard editing supports insertion and selection in the display. Requests disable controls until completion or a 10-second timeout.

## Optional operations

| Operation | Examples | Rules |
| --- | --- | --- |
| Power `^` | `2^3 = 8`, `2^-3 = 0.125` | Integer exponents from -100 through 100. Reject `0^0` and zero raised to a negative exponent. |
| Square root | `sqrt(9) = 3`, `sqrt(2+7) = 3` | Parentheses required. Reject negative arguments. `√(9)` also works. |
| Percentage `%` | `10% = 0.1`, `200*10% = 20` | Postfix divide-by-100; no context-dependent percentage behavior. |

Powers associate right to left: `2^3^2 = 512`. They bind more strongly than unary signs: `-2^2 = -4`, but `(-2)^2 = 4`. Parenthesize negative exponents for readability, although `2^-3` is also accepted.

Percentage binds to the immediately preceding value or group before powers and unary signs. Thus `200+10% = 200.1`; use `200*(1+10%)` for a 10% increase. `(2^3)% = 0.08`, while `2^100% = 2^(100%) = 2`. Repeated percentages divide by 100 repeatedly.

The √ keypad button inserts `sqrt(` during entry. After a result, √ wraps that result and % appends to it; press equals to evaluate. Power continuation parenthesizes a negative displayed result, so squaring -2 yields 4. Typing `sqrt(...)` starts a new expression after a result.

Square roots of perfect rational squares are exact, including `sqrt(4/9)*3 = 2`. Other roots round to 12 fractional places **at the square-root operation**. Later arithmetic uses that approximate decimal: `sqrt(2)^2 = 1.999999999999731161391129`. The UI shows this precision rule whenever a square root is present.

## API

### POST /api/calculate

Send `Content-Type: application/json` and an expression string:

```sh
curl -i http://127.0.0.1:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"(2 + 3) * 4"}'
```

Status `200`:

```json
{"result":"20"}
```

Division-by-zero example:

```sh
curl -i http://127.0.0.1:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"1 / (2 - 2)"}'
```

Status `400`:

```json
{"error":{"code":"DIVISION_BY_ZERO","message":"Cannot divide by zero."}}
```

Advanced operations use the same expression endpoint:

```sh
curl -i http://127.0.0.1:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"sqrt(25%) + 2^3"}'
```

Response: `{"result":"8.5"}`.

The original two-operand API remains supported for compatibility:

```json
{"operation":"add","a":"0.1","b":"0.2"}
```

It accepts `add`, `subtract`, `multiply`, or `divide` and requires string operands. Send either this format or an expression, never both.

| Status | Error codes |
| --- | --- |
| 400 | `INVALID_REQUEST`, `INVALID_OPERAND`, `INVALID_EXPRESSION`, `EXPRESSION_TOO_LONG`, `UNSUPPORTED_OPERATION`, `DIVISION_BY_ZERO`, `INVALID_EXPONENT`, `EXPONENT_OUT_OF_RANGE`, `UNDEFINED_POWER`, `NEGATIVE_SQUARE_ROOT`, `CALCULATION_TOO_LARGE` |
| 404 | `NOT_FOUND` |
| 405 | `METHOD_NOT_ALLOWED`, with `Allow: POST` |
| 413 | `REQUEST_TOO_LARGE` |
| 415 | `UNSUPPORTED_MEDIA_TYPE` |

The API accepts one JSON object, rejects unknown fields and trailing JSON, and limits bodies to 4 KB. Go's standard JSON decoder semantics apply, including case-insensitive field matching and last-value handling for duplicate keys.

## Numeric contract and limits

- Expressions contain at most 256 characters, including spaces and parentheses. This bounds parsing and nesting. Each intermediate rational numerator and denominator is also limited to 4096 bits. Powers check a conservative size estimate before exponentiation; other operations check their bounded results.
- Expression literals are bounded by the overall 256-character limit, so long results can be reused without losing digits. The legacy two-operand API still limits each operand to 16 integer digits and 12 fractional digits. If a result plus the next operation exceeds 256 characters, clear and start a new calculation.
- Accept decimal forms such as `12`, `-3.5`, `.5`, and `1.`. Spaces, tabs, and line breaks separate tokens, but cannot split a number (`1 2` is invalid).
- Accept ASCII `+ - * / ^ %`, the function `sqrt(...)`, and display symbols `− × ÷ √`. Unary `+` and `-` bind more strongly than multiplication/division; repeated unary signs are valid.
- Require explicit multiplication: use `2*(3+4)`, not `2(3+4)`.
- Reject scientific notation, commas, unsupported functions/operators, infinity, `NaN`, missing operands, and unbalanced parentheses.
- Intermediate expression values are **exact rational numbers**, including divisions, except for rounded non-exact square roots. Thus `1/3*3` is `1`, and `1/(1/3*3-1)` correctly fails with division by zero.
- Final terminating decimals are returned exactly. Final repeating decimals round once to 12 fractional places, ties away from zero. Very small repeating values may round to zero. For example, `1/3` is `"0.333333333333"`, while `0.000000000001/2` is exactly `"0.0000000000005"`.
- The legacy two-operand division keeps its original rule: round to 12 fractional places, including terminating fractions with more places.
- Results use plain decimal notation with unnecessary fractional zeros removed; signed zero becomes `"0"`.

The frontend sends expression text and displays result strings without JavaScript arithmetic or numeric conversion. The backend uses `shopspring/decimal` to validate decimal literals and for the original two-operand API. Expression evaluation uses Go's `math/big.Rat` for exact intermediate values and its decimal conversion to format final results.

## Design

- **Parser:** a shunting-yard evaluator with separate operator and value stacks. Operator precedence controls when values are reduced; parentheses defer reduction. No `eval`, execution of user code, recursion, or external parsing dependency.
- **Backend:** standard-library `net/http`, a small JSON handler, and calculation functions independent of HTTP. No database, framework, or service interfaces.
- **Frontend:** React local state, native fetch, plain CSS, Vite. Client-side checks cover supported characters and expression limits; the backend owns grammar and mathematical validation.
- **Visual design:** the interface uses Sezzle's published deep-purple, lilac, mint, peach, and coral palette. The calculator remains product-focused and does not reproduce or alter Sezzle logo assets.
- **Square-root algorithm:** integer square roots identify exact rational roots. For irrational roots, integer comparisons round the scaled root without any binary floating-point conversion.
- **Scope:** history, persistence, authentication, and fractional exponents remain out of scope.

### Code navigation

| Responsibility | Location |
| --- | --- |
| Token scanning and precedence | `backend/internal/calculator/tokens.go` |
| Grammar and stack reduction | `backend/internal/calculator/expression.go` |
| Arithmetic and formatting | `operators.go`, `advanced.go`, `format.go` in the calculator package |
| Named limits and error codes | `limits.go`, `errors.go` in the calculator package |
| HTTP handling, validation, responses | `backend/internal/httpapi/handler.go`, `request.go`, `responses.go` |
| UI composition and presentation | `frontend/src/App.tsx`, `components/` |
| State and keyboard events | `frontend/src/hooks/` |
| Pure input rules and typed commands | `frontend/src/calculatorInput.ts` |
| Transport and client validation | `frontend/src/api.ts`, `validation.ts` |

Presentation components receive data and callbacks; input rules do not depend on React or HTTP. Named tokens, commands, error codes, and limits replace internal marker strings and repeated policy values. Small functions and modules provide these boundaries without service classes or interfaces with only one implementation.

## Checks and coverage

Backend:

```sh
cd backend
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
go vet ./...
go test ./internal/calculator -run='^$' -fuzz=FuzzEvaluate -fuzztime=5s
```

Frontend:

```sh
cd frontend
npm test
npm run coverage
npm run build
npm run format:check
```

The coverage summary is in [docs/coverage.md](docs/coverage.md). Detailed reports are generated at `backend/coverage.html` and `frontend/coverage/index.html` and are Git-ignored.

`npm run build` checks TypeScript and writes frontend assets to `frontend/dist`. Set `STATIC_DIR` to that directory to serve the build with Go. The Docker image configures this automatically. Static files support GET/HEAD; unknown files and API routes return 404, and directory listings are disabled.

## AI assistance

Codex assisted with planning, implementation, and verification. See [PROMPTS.md](PROMPTS.md).
