import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';
import { loginViaAPI, createUserViaAPI } from './utils/api_auth';

const adminUser = process.env.E2E_ADMIN_USER || 'ta2087';
const adminPass = process.env.E2E_ADMIN_PASS || 'meadow-pluto-pioneer48';
const studentUser = process.env.E2E_STUDENT_USER || 'ts5251';
const studentPass = process.env.E2E_STUDENT_PASS || 'pioneer-canvas-silver52';

test.describe('Authentication & RBAC', () => {
  
  test('Login Page Loads', async ({ page }) => {
    await page.goto('/login');
    await expect(page).toHaveTitle(/PhD Student Portal/i);
    await expect(page.getByLabel(/username/i)).toBeVisible();
    await expect(page.locator('input[name="password"]')).toBeVisible();
  });

  test('Admin Login & Dashboard', async ({ page }) => {
    // Attempt login
    await loginViaUI(page, adminUser, adminPass);

    // Verify redirect to dashboard or home
    try {
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
    } catch (e) {
      // Log any error message on the page
      const error = await page.getByRole('alert').textContent().catch(() => 'No alert found');
      console.log(`Login failed. Error on page: ${error}`);
      throw e;
    }
    
    // Verify Admin-specific elements
    // Adjust selector based on actual UI (e.g., sidebar links)
    // Verify Admin-specific elements
    // Click Admin link in nav first
    await page.getByRole('link', { name: /admin|админ/i }).click();
    
    // Then verify Users link/card is visible on the admin dashboard
    await expect(page.getByRole('link', { name: /users|пользователи|пайдаланушылар/i })).toBeVisible();
    
    // Save state for next tests if needed, or just return token
  });

  test('Student Login & Journey', async ({ page }) => {
    // Login as Student
    await loginViaUI(page, studentUser, studentPass);

    // Verify redirect
    try {
        await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
    } catch (e) {
        const error = await page.getByRole('alert').textContent().catch(() => 'No alert found');
        console.log(`Student Login failed. Error: ${error}`);
        throw e;
    }

    // Verify Student-specific elements
    await expect(page.getByRole('link', { name: /journey|путь|жол/i })).toBeVisible();
    
    // Verify Admin elements are NOT visible
    await expect(page.getByRole('link', { name: /users|пользователи|пайдаланушылар/i })).not.toBeVisible();
  });


  test.skip('Logout', async ({ page }) => {
    await loginViaUI(page, adminUser, adminPass);
    
    // Find logout button (usually in a dropdown or sidebar)
    // The trigger is an Avatar button in the header
    // Debug: Log buttons in header
    const buttons = page.locator('header button');
    const count = await buttons.count();
    console.log(`Found ${count} buttons in header`);
    
    // The UserMenu is usually the last visible button on desktop
    const profileMenu = buttons.locator('visible=true').last();
    await expect(profileMenu).toBeVisible();
    await profileMenu.click();
    
    await page.getByText(/logout|выйти|шығу/i).click();
    
    await expect(page).toHaveURL(/\/login/);
  });

  test('Protected Route Redirect', async ({ page }) => {
    await page.goto('/dashboard'); // or any protected route
    await expect(page).toHaveURL(/\/login/);
  });
});
