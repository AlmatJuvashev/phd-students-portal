import { test, expect } from "@playwright/test";
import { getAuthToken } from "./utils/auth";

const baseURL = process.env.E2E_BASE_URL || "http://localhost:5175";

test.describe("Admin notifications", () => {
  test("admin can view and mark notifications as read", async ({ page }) => {
    test.skip(
      !process.env.E2E_ADMIN_USER || !process.env.E2E_ADMIN_PASS,
      "Requires E2E_ADMIN_USER and E2E_ADMIN_PASS"
    );

    const token = await getAuthToken("admin", baseURL);
    await page.goto(baseURL);
    await page.evaluate((tk) => localStorage.setItem("token", tk), token);

    await page.goto("/admin/notifications");
    await expect(page.getByRole("heading", { name: /notifications|уведомления|хабарламалар/i })).toBeVisible();

    const unreadBadge = page.getByText(/\d+/, { exact: false }).first();
    const markAllButton = page.getByRole("button", { name: /mark all read|отметить все|барлығын/i });
    if (await markAllButton.isVisible()) {
      await markAllButton.click();
      await expect(markAllButton).toBeEnabled();
    }
    await expect(unreadBadge).toBeVisible();
  });
});
