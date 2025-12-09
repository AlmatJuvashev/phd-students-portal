# Frontend Testing Guidelines

This document outlines the standards and best practices for writing End-to-End (E2E) tests in this project using Playwright.

## 1. Selector Strategy

We prioritize stability and resilience in our selectors.

### Priority Order:
1.  **`data-testid`**: The gold standard. Use this attribute explicitly for testing.
    *   *Example*: `page.getByTestId('user-menu-button')`
    *   *Implementation*: Add `data-testid="your-id"` to the React component or HTML element.
2.  **ARIA Roles**: Use standard accessibility roles when `data-testid` is not available.
    *   *Example*: `page.getByRole('button', { name: /save/i })`
3.  **Text Content**: Use text with localization-resilient Regex.
    *   *Example*: `page.getByText(/submit|отправить/i)`
4.  **CSS Selectors**: Use only as a last resort for structural wrappers.
    *   *Example*: `.fixed.z-50 > .relative` (Avoid if possible)

### Action Item:
When you encounter a fragile selector (like a complex CSS chain), **refactor the component** to add a `data-testid`.

## 2. Data Management

Tests must be independent and deterministic.

*   **Fresh Data**: **ALWAYS** create fresh test data for each test run. Do not rely on seed data (like `ts5251`) which may be modified by other tests or manual usage.
*   **API Helpers**: Use the `tests/utils/api_auth.ts` helpers to create users, entities, etc., before interacting with the UI.
    *   *Pattern*:
        ```typescript
        const token = await loginViaAPI(request, adminUser, adminPass);
        const user = await createUserViaAPI(token, { ... });
        await loginViaUI(page, user.username, user.password);
        ```

## 3. Resilience & Timing

*   **Explicit Waits**: Do not rely on implicit timeouts for complex state changes. Assert the end state.
    *   *Bad*: `await page.waitForTimeout(1000)`
    *   *Good*: `await expect(modal).not.toBeVisible()`
*   **Localization**: Always use Case-Insensitive Regex for text matching to handle multiple languages and English fallbacks.
    *   *Example*: `getByRole('heading', { name: /chat|чат/i })`

## 4. Layout Awareness

*   **Context Matters**: Be aware that different layouts (Admin vs. App) may use different tags (e.g., `div` vs `header`).
*   **Sticky Elements**: Use `{ force: true }` if an element might be obscured by a sticky header or animation, but prefer waiting for the overlay to disappear.
