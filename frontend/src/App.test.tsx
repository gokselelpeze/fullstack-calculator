import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, expect, it, vi } from "vitest";
import App from "./App";

afterEach(() => vi.unstubAllGlobals());
function fill(expression: string) {
  fireEvent.change(screen.getByLabelText("Expression"), {
    target: { value: expression },
  });
}
function mockResult(result: string) {
  const fetch = vi
    .fn()
    .mockImplementation(() =>
      Promise.resolve(new Response(JSON.stringify({ result }))),
    );
  vi.stubGlobal("fetch", fetch);
  return fetch;
}
async function click(name: string) {
  await userEvent.click(screen.getByRole("button", { name }));
}
it("builds a chained expression with the keypad and sends it to Go", async () => {
  const fetch = mockResult("14");
  render(<App />);
  for (const key of ["2", "Add", "3", "Multiply", "4"]) await click(key);
  expect(screen.getByLabelText("Expression")).toHaveValue("2+3×4");
  await click("Calculate");
  expect(await screen.findByText("14")).toBeInTheDocument();
  expect(JSON.parse(fetch.mock.calls[0][1].body)).toEqual({
    expression: "2+3×4",
  });
});
it("supports parentheses, decimals, and negative numbers", async () => {
  mockResult("-1.5");
  render(<App />);
  for (const key of [
    "Subtract",
    "Open parenthesis",
    "1",
    "Add",
    "Decimal point",
    "5",
    "Close parenthesis",
  ])
    await click(key);
  expect(screen.getByLabelText("Expression")).toHaveValue("−(1+.5)");
  await click("Calculate");
  expect(await screen.findByText("-1.5")).toBeInTheDocument();
});
it("submits typed expressions with Enter and clears stale results on edit", async () => {
  mockResult("20");
  render(<App />);
  fill("(2+3)*4");
  await userEvent.click(screen.getByLabelText("Expression"));
  await userEvent.keyboard("{Enter}");
  expect(await screen.findByText("20")).toBeInTheDocument();
  fill("2+3");
  expect(screen.queryByText("20")).not.toBeInTheDocument();
});
it("continues from the displayed result and replaces the previous operation", async () => {
  const fetch = vi
    .fn()
    .mockResolvedValueOnce(new Response('{"result":"8"}'))
    .mockResolvedValueOnce(new Response('{"result":"11"}'));
  vi.stubGlobal("fetch", fetch);
  render(<App />);
  fill("3+5");
  await click("Calculate");
  expect(await screen.findByRole("status")).toHaveTextContent("8");
  expect(screen.getByLabelText("Expression")).toHaveValue("3+5");
  await click("Add");
  await click("3");
  expect(screen.getByLabelText("Expression")).toHaveValue("8 + 3");
  await click("Calculate");
  await screen.findByText("11");
  expect(screen.getByLabelText("Expression")).toHaveValue("8 + 3");
  expect(JSON.parse(fetch.mock.calls[1][1].body)).toEqual({
    expression: "8 + 3",
  });
  await click("7");
  expect(screen.getByLabelText("Expression")).toHaveValue("7");
});
it.each([
  "0.333333333333",
  "-8",
  "99999999999999980000000000000001",
  "0.000000000000000000000001",
])("reuses displayed result %s without numeric conversion", async (value) => {
  const fetch = mockResult(value);
  render(<App />);
  fill("1+2");
  await click("Calculate");
  await screen.findByText(value, { selector: '[role="status"]' });
  await click("Multiply");
  await click("3");
  await click("Calculate");
  expect(JSON.parse(fetch.mock.calls[1][1].body)).toEqual({
    expression: value + " × 3",
  });
});
it("supports keyboard continuation and equals", async () => {
  mockResult("5");
  render(<App />);
  fill("2+3");
  const input = screen.getByLabelText("Expression");
  await userEvent.click(input);
  await userEvent.keyboard("=");
  await screen.findByText("5", { selector: '[role="status"]' });
  await userEvent.keyboard("*2");
  expect(input).toHaveValue("5 × 2");
  await userEvent.keyboard("{Escape}");
  expect(input).toHaveValue("");
});
it("allows backspace and clear", async () => {
  render(<App />);
  fill("12+3");
  await click("Backspace");
  expect(screen.getByLabelText("Expression")).toHaveValue("12+");
  await click("Clear");
  expect(screen.getByLabelText("Expression")).toHaveValue("");
});
it("validates empty input and unsupported characters before a request", async () => {
  const fetch = vi.fn();
  vi.stubGlobal("fetch", fetch);
  render(<App />);
  await click("Calculate");
  expect(screen.getByRole("alert")).toHaveTextContent("Enter an expression");
  fill("1e3");
  await click("Calculate");
  expect(screen.getByRole("alert")).toHaveTextContent("Use numbers");
  expect(fetch).not.toHaveBeenCalled();
});
it("shows backend syntax errors and permits correction", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue(
      new Response('{"error":{"message":"Enter a complete expression."}}', {
        status: 400,
      }),
    ),
  );
  render(<App />);
  fill("1+");
  await click("Calculate");
  expect(await screen.findByRole("alert")).toHaveTextContent(
    "Enter a complete expression",
  );
  await click("2");
  expect(screen.queryByRole("alert")).not.toBeInTheDocument();
});
it("locks controls during requests and restores them", async () => {
  let resolve!: (response: Response) => void;
  vi.stubGlobal(
    "fetch",
    vi.fn(
      () =>
        new Promise<Response>((done) => {
          resolve = done;
        }),
    ),
  );
  render(<App />);
  fill("2+3");
  await click("Calculate");
  expect(screen.getByLabelText("Expression")).toBeDisabled();
  expect(screen.getByRole("button", { name: "Clear" })).toBeDisabled();
  expect(screen.getByRole("button", { name: "Calculate" })).toBeDisabled();
  expect(screen.getByLabelText("Expression")).toHaveValue("2+3");
  expect(screen.getByLabelText("Expression").parentElement).not.toHaveClass(
    "has-result",
  );
  expect(screen.getByRole("status")).toBeEmptyDOMElement();
  expect(screen.queryByText("Calculating…")).not.toBeInTheDocument();
  await click("Calculate");
  await userEvent.keyboard("{Enter}=");
  expect(fetch).toHaveBeenCalledTimes(1);
  resolve(new Response('{"result":"5"}'));
  await screen.findByText("5", { selector: '[role="status"]' });
  expect(screen.getByLabelText("Expression")).toBeEnabled();
});
it("keeps the displayed result while recalculating", async () => {
  let resolve!: (response: Response) => void;
  vi.stubGlobal(
    "fetch",
    vi
      .fn()
      .mockResolvedValueOnce(new Response('{"result":"5"}'))
      .mockImplementationOnce(
        () =>
          new Promise<Response>((done) => {
            resolve = done;
          }),
      ),
  );
  render(<App />);
  fill("2+3");
  await click("Calculate");
  await screen.findByText("5", { selector: '[role="status"]' });
  await click("Calculate");
  expect(screen.getByRole("status")).toHaveTextContent("5");
  expect(screen.getByLabelText("Expression")).toHaveValue("2+3");
  expect(screen.getByLabelText("Expression").parentElement).toHaveClass(
    "has-result",
  );
  expect(screen.queryByText("Calculating…")).not.toBeInTheDocument();
  resolve(new Response('{"result":"5"}'));
  await waitFor(() =>
    expect(screen.getByRole("button", { name: "Calculate" })).toBeEnabled(),
  );
  expect(screen.getByRole("status")).toHaveTextContent("5");
});
it("handles network errors", async () => {
  vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new TypeError("offline")));
  render(<App />);
  fill("1+2");
  await click("Calculate");
  expect(await screen.findByRole("alert")).toHaveTextContent("Cannot reach");
  expect(screen.getByLabelText("Expression")).toBeEnabled();
});
it("prevents keypad entry beyond the expression limit", async () => {
  render(<App />);
  fill("1".repeat(256));
  await click("2");
  expect(screen.getByLabelText("Expression")).toHaveValue("1".repeat(256));
  expect(screen.getByRole("alert")).toHaveTextContent("256 characters");
});

