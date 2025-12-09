import { test, expect } from "@playwright/test";
import { getAuthToken } from "./utils/auth";

const baseURL = process.env.E2E_BASE_URL || "http://localhost:5175";

test.describe("Chat", () => {
  test("student can open chat and send a message", async ({ page, request }) => {
    test.skip(
      !process.env.E2E_STUDENT_USER || !process.env.E2E_STUDENT_PASS,
      "Requires E2E_STUDENT_USER and E2E_STUDENT_PASS"
    );

    const token = await getAuthToken("student", baseURL);
    await page.goto(baseURL);
    await page.evaluate((tk) => localStorage.setItem("token", tk), token);

    await page.goto("/chat");
    await expect(page.getByRole("heading", { name: /messages|сообщения|хабарламалар/i })).toBeVisible();

    // Select the first available room
    const roomButton = page.locator("section", { hasText: /rooms|комнаты|бөлмелер/i }).locator("button").first();
    await roomButton.click();

    const message = `e2e message ${Date.now()}`;
    await page.getByRole("textbox").fill(message);
    await page.getByRole("button", { name: /send|отправить|жіберу/i }).click();

    await expect(page.getByText(message)).toBeVisible({ timeout: 10000 });
  });
});
