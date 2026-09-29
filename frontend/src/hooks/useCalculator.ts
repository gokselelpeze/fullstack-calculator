import { useCallback, useState } from "react";
import { calculate, CALCULATION_FAILED_MESSAGE } from "../api";
import {
  Command,
  nextExpression,
  type CalculatorCommand,
} from "../calculatorInput";
import { MAX_EXPRESSION_LENGTH, validateExpression } from "../validation";

const continuationLimitMessage = `Use at most ${MAX_EXPRESSION_LENGTH} characters. Clear to start a new calculation.`;

export function useCalculator() {
  const [expression, setExpression] = useState("");
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const edit = useCallback((value: string) => {
    setResult(null);
    setError("");
    setExpression(value);
  }, []);

  const evaluate = useCallback(async () => {
    if (loading) return;
    const validation = validateExpression(expression);
    setError(validation ?? "");
    if (validation) return;
    setLoading(true);
    try {
      setResult(await calculate({ expression }));
    } catch (err) {
      setResult(null);
      setError(err instanceof Error ? err.message : CALCULATION_FAILED_MESSAGE);
    } finally {
      setLoading(false);
    }
  }, [expression, loading]);

  const execute = useCallback(
    (command: CalculatorCommand) => {
      if (loading) return;
      if (command.type === Command.Evaluate) {
        void evaluate();
        return;
      }
      const next = nextExpression(expression, result, command);
      // Clearing and deletion must remain available even after oversized pasted input.
      const addsInput =
        command.type === Command.Insert || command.type === Command.SquareRoot;
      if (addsInput && next.length > MAX_EXPRESSION_LENGTH) {
        setError(continuationLimitMessage);
        return;
      }
      edit(next);
    },
    [edit, evaluate, expression, loading, result],
  );

  return { expression, result, error, loading, edit, execute };
}
