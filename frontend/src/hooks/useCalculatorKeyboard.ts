import {
  useCallback,
  useEffect,
  type KeyboardEvent as ReactKeyboardEvent,
} from "react";
import {
  Command,
  keyboardCommand,
  type CalculatorCommand,
} from "../calculatorInput";

const editableSelector = "input, textarea, select, [contenteditable='true']";

type Options = {
  loading: boolean;
  hasResult: boolean;
  execute: (command: CalculatorCommand) => void;
};

function ignoreShortcut(event: KeyboardEvent): boolean {
  return (
    event.defaultPrevented ||
    event.isComposing ||
    event.ctrlKey ||
    event.metaKey ||
    event.altKey
  );
}

export function useCalculatorKeyboard({
  loading,
  hasResult,
  execute,
}: Options) {
  useEffect(() => {
    function shortcut(event: KeyboardEvent) {
      if (loading || ignoreShortcut(event)) return;
      if (
        event.target instanceof HTMLElement &&
        event.target.closest(editableSelector)
      )
        return;
      const command = keyboardCommand(event.key);
      if (!command) return;
      event.preventDefault();
      execute(command);
    }
    window.addEventListener("keydown", shortcut);
    return () => window.removeEventListener("keydown", shortcut);
  }, [loading, execute]);

  return useCallback(
    (event: ReactKeyboardEvent<HTMLTextAreaElement>) => {
      if (loading || ignoreShortcut(event.nativeEvent)) return;
      const command = keyboardCommand(event.key);
      if (!command || command.type === Command.Backspace) return;
      // While editing, native input preserves caret/selection and paste behavior.
      if (command.type === Command.Insert && !hasResult) return;
      event.preventDefault();
      execute(command);
    },
    [loading, hasResult, execute],
  );
}
