import { SQUARE_ROOT_FUNCTION, Symbols } from "../calculatorInput";

type Props = { expression: string; error: string };

export function CalculationFeedback({ expression, error }: Props) {
  const usesRoot =
    expression.includes(SQUARE_ROOT_FUNCTION) ||
    expression.includes(Symbols.SquareRoot);
  const usesDivision =
    expression.includes("/") || expression.includes(Symbols.Divide);
  return (
    <>
      {error && (
        <p id="calculation-error" role="alert" className="error">
          {error}
        </p>
      )}
      <p id="input-help" className="help">
        Powers before × and ÷, then + and −. Use parentheses to group.
      </p>
      {expression.includes(Symbols.Percent) && (
        <p className="help">% divides the preceding value by 100.</p>
      )}
      {usesRoot && (
        <p className="help">
          Non-exact square roots round to 12 decimal places before further
          calculations.
        </p>
      )}
      {usesDivision && (
        <p className="help">
          Repeating decimals round to 12 places at the final result.
        </p>
      )}
    </>
  );
}
