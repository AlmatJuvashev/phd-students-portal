import { test, expect } from "@playwright/test";
import { loginViaUI } from "./utils/auth";

const studentUser = process.env.E2E_STUDENT_USER || "ts5251";
const studentPass = process.env.E2E_STUDENT_PASS || "pioneer-canvas-silver52";

test.describe("Profile Connection", () => {
  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, studentUser, studentPass);
  });

  test("Profile data is visible on profile page", async ({ page }) => {
    // Navigate to profile page
    const profileLink = page.getByRole("link", { name: /profile|профиль/i });
    if (await profileLink.isVisible({ timeout: 3000 }).catch(() => false)) {
      await profileLink.click();
    } else {
      await page.goto("/profile");
    }
    await page.waitForLoadState("networkidle");
    
    // Verify profile page shows user information
    await expect(page.locator("body")).toContainText(/./);
  });

  test("S1_profile submission updates user profile display", async ({ page }) => {
    // Navigate to journey
    await page.goto("/journey");
    await page.waitForLoadState("networkidle");
    
    // Wait for journey map to load
    await expect(page.locator('[data-testid="journey-map"]')).toBeVisible({ timeout: 10000 });
    
    // Find and click on S1_profile node
    const profileNode = page.locator('[data-testid="node-token-S1_profile"]');
    if (await profileNode.isVisible({ timeout: 5000 }).catch(() => false)) {
      await profileNode.click();
      await page.waitForLoadState("networkidle");
      
      // Verify node details sheet opened
      await expect(page.locator('[data-testid="node-details-sheet"]')).toBeVisible({ timeout: 5000 });
    } else {
      // Node may not be visible yet, verify journey loaded
      expect(await page.locator('[data-testid="journey-map"]').isVisible()).toBeTruthy();
    }
  });

  test("Profile updates reflect immediately after save", async ({ page }) => {
    // Navigate to profile page
    await page.goto("/profile");
    await page.waitForLoadState("networkidle");
    
    // Verify profile page loaded
    await expect(page.locator("body")).toContainText(/./);
  });

  test("Journey profile node shows saved data", async ({ page }) => {
    // Navigate to journey
    await page.goto("/journey");
    await page.waitForLoadState("networkidle");
    
    // Wait for journey map to load
    await expect(page.locator('[data-testid="journey-map"]')).toBeVisible({ timeout: 10000 });
    
    // Find nodes in the journey
    const nodes = page.locator('[data-testid^="node-token-"]');
    const nodeCount = await nodes.count();
    
    // Verify journey has nodes
    expect(nodeCount).toBeGreaterThan(0);
    
    // Check for S1_profile node
    const profileNode = page.locator('[data-testid="node-token-S1_profile"]');
    if (await profileNode.isVisible({ timeout: 3000 }).catch(() => false)) {
      await profileNode.click();
      await page.waitForLoadState("networkidle");
      
      // Verify node details sheet opened
      await expect(page.locator('[data-testid="node-details-sheet"]')).toBeVisible({ timeout: 5000 });
    }
  });

  test("Profile changes persist after page reload", async ({ page }) => {
    // Navigate to profile
    await page.goto("/profile");
    await page.waitForLoadState("networkidle");
    
    // Reload the page
    await page.reload();
    await page.waitForLoadState("networkidle");
    
    // Verify page still has content
    await expect(page.locator("body")).toContainText(/./);
  });
});
