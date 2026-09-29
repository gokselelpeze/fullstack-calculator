export type Calculation = { expression: string };

const CALCULATE_ENDPOINT = "/api/calculate";
const REQUEST_TIMEOUT_MS = 10_000;
const resultPattern = /^-?\d+(\.\d+)?$/;
export const CALCULATION_FAILED_MESSAGE =
  "The calculation failed. Please try again.";
const messages = {
  timeout: "The calculation timed out. Please try again.",
  network: "Cannot reach the calculator. Check your connection and try again.",
  invalidResponse: "The server returned an invalid response. Please try again.",
  invalidResult: "The server returned an invalid result. Please try again.",
} as const;

export async function calculate(input: Calculation): Promise<string> {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);
  try {
    const response = await sendCalculation(input, controller.signal);
    const payload = await readPayload(response, controller.signal);
    return parseResult(payload, response.ok);
  } finally {
    clearTimeout(timeout);
  }
}

async function sendCalculation(
  input: Calculation,
  signal: AbortSignal,
): Promise<Response> {
  try {
    return await fetch(CALCULATE_ENDPOINT, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(input),
      signal,
    });
  } catch {
    throw new Error(signal.aborted ? messages.timeout : messages.network);
  }
}

async function readPayload(
  response: Response,
  signal: AbortSignal,
): Promise<unknown> {
  try {
    return await response.json();
  } catch {
    throw new Error(
      signal.aborted ? messages.timeout : messages.invalidResponse,
    );
  }
}

function parseResult(payload: unknown, ok: boolean): string {
  const data = asRecord(payload);
  if (!ok) {
    const message = asRecord(data?.error)?.message;
    throw new Error(
      typeof message === "string" ? message : CALCULATION_FAILED_MESSAGE,
    );
  }
  if (typeof data?.result !== "string" || !resultPattern.test(data.result))
    throw new Error(messages.invalidResult);
  return data.result;
}

function asRecord(value: unknown): Record<string, unknown> | undefined {
  return typeof value === "object" && value !== null
    ? (value as Record<string, unknown>)
    : undefined;
}
