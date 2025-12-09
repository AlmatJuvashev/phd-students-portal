import { test as setup } from "@playwright/test";
import { getAuthToken } from "./utils/auth";

const authFile = "playwright/.auth/user.json";
const baseURL = process.env.E2E_BASE_URL || "http://localhost:5175";

setup("authenticate (optional)", async ({ page }) => {
  if (!process.env.E2E_STUDENT_USER || !process.env.E2E_STUDENT_PASS) return;
  const token = await getAuthToken("student", baseURL);
  await page.goto(baseURL);
  await page.evaluate((tk) => localStorage.setItem("token", tk), token);
  await page.context().storageState({ path: authFile });
});
