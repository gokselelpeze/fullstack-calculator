export const Symbols = {
  Add: "+",
  Subtract: "−",
  Multiply: "×",
  Divide: "÷",
  Power: "^",
  Percent: "%",
  SquareRoot: "√",
  OpenGroup: "(",
  CloseGroup: ")",
} as const;

export const SQUARE_ROOT_FUNCTION = "sqrt";

export const Command = {
  Insert: "insert",
  Clear: "clear",
  Backspace: "backspace",
  SquareRoot: "square-root",
  Evaluate: "evaluate",
} as const;

export type EntryCommand =
  | { type: typeof Command.Insert; value: string }
  | {
      type:
        | typeof Command.Clear
        | typeof Command.Backspace
        | typeof Command.SquareRoot;
    };
export type CalculatorCommand =
  EntryCommand | { type: typeof Command.Evaluate };

const binaryOperators = new Set<string>([
  Symbols.Add,
  Symbols.Subtract,
  Symbols.Multiply,
  Symbols.Divide,
  Symbols.Power,
]);
const keyboardInputPattern = /^[0-9.()+*/−×÷^%sqrt-]$/;

export function keyboardCommand(key: string): CalculatorCommand | null {
  switch (key) {
    case "Enter":
    case "=":
      return { type: Command.Evaluate };
    case "Escape":
      return { type: Command.Clear };
    case "Backspace":
      return { type: Command.Backspace };
    case Symbols.SquareRoot:
      return { type: Command.SquareRoot };
    default:
      if (!keyboardInputPattern.test(key)) return null;
      return { type: Command.Insert, value: normalizeOperator(key) };
  }
}

function normalizeOperator(key: string): string {
  switch (key) {
    case "*":
      return Symbols.Multiply;
    case "/":
      return Symbols.Divide;
    case "-":
      return Symbols.Subtract;
    default:
      return key;
  }
}

// Pure entry rules, shared by the keypad and keyboard. No React or API access.
export function nextExpression(
  expression: string,
  result: string | null,
  command: EntryCommand,
): string {
  switch (command.type) {
    case Command.Clear:
      return "";
    case Command.Backspace:
      return expression.slice(0, -1);
    case Command.SquareRoot:
      return result !== null
        ? `${SQUARE_ROOT_FUNCTION}(${result})`
        : expression + `${SQUARE_ROOT_FUNCTION}(`;
    case Command.Insert:
      return result !== null
        ? continueFromResult(result, command.value)
        : expression + command.value;
  }
}

function continueFromResult(result: string, value: string): string {
  if (value === Symbols.Percent) return `(${result})${Symbols.Percent}`;
  if (!binaryOperators.has(value)) return value;
  const base =
    value === Symbols.Power && result.startsWith("-") ? `(${result})` : result;
  return `${base} ${value} `;
}