it("shows the active expression prominently, then moves it above the result", async () => {
  mockResult("14");
  render(<App />);
  fill("8+9-3");
  const input = screen.getByLabelText("Expression");
  expect(input.parentElement).not.toHaveClass("has-result");
  expect(screen.getByRole("status")).toBeEmptyDOMElement();
  await click("Calculate");
  await screen.findByText("14");
  expect(input.parentElement).toHaveClass("has-result");
  expect(input).toHaveValue("8+9-3");
  await click("Subtract");
  expect(input).toHaveValue("14 − ");
  expect(input.parentElement).not.toHaveClass("has-result");
  expect(screen.getByRole("status")).toBeEmptyDOMElement();
});

it("accepts keyboard input immediately without clicking the display", async () => {
  const fetch = mockResult("14");
  render(<App />);
  await userEvent.keyboard("2+3*4{Enter}");
  await screen.findByText("14");
  expect(JSON.parse(fetch.mock.calls[0][1].body)).toEqual({
    expression: "2+3×4",
  });
  expect(fetch).toHaveBeenCalledTimes(1);
});

it("keeps keyboard input working after clicking keypad buttons and calculating", async () => {
  const fetch = vi
    .fn()
    .mockResolvedValueOnce(new Response('{"result":"8"}'))
    .mockResolvedValueOnce(new Response('{"result":"11"}'));
  vi.stubGlobal("fetch", fetch);
  render(<App />);
  await click("3");
  await userEvent.keyboard("+5{Enter}");
  await screen.findByText("8", { selector: '[role="status"]' });
  await userEvent.keyboard("+3=");
  await screen.findByText("11");
  expect(screen.getByLabelText("Expression")).toHaveValue("8 + 3");
  expect(fetch).toHaveBeenCalledTimes(2);
});

