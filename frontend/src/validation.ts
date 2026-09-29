import { SQUARE_ROOT_FUNCTION, Symbols } from "./calculatorInput";

const LEGACY_INTEGER_DIGITS = 16;
const LEGACY_FRACTIONAL_DIGITS = 12;
const decimalMessage = "Use a decimal number, such as -3.5.";
const signedDecimalPattern = /^[+-]?(\d+(\.\d*)?|\.\d+)$/;
const decimalLiteralPattern = /^(\d+(\.\d*)?|\.\d+)$/;
const expressionCharactersPattern = /^[0-9.+\-*/()×÷−^% \t\r\n]*$/;

export function validateOperand(value: string): string | null {
  const text = value.trim();
  if (!text) return "Enter a number.";
  if (!signedDecimalPattern.test(text)) return decimalMessage;
  const [integer, fraction = ""] = text.replace(/^[+-]/, "").split(".");
  if (
    integer.length > LEGACY_INTEGER_DIGITS ||
    fraction.length > LEGACY_FRACTIONAL_DIGITS
  )
    return `Use at most ${LEGACY_INTEGER_DIGITS} integer digits and ${LEGACY_FRACTIONAL_DIGITS} decimal places.`;
  return null;
}

export const MAX_EXPRESSION_LENGTH = 256;

// Full grammar and mathematical validation belong to the backend.
export function validateExpression(value: string): string | null {
  if (!value.trim()) return "Enter an expression.";
  if (value.length > MAX_EXPRESSION_LENGTH)
    return `Use at most ${MAX_EXPRESSION_LENGTH} characters.`;
  const withoutFunctions = value
    .replaceAll(SQUARE_ROOT_FUNCTION, "")
    .replaceAll(Symbols.SquareRoot, "");
  if (!expressionCharactersPattern.test(withoutFunctions))
    return "Use numbers, +, −, ×, ÷, ^, %, sqrt(), and parentheses.";
  for (const literal of value.match(/[0-9.]+/g) ?? []) {
    if (!decimalLiteralPattern.test(literal)) return decimalMessage;
  }
  return null;
}
