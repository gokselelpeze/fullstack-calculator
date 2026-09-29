import { describe, expect, it } from "vitest";
import { validateOperand } from "./validation";

describe("decimal input", () => {
  it.each([
    "0",
    "-0",
    "+2",
    ".5",
    "-.5",
    "1.",
    " 1.20 ",
    "9999999999999999.123456789012",
  ])("accepts %s", (value) => {
    expect(validateOperand(value)).toBeNull();
  });
  it.each([
    "",
    " ",
    ".",
    "-",
    "1e3",
    "NaN",
    "Infinity",
    "1,000",
    "1 2",
    "1.2.3",
    "10000000000000000",
    "0.1234567890123",
  ])("rejects %s", (value) => {
    expect(validateOperand(value)).not.toBeNull();
  });
});

import { validateExpression } from "./validation";
it.each([
  "2+3*4",
  "-(2+.5)",
  "1 ÷ 3 × 3",
  "1+",
  "(",
  "99999999999999980000000000000001 + 1",
  "0.000000000000000000000001 * 2",
])(
  "accepts supported expression characters for backend parsing: %s",
  (expression) => {
    expect(validateExpression(expression)).toBeNull();
  },
);
it.each(["", " ", "1e3", "1,000", "1..2", "1".repeat(257)])(
  "rejects invalid expression input: %s",
  (expression) => {
    expect(validateExpression(expression)).not.toBeNull();
  },
);

it.each(["2^3", "sqrt(9)", "√(9)", "200*10%", "sqrt(25%)+2^3"])(
  "accepts optional operation %s",
  (expression) => {
    expect(validateExpression(expression)).toBeNull();
  },
);
it.each(["sin(2)", "sqr(9)", "sqrt(1e3)"])(
  "rejects unsupported function input %s",
  (expression) => {
    expect(validateExpression(expression)).not.toBeNull();
  },
);