it("supports keyboard parentheses, decimals, backspace and escape outside the display", async () => {
  render(<App />);
  await userEvent.keyboard("(1.5+2)/3{Backspace}4");
  expect(screen.getByLabelText("Expression")).toHaveValue("(1.5+2)÷4");
  await userEvent.keyboard("{Escape}");
  expect(screen.getByLabelText("Expression")).toHaveValue("");
});

it("does not intercept modifier shortcuts or typing in another editable field", async () => {
  render(
    <>
      <App />
      <input aria-label="Other field" />
    </>,
  );
  await userEvent.keyboard("{Control>}1{/Control}{Meta>}2{/Meta}{Alt>}3{/Alt}");
  expect(screen.getByLabelText("Expression")).toHaveValue("");
  await userEvent.click(screen.getByLabelText("Other field"));
  await userEvent.keyboard("123");
  expect(screen.getByLabelText("Other field")).toHaveValue("123");
  expect(screen.getByLabelText("Expression")).toHaveValue("");
});

it("preserves text selection editing without duplicating keystrokes", async () => {
  render(<App />);
  const input = screen.getByLabelText("Expression");
  await userEvent.click(input);
  await userEvent.keyboard("12+3");
  await userEvent.keyboard("{Control>}a{/Control}4+5");
  expect(input).toHaveValue("4+5");
});

it("ignores global shortcuts while a request is pending", async () => {
  let resolve!: (response: Response) => void;
  vi.stubGlobal(
    "fetch",
    vi.fn(
      () =>
        new Promise<Response>((done) => {
          resolve = done;
        }),
    ),
  );
  render(<App />);
  await userEvent.keyboard("2+3{Enter}");
  await userEvent.keyboard("9{Escape}{Backspace}{Enter}");
  expect(screen.getByLabelText("Expression")).toHaveValue("2+3");
  resolve(new Response('{"result":"5"}'));
  await screen.findByText("5", { selector: '[role="status"]' });
});

it("builds optional operations using the keypad", async () => {
  const fetch = mockResult("8.5");
  render(<App />);
  for (const key of [
    "Square root",
    "2",
    "5",
    "Percentage",
    "Close parenthesis",
    "Add",
    "2",
    "Power",
    "3",
  ])
    await click(key);
  expect(screen.getByLabelText("Expression")).toHaveValue("sqrt(25%)+2^3");
  await click("Calculate");
  await screen.findByText("8.5");
  expect(JSON.parse(fetch.mock.calls[0][1].body)).toEqual({
    expression: "sqrt(25%)+2^3",
  });
});

it.each([
  ["9", "Square root", "sqrt(9)"],
  ["50", "Percentage", "(50)%"],
  ["-2", "Power", "(-2) ^ "],
])(
  "applies %s via %s to the displayed result",
  async (value, key, expression) => {
    mockResult(value);
    render(<App />);
    fill("1+2");
    await click("Calculate");
    await screen.findByText(value, { selector: '[role="status"]' });
    await click(key);
    expect(screen.getByLabelText("Expression")).toHaveValue(expression);
  },
);

it("accepts sqrt, powers and percentages from the keyboard without display focus", async () => {
  const fetch = mockResult("11.1");
  render(<App />);
  await userEvent.keyboard("sqrt(9)+2^3+10%{Enter}");
  await screen.findByText("11.1");
  expect(JSON.parse(fetch.mock.calls[0][1].body)).toEqual({
    expression: "sqrt(9)+2^3+10%",
  });
});

it("supports percentage continuation in the expression field", async () => {
  mockResult("50");
  render(<App />);
  fill("25+25");
  await userEvent.click(screen.getByLabelText("Expression"));
  await userEvent.keyboard("=");
  await screen.findByText("50");
  await userEvent.keyboard("%");
  expect(screen.getByLabelText("Expression")).toHaveValue("(50)%");
});

it("accepts the square-root character while editing", async () => {
  render(<App />);
  await userEvent.click(screen.getByLabelText("Expression"));
  await userEvent.keyboard("√9)");
  expect(screen.getByLabelText("Expression")).toHaveValue("sqrt(9)");
});

it("shows mathematical errors for optional operations and recovers", async () => {
  vi.stubGlobal(
    "fetch",
    vi
      .fn()
      .mockResolvedValue(
        new Response(
          '{"error":{"message":"Cannot take the square root of a negative number."}}',
          { status: 400 },
        ),
      ),
  );
  render(<App />);
  fill("sqrt(-1)");
  await click("Calculate");
  expect(await screen.findByRole("alert")).toHaveTextContent("negative number");
  await click("Clear");
  expect(screen.queryByRole("alert")).not.toBeInTheDocument();
});
