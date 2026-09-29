import { Command, Symbols, type CalculatorCommand } from "../calculatorInput";

type Key = {
  text: string;
  label: string;
  kind: "utility" | "operator" | "number" | "zero" | "equals";
  command: CalculatorCommand;
};

function insertKey(
  text: string,
  label = text,
  kind: Key["kind"] = "number",
): Key {
  return { text, label, kind, command: { type: Command.Insert, value: text } };
}

const keys: Key[] = [
  {
    text: "AC",
    label: "Clear",
    kind: "utility",
    command: { type: Command.Clear },
  },
  {
    text: Symbols.SquareRoot,
    label: "Square root",
    kind: "utility",
    command: { type: Command.SquareRoot },
  },
  insertKey(Symbols.Percent, "Percentage", "utility"),
  {
    text: "⌫",
    label: "Backspace",
    kind: "utility",
    command: { type: Command.Backspace },
  },
  insertKey(Symbols.OpenGroup, "Open parenthesis", "utility"),
  insertKey(Symbols.CloseGroup, "Close parenthesis", "utility"),
  insertKey(Symbols.Power, "Power", "operator"),
  insertKey(Symbols.Divide, "Divide", "operator"),
  insertKey("7"),
  insertKey("8"),
  insertKey("9"),
  insertKey(Symbols.Multiply, "Multiply", "operator"),
  insertKey("4"),
  insertKey("5"),
  insertKey("6"),
  insertKey(Symbols.Subtract, "Subtract", "operator"),
  insertKey("1"),
  insertKey("2"),
  insertKey("3"),
  insertKey(Symbols.Add, "Add", "operator"),
  insertKey("0", "0", "zero"),
  insertKey(".", "Decimal point"),
  {
    text: "=",
    label: "Calculate",
    kind: "equals",
    command: { type: Command.Evaluate },
  },
];

export function CalculatorKeypad({
  onCommand,
}: {
  onCommand: (command: CalculatorCommand) => void;
}) {
  return (
    <div className="keypad">
      {keys.map((key) => (
        <button
          key={key.label}
          type={key.command.type === Command.Evaluate ? "submit" : "button"}
          className={`key ${key.kind}`}
          aria-label={key.label}
          onClick={
            key.command.type === Command.Evaluate
              ? undefined
              : () => onCommand(key.command)
          }
        >
          {key.text}
        </button>
      ))}
    </div>
  );
}
