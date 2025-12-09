import { test, expect } from '@playwright/test';
import path from 'path';

test.describe('Document Uploads', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    
    // Wait for login form to be ready
    await page.waitForSelector('input#username', { state: 'visible' });
    await page.waitForSelector('input#password', { state: 'visible' });
    
    const username = process.env.TEST_USERNAME || 'tu6260';
    const password = process.env.TEST_PASSWORD || 'thunder-pluto-river72';
    
    await page.locator('#username').fill(username);
    await page.locator('#password').fill(password);
    
    // Click submit button - exclude the "Show password" button which has type="button"
    await page.locator('form button:not([type="button"])').click();
    
    // Wait for navigation to complete
    await page.waitForLoadState('networkidle');
    
    // Give extra time for any client-side routing
    await page.waitForTimeout(1000);
  });

  test('should navigate to Doctoral profile and see form', async ({ page }) => {
    // Navigate to S1_profile node (which should be unlocked by default)
    await page.getByText('Doctoral profile').click();
    
    // Wait for node details to open
    await page.waitForTimeout(500);
    
    // Verify we can see the form by checking for the full_name input field
    await expect(page.locator('#full_name')).toBeVisible({ timeout: 10000 });
  });

  test('should see Publications List node', async ({ page }) => {
    // Try to find Publications node (may be locked)
    const pubNode = page.getByText('Publications List').or(page.getByText('Список публикаций'));
    await expect(pubNode).toBeVisible({ timeout: 10000 });
  });
});
