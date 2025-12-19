import { test, expect } from '@playwright/test';

test.describe('Profile Settings', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/login');
    await page.fill('input[name="username"]', 'demo.student1');
    await page.fill('input[type="password"]', 'demopassword123!');
    await page.click('button[type="submit"]');
    // Wait for login to succeed (dashboard)
    await page.waitForURL((url) => url.pathname === '/');
  });

  test('should navigate to profile and update bio', async ({ page }) => {
    // Navigate to Profile page
    // Assuming there is a menu item or we click on user avatar
    // Or we can go directly to /profile
    await page.goto('/profile');

    // Verify we are on profile page
    await expect(page.getByRole('heading', { name: 'Personal Information' })).toBeVisible();

    // Click Edit Profile to show form
    await page.click('button:has-text("Edit Profile")');

    // Fill Bio
    const bioInput = page.locator('textarea[name="bio"]');
    await bioInput.fill('This is a test bio updated by Playwright');

    // Fill Password for confirmation
    await page.fill('input[name="current_password"]', 'demopassword123!');

    // Submit
    await page.click('button[type="submit"]');

    // Verify Success Toast
    // Using toast locator pattern
    // Verify Success Toast (loose check first)
    // await expect(page.getByText('Profile updated successfully')).toBeVisible({ timeout: 10000 });
    const toast = page.locator('li[role="status"]'); // or .toast
    await expect(toast).toBeVisible({ timeout: 10000 });
    console.log('Toast text:', await toast.innerText());
  });

  test('should fail to update without current password', async ({ page }) => {
    await page.goto('/profile');

    // Click Edit Profile to show form
    await page.click('button:has-text("Edit Profile")');

    await page.fill('textarea[name="bio"]', 'Should fail bio');
    
    // Clear password if prefilled or ensure it's empty
    await page.fill('input[name="current_password"]', '');

    await page.click('button[type="submit"]');

    // Verify Validation Error
    // Should see "Password required" or similar under the field
    await expect(page.locator('text=Password is required')).toBeVisible();
  });
});
