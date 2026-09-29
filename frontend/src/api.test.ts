import { afterEach, expect, it, vi } from "vitest";
import { calculate } from "./api";

afterEach(() => {
  vi.unstubAllGlobals();
  vi.useRealTimers();
});
const input = { expression: "9007199254740992 + 1" };
it("sends and returns decimal strings without numeric conversion", async () => {
  const fetch = vi
    .fn()
    .mockResolvedValue(
      new Response(JSON.stringify({ result: "9007199254740993" })),
    );
  vi.stubGlobal("fetch", fetch);
  expect(await calculate(input)).toBe("9007199254740993");
  expect(fetch).toHaveBeenCalledWith(
    "/api/calculate",
    expect.objectContaining({ method: "POST", body: JSON.stringify(input) }),
  );
});
it("returns the API error message", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue(
      new Response(
        JSON.stringify({
          error: {
            code: "DIVISION_BY_ZERO",
            message: "Cannot divide by zero.",
          },
        }),
        { status: 400 },
      ),
    ),
  );
  await expect(calculate(input)).rejects.toThrow("Cannot divide by zero.");
});
it("handles an unstructured server error", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue(new Response("{}", { status: 500 })),
  );
  await expect(calculate(input)).rejects.toThrow("The calculation failed");
});
it.each(["not json", '{"result":123}', '{"result":"NaN"}', "null"])(
  "rejects malformed response %s",
  async (body) => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(new Response(body)));
    await expect(calculate(input)).rejects.toThrow(/invalid/);
  },
);
it("handles network failures", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockRejectedValue(new TypeError("Failed to fetch")),
  );
  await expect(calculate(input)).rejects.toThrow("Cannot reach the calculator");
});
it("aborts slow requests with a useful error", async () => {
  vi.useFakeTimers();
  vi.stubGlobal(
    "fetch",
    vi.fn(
      (_url, options: RequestInit) =>
        new Promise((_resolve, reject) => {
          options.signal!.addEventListener("abort", () =>
            reject(new DOMException("Aborted", "AbortError")),
          );
        }),
    ),
  );
  const expectation = expect(calculate(input)).rejects.toThrow("timed out");
  await vi.advanceTimersByTimeAsync(10_000);
  await expectation;
  expect(vi.getTimerCount()).toBe(0);
});

it("times out if response headers arrive but the body stalls", async () => {
  vi.useFakeTimers();
  vi.stubGlobal(
    "fetch",
    vi.fn((_url, options: RequestInit) =>
      Promise.resolve({
        ok: true,
        json: () =>
          new Promise((_resolve, reject) => {
            options.signal!.addEventListener("abort", () =>
              reject(new DOMException("Aborted", "AbortError")),
            );
          }),
      }),
    ),
  );
  const expectation = expect(calculate(input)).rejects.toThrow("timed out");
  await vi.advanceTimersByTimeAsync(10_000);
  await expectation;
  expect(vi.getTimerCount()).toBe(0);
});
