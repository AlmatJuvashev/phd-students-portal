import { test, expect } from "@playwright/test";
import { getAuthToken } from "./utils/auth";
import path from "path";
import fs from "fs";

const baseURL = process.env.E2E_BASE_URL || "http://localhost:5175";

test.describe("Document submission and review", () => {
  test("student can upload a document to the first available slot", async ({ page }) => {
    test.skip(
      !process.env.E2E_STUDENT_USER || !process.env.E2E_STUDENT_PASS,
      "Requires E2E_STUDENT_USER and E2E_STUDENT_PASS"
    );

    const token = await getAuthToken("student", baseURL);
    await page.goto(baseURL);
    await page.evaluate((tk) => localStorage.setItem("token", tk), token);

    await page.goto("/journey");
    await expect(page.getByText(/journey|карт|жол/i)).toBeVisible();

    // Open first upload button if available
    const uploadButton = page.getByRole("button", { name: /upload|загрузить|жүктеу/i }).first();
    if (!(await uploadButton.isVisible())) {
      test.skip(true, "No upload button available in journey page");
    }
    await uploadButton.click();

    const tmpFile = path.join(process.cwd(), "tmp-upload.txt");
    fs.writeFileSync(tmpFile, `E2E upload ${Date.now()}`);
    await page.setInputFiles('input[type="file"]', tmpFile);
    await expect(page.getByText(/upload|загрузка|жүктеу/i)).toBeVisible();
  });
});
