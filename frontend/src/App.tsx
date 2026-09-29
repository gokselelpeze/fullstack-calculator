import { Command } from "./calculatorInput";
import { CalculationFeedback } from "./components/CalculationFeedback";
import { CalculatorDisplay } from "./components/CalculatorDisplay";
import { CalculatorKeypad } from "./components/CalculatorKeypad";
import { useCalculator } from "./hooks/useCalculator";
import { useCalculatorKeyboard } from "./hooks/useCalculatorKeyboard";

export default function App() {
  const calculator = useCalculator();
  const keyboard = useCalculatorKeyboard({
    loading: calculator.loading,
    hasResult: calculator.result !== null,
    execute: calculator.execute,
  });

  return (
    <main className="page">
      <section className="calculator" aria-label="Calculator">
        <header className="card-heading">
          <h1>Calculator</h1>
        </header>
        <form
          noValidate
          onSubmit={(event) => {
            event.preventDefault();
            calculator.execute({ type: Command.Evaluate });
          }}
        >
          <fieldset disabled={calculator.loading}>
            <CalculatorDisplay
              expression={calculator.expression}
              result={calculator.result}
              error={calculator.error}
              onEdit={calculator.edit}
              onKeyDown={keyboard}
            />
            <CalculatorKeypad onCommand={calculator.execute} />
          </fieldset>
        </form>
        <CalculationFeedback
          expression={calculator.expression}
          error={calculator.error}
        />
      </section>
    </main>
  );
}
