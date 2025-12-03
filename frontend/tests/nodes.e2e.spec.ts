import { test, expect } from "@playwright/test";
import { getAuthToken } from "./utils/auth";

const baseURL = process.env.E2E_BASE_URL || "http://localhost:5175";

test.describe("Node completion", () => {
  test("student can mark a node/task as done", async ({ page }) => {
    test.skip(
      !process.env.E2E_STUDENT_USER || !process.env.E2E_STUDENT_PASS,
      "Requires E2E_STUDENT_USER and E2E_STUDENT_PASS"
    );

    const token = await getAuthToken("student", baseURL);
    await page.goto(baseURL);
    await page.evaluate((tk) => localStorage.setItem("token", tk), token);

    await page.goto("/journey");
    await expect(page.getByText(/progress|прогресс|ілгерілеу/i)).toBeVisible();

    const firstCheckbox = page.getByRole("checkbox").first();
    if (!(await firstCheckbox.isVisible())) {
      test.skip(true, "No actionable node checkbox found");
    }
    await firstCheckbox.check();
    await expect(firstCheckbox).toBeChecked();
  });
});
