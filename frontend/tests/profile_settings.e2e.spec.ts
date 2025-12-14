import { test, expect } from '@playwright/test';

test.describe('Profile Settings', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/login');
    await page.fill('input[type="email"]', 'student@example.com');
    await page.fill('input[type="password"]', 'password');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/map');
  });

  test('should navigate to profile and update bio', async ({ page }) => {
    // Navigate to Profile page
    // Assuming there is a menu item or we click on user avatar
    // Or we can go directly to /profile
    await page.goto('/profile');

    // Verify Header
    await expect(page.locator('h1')).toContainText(/Profile|Профиль/i);

    // Fill Bio
    const bioInput = page.locator('textarea[name="bio"]');
    await bioInput.fill('This is a test bio updated by Playwright');

    // Fill Password for confirmation
    await page.fill('input[name="current_password"]', 'password');

    // Submit
    await page.click('button[type="submit"]');

    // Verify Success Toast
    // Using toast locator pattern
    await expect(page.locator('text=Profile Updated')).toBeVisible({ timeout: 10000 });
  });

  test('should fail to update without current password', async ({ page }) => {
    await page.goto('/profile');

    await page.fill('textarea[name="bio"]', 'Should fail bio');
    
    // Clear password if prefilled or ensure it's empty
    await page.fill('input[name="current_password"]', '');

    await page.click('button[type="submit"]');

    // Verify Validation Error
    // Should see "Password required" or similar under the field
    await expect(page.locator('text=Password is required')).toBeVisible();
  });
});
