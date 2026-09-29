# Verification

Verified on 2026-09-29 with Go 1.27.1, Node.js 25.2.1, and npm 11.6.2 on macOS arm64. Docker verification used Linux arm64, Go 1.27.1, and Node 24 from the Dockerfile. Go runs from a temporary tool directory on the host.

## Backend

`go test -race -coverprofile=coverage.out ./...` and `go vet ./...` passed.

| Package | Statement coverage |
| --- | ---: |
| internal/calculator | 96.4% |
| internal/httpapi | 100.0% |
| cmd/server | 0.0% |
| Total | 94.8% |

The advanced arithmetic helpers (power, square root, size checks) each have 100% statement coverage. Tests cover:

- Basic arithmetic, exact decimals, large values, rational intermediates, and formatting.
- Powers: right associativity, unary-minus precedence, integer-only and bounded exponents, negative powers, zero-domain errors, and size guards before exponentiation.
- Square roots: perfect rational squares, nested roots, irrational rounding up/down, very small values, negative arguments, and approximate-root composition.
- Percentages: literals, groups, repeated percent, precedence, and combinations with powers and roots.
- Parser syntax, nesting, expression limits, computed zero divisors, and API error/status contracts.
- Production serving: index and assets, HEAD requests, API routing, missing paths, method checks, directory-listing prevention, and API-only mode.

Uncovered calculator paths primarily cover defensive stack and unknown-token checks. The server entry point is exercised through the live application rather than unit coverage.

A five-second fuzz run including the new operators passed **635,229 executions** without a crash or empty successful result. This supplements semantic tests, not a proof of mathematical correctness.

## Frontend

`npm run coverage`: **81 tests passed** across three files.

| Metric | Coverage |
| --- | ---: |
| Statements | 96.22% |
| Branches | 94.92% |
| Functions | 100.00% |
| Lines | 99.29% |

Coverage includes all application TypeScript/TSX files, including extracted components, hooks, and input helpers. DOM bootstrap, CSS, and configuration are excluded. Tests cover keypad and global keyboard entry, optional operations, result continuation (including a negative power base), display transitions, text editing, modifier shortcuts, clearing/backspace, loading, API/network errors, timeouts, input limits, and malformed responses.

TypeScript checking, the production build, and formatting checks passed.

## Docker and live browser

A post-refactor compatibility check against the previous Docker API compared 336 expressions, including invalid input and domain errors. Every HTTP status and JSON result/error matched. The deployment and browser checks below were performed before the structural refactor.

`docker build -t calculator:local .` succeeded. The image runs the compiled Go binary as user 65532 and serves the production frontend and API from the same port. The build's dependency installation reported zero vulnerabilities.

Started the image on localhost:8081 and verified in the browser:

- Keypad entry `sqrt(25%)+2^3` returned `8.5`.
- `2^3^2` returned `512`.
- At 390 px width, `sqrt(-1)` displayed the expected domain error; correcting the expression recovered normally.
- The production JavaScript/CSS loaded through Go, and API requests reached the bundled backend.
- Desktop/mobile keypad layout was visually inspected; temporary viewport overrides were reset.

Earlier checks also verified basic arithmetic, rounding, calculated zero divisors, and `3+5 = 8` followed by `+3 = 11` through mixed mouse/keyboard input.

These are manual browser smoke checks, not a committed end-to-end suite. Native mobile keyboards and multiple browser engines were not tested.

## Regenerate

See README.md for commands. Generated reports are Git-ignored:

- backend/coverage.out
- backend/coverage.html
- frontend/coverage/index.html
- frontend/coverage/coverage-summary.json
