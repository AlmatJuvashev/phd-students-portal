import { request as playwrightRequest, APIRequestContext, expect, Page } from "@playwright/test";

type Role = "admin" | "superadmin" | "advisor" | "student";

function envFor(role: Role, key: "user" | "pass") {
  const suffix = role.toUpperCase();
  return (
    process.env[`E2E_${suffix}_USER`] && key === "user"
      ? process.env[`E2E_${suffix}_USER`]
      : process.env[`E2E_${suffix}_PASS`] && key === "pass"
        ? process.env[`E2E_${suffix}_PASS`]
        : undefined
  ) || process.env[`E2E_${suffix}_${key === "user" ? "USERNAME" : "PASSWORD"}`];
}

// Returns the cookie value or state
export async function getAuthToken(role: Role, baseURL: string): Promise<string> {
  const username = envFor(role, "user");
  const password = envFor(role, "pass");
  if (!username || !password) {
    throw new Error(`Missing credentials for role ${role}. Set E2E_${role.toUpperCase()}_USER and E2E_${role.toUpperCase()}_PASS`);
  }
  const req: APIRequestContext = await playwrightRequest.newContext({ baseURL });
  const res = await req.post("/api/auth/login", {
    data: { username, password },
  });
  expect(res.ok()).toBeTruthy();
  
  // Extract set-cookie header
  const setCookie = res.headers()["set-cookie"];
  if (!setCookie) {
    throw new Error(`Login failed for ${role}: no set-cookie header`);
  }
  
  // Simple parse for jwt_token
  const match = setCookie.match(/jwt_token=([^;]+)/);
  if (!match) {
    throw new Error(`Login failed for ${role}: jwt_token cookie not found`);
  }
  
  return match[1];
}

export async function loginViaUI(page: Page, username: string, password: string) {
  await page.goto("/login");
  await page.getByLabel(/username/i).fill(username);
  await page.locator('input[name="password"]').fill(password);
  await page.getByRole("button", { name: /sign in|войти|кіру/i }).click();
  await page.waitForLoadState("networkidle");
  // Wait for redirect to dashboard or home
  await page.waitForURL(/dashboard|journey|profile/);
}
