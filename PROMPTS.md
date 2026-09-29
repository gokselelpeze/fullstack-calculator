# AI assistance and prompts

Tool: OpenAI Codex. Used for MVP analysis, numeric/API design, implementation, test development, and verification. No subagents were used. The user requested phase one only; optional operations and Docker are deferred.

## Initial planning prompt

> I have a technical assesment where I need to create a full-stack calculator application. Im pasting the task description down. we will start with the must components and integrate the optional ones later. lets conduct an analysis for an mvp implementation using Go for the backend and React with TypeScript for the frontend. Since there are mathematical operations involved, we need to pay close attention to calculation errors and edge cases, such as division by zero. try to plan as simply as possible but in the right way

The supplied assessment required React (TypeScript preferred), a Go REST backend, addition/subtraction/multiplication/division, input validation and JSON errors, a responsive UI, unit tests and coverage reports, setup/API/design documentation, a Git repository, and disclosure of AI prompts. Optional features: exponentiation, square root, percentage, and Docker. Suggested time budget: 2–4 hours.

## Follow-up planning prompt

> we will also implement the optional things after the initial implementation. will we need to make any fundamental changes while doing this? is there anything we should consider right now? or will we simply be adding new features?

Accepted plan: decimal strings over HTTP, decimal arithmetic in Go, bounded operands, exact basic arithmetic except division rounded to 12 fractional places, a single calculation endpoint, and operand presence represented explicitly for later unary operations.

## Implementation prompt

> okay implement the first phase

## Workspace instruction

> Use AGENTS.md strictly. First inspect the relevant files, then make the smallest safe change. Do not refactor unrelated code. Do not add new abstractions unless necessary.

The workspace was empty and no additional AGENTS.md files were found during inspection. Codex implemented the application and tests, then ran build, coverage, static checks, and integration verification. See docs/coverage.md for actual results and limitations.

## Product copy revision

> firstly remove all the texts like **SMALL TOOL. CLEAR ANSWERS.  just focus on the product**

Removed the promotional copy, intro, and footer.

## Expression calculator revision

> nice, first design and implementation is ok but i have few questions. right now we can not use more than 2 numbers and can not chain operations. this is so basic we select and send the request. can we implement a calculator like supports chaining operations? maybe we can use stack to ensure the operation is right. and the design can be similar to calculator not like this

Replaced the two-input form with a keypad and editable expression display. Added a stack-based Go evaluator with precedence, parentheses, unary signs, exact rational intermediates, and rounding of final repeating decimals. Preserved the original two-operand API. Added parser, API, UI, and fuzz tests, then verified the running calculator in the browser.

## Continuation display refinement

> you can show the last operation on the top but when I try to do a new calculation you can reset the top part with the last operation
>
> for example first calculation is 3+5. after i press = the top part shows 3+5 and result is below 8. then i add 3 more. top part is 8 + 3 below 11

> this is the behavior I want

The second prompt included a screen recording of a calculator: active expressions appear as the main display; after equals, the expression appears above the result; subsequent operators start from the displayed result. Implemented this display transition and decimal-string result reuse. Expression literals now use the existing overall length bound so computed results can be reused without the legacy per-operand limits.

## Keyboard focus fix

> I can not type with keyboard in the calculator can you fix that

Added keyboard shortcuts outside the expression field so typing works immediately and after keypad clicks. Preserved textarea selection editing, modifier shortcuts, and request loading behavior. Added regression tests and verified keyboard-only and mixed mouse/keyboard calculations in the running browser.

## Optional features

> looking good, lets implement the optional features. can you walk me through how we gonna add these

The proposed rules were integer powers with conventional precedence and right associativity; nonnegative square roots with exact or 12-place rounded results; postfix percentage as divide-by-100; and a multi-stage Docker build serving the frontend through Go.

> lets add it

Added the optional expression operations, calculation-size limits, keypad/keyboard support, tests, Go static-file serving, a non-root Docker image, and deployment documentation.

## Readability refactor

> functionality is good working really well but when I try to review the code you wrote too long functions and magic strings, can you review your code and refactor as clean as possible consider solid principles and best practices

Separated token scanning, grammar, arithmetic, HTTP validation/response handling, React presentation, state, keyboard events, and pure input rules. Added named tokens, commands, limits, and error codes. Verified with the existing tests, fuzz testing, and 336 API response comparisons against the preceding implementation. Expanded frontend coverage collection to include all extracted application modules.

## Sezzle-inspired visual design

> Im doing this task for https://sezzle.com/, search their color palette and redesign the frontend page

Redesigned the calculator using Sezzle's published deep purple, lilac, mint, peach, and coral palette. Kept the interface focused on the calculator, without using or recreating Sezzle's protected logo assets. Verified the result at desktop and mobile viewport sizes and reran the frontend tests and production build.
