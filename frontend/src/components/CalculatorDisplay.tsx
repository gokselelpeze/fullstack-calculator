import type { KeyboardEventHandler } from "react";

type Props = {
  expression: string;
  result: string | null;
  loading: boolean;
  error: string;
  onEdit: (value: string) => void;
  onKeyDown: KeyboardEventHandler<HTMLTextAreaElement>;
};

export function CalculatorDisplay({
  expression,
  result,
  loading,
  error,
  onEdit,
  onKeyDown,
}: Props) {
  const showResult = result !== null || loading;
  return (
    <div className={`display ${showResult ? "has-result" : ""}`}>
      <label className="sr-only" htmlFor="expression">
        Expression
      </label>
      <textarea
        id="expression"
        rows={1}
        placeholder="0"
        value={expression}
        spellCheck={false}
        autoComplete="off"
        aria-invalid={!!error}
        aria-describedby={error ? "calculation-error" : undefined}
        onChange={(event) => onEdit(event.target.value)}
        onKeyDown={onKeyDown}
      />
      <div
        className={`result ${showResult ? "" : "sr-only"}`}
        role="status"
        aria-live="polite"
        aria-atomic="true"
      >
        {showResult && <span className="sr-only">Result: </span>}
        {loading ? <span className="pending">Calculating…</span> : result}
      </div>
    </div>
  );
}
