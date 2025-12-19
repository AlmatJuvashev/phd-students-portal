import { test as setup } from "@playwright/test";
import { getAuthToken } from "./utils/auth";

const authFile = "playwright/.auth/user.json";
const baseURL = process.env.E2E_BASE_URL || "http://localhost:5175";

setup("authenticate (optional)", async ({ page, context }) => {
  if (!process.env.E2E_STUDENT_USER || !process.env.E2E_STUDENT_PASS) return;
  const token = await getAuthToken("student", baseURL);
  
  // Set cookie in context
  await context.addCookies([{
    name: "jwt_token",
    value: token,
    domain: "localhost", // Adjust if needed
    path: "/",
    httpOnly: true,
    secure: false, // Localhost usually false
    sameSite: "Lax"
  }]);
  
  await page.context().storageState({ path: authFile });
});
