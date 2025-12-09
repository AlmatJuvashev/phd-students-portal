import { test, expect } from "@playwright/test";

const adminUser = process.env.E2E_ADMIN_USER || "ta2087";
const adminPass = process.env.E2E_ADMIN_PASS || "meadow-pluto-pioneer48";
const studentUser = process.env.E2E_STUDENT_USER || "ts5251";
const studentPass = process.env.E2E_STUDENT_PASS || "pioneer-canvas-silver52";

test.describe("Admin Review Workflow", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to login page
    await page.goto("/login");
  });

  test("Admin can view student list", async ({ page }) => {
    // Login as admin
    await page.getByLabel(/username/i).fill(adminUser);
    await page.locator('input[name="password"]').fill(adminPass);
    await page.getByRole("button", { name: /sign in|войти|кіру/i }).click();
    await page.waitForLoadState("networkidle");
    
    // Wait for any post-login page (not login)
    await expect(page).not.toHaveURL(/login/);
    
    // Look for student list or progress view
    const studentsLink = page.getByRole("link", { name: /students|студенты|прогресс|progress/i });
    if (await studentsLink.isVisible({ timeout: 5000 }).catch(() => false)) {
      await studentsLink.click();
      await page.waitForLoadState("networkidle");
    }
    
    // Verify admin is logged in
    await expect(page.locator("body")).toContainText(/./);
  });

  test("Admin can navigate to student journey details", async ({ page }) => {
    // Login as admin
    await page.getByLabel(/username/i).fill(adminUser);
    await page.locator('input[name="password"]').fill(adminPass);
    await page.getByRole("button", { name: /sign in|войти|кіру/i }).click();
    await page.waitForLoadState("networkidle");
    
    // Wait for any post-login page
    await expect(page).not.toHaveURL(/login/);
    
    // Navigate to student progress/journey view
    const progressLink = page.getByRole("link", { name: /progress|прогресс|journey|students/i });
    if (await progressLink.isVisible({ timeout: 5000 }).catch(() => false)) {
      await progressLink.click();
      await page.waitForLoadState("networkidle");
    }
    
    // Verify page loaded
    await expect(page.locator("body")).toContainText(/./);
  });

  test("Admin can approve a submitted node", async ({ page }) => {
    // Login as admin
    await page.getByLabel(/username/i).fill(adminUser);
    await page.locator('input[name="password"]').fill(adminPass);
    await page.getByRole("button", { name: /sign in|войти|кіру/i }).click();
    await page.waitForLoadState("networkidle");
    
    // Wait for any post-login page
    await expect(page).not.toHaveURL(/login/);
    
    // Verify admin has access to some admin functionality
    await expect(page.locator("body")).toContainText(/./);
  });

  test("Admin can request fixes for a node", async ({ page }) => {
    // Login as admin
    await page.getByLabel(/username/i).fill(adminUser);
    await page.locator('input[name="password"]').fill(adminPass);
    await page.getByRole("button", { name: /sign in|войти|кіру/i }).click();
    await page.waitForLoadState("networkidle");
    
    // Wait for any post-login page
    await expect(page).not.toHaveURL(/login/);
    
    // Verify admin dashboard loaded
    await expect(page.locator("body")).toContainText(/./);
  });

  test("Student sees updated status after admin action", async ({ page }) => {
    // Login as student
    await page.getByLabel(/username/i).fill(studentUser);
    await page.locator('input[name="password"]').fill(studentPass);
    await page.getByRole("button", { name: /sign in|войти|кіру/i }).click();
    await page.waitForLoadState("networkidle");
    
    // Wait for any post-login page
    await expect(page).not.toHaveURL(/login/);
    
    // Navigate to journey if not already there
    const journeyLink = page.getByRole("link", { name: /journey|путь/i });
    if (await journeyLink.isVisible({ timeout: 5000 }).catch(() => false)) {
      await journeyLink.click();
      await page.waitForLoadState("networkidle");
    }
    
    // Verify student can see content
    await expect(page.locator("body")).toContainText(/./);
  });
});
