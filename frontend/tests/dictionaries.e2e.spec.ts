import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';

const adminUser = process.env.E2E_ADMIN_USER || 'ta2087';
const adminPass = process.env.E2E_ADMIN_PASS || 'meadow-pluto-pioneer48';

test.describe('Dictionary Management', () => {
  
  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, adminUser, adminPass);
    // Navigate to Admin Dashboard -> Dictionaries
    await page.getByRole('link', { name: /admin|админ/i }).click();
    await page.getByRole('link', { name: /dictionaries|справочники/i }).click();
  });

  test('Create New Program', async ({ page }) => {
    // Verify we are on Dictionaries page
    await expect(page.getByRole('heading', { name: /dictionaries|справочники/i })).toBeVisible();
    
    // Ensure "Programs" tab is active (default)
    await expect(page.getByRole('tab', { name: /programs|программы/i })).toHaveAttribute('data-state', 'active');

    // Click "Add Program"
    await page.getByRole('button', { name: /add program|добавить программу/i }).click();
    
    // Fill form (Modal)
    // Scope to modal (Radix UI uses role="dialog")
    const modal = page.getByRole('dialog');
    await expect(modal).toBeVisible();
    
    const timestamp = Date.now();
    const programName = `Test Program ${timestamp}`;
    
    await modal.getByPlaceholder(/name|название/i).fill(programName);
    
    // Submit
    await modal.getByRole('button', { name: /create|создать/i }).click();
    
    // Verify success (Toast might be flaky, so we also check the list)
    // await expect(page.getByText(/program created|программа создана/i)).toBeVisible();
    
    // Verify in list
    await expect(page.getByText(programName)).toBeVisible();
  });
});
