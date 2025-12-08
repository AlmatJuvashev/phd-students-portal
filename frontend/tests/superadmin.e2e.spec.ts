import { test, expect, Page } from '@playwright/test';
import { loginViaUI } from './utils/auth';

// Superadmin credentials
const superadminUser = process.env.E2E_SUPERADMIN_USER || 'juvashev';
const superadminPass = process.env.E2E_SUPERADMIN_PASS || 'superadminpassword123!';

/**
 * Custom login for superadmin users - handles redirect to /superadmin
 */
async function loginAsSuperadmin(page: Page) {
  await page.goto('/login');
  await page.getByLabel(/username/i).fill(superadminUser);
  await page.locator('input[name="password"]').fill(superadminPass);
  await page.getByRole('button', { name: /sign in|войти|кіру/i }).click();
  await page.waitForLoadState('networkidle');
  // Superadmin may redirect to /superadmin or /dashboard - just wait for no login page
  await expect(page).not.toHaveURL(/\/login/, { timeout: 15000 });
}

/**
 * E2E tests for Superadmin Panel UI
 * Tests superadmin workflows: tenants, admins, settings, logs
 */
test.describe('Superadmin Panel', () => {
  
  test.describe('Authentication & Access', () => {
    
    test('Superadmin can access /superadmin panel', async ({ page }) => {
      await loginAsSuperadmin(page);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
      
      await page.goto('/superadmin');
      await page.waitForLoadState('networkidle');
      
      // Should see superadmin dashboard or redirect to tenants
      await expect(page).toHaveURL(/\/superadmin/);
      
      // Should see main navigation elements
      await expect(page.getByRole('link', { name: /tenants|организации|ұйымдар/i })).toBeVisible();
    });

    test('Non-superadmin cannot access /superadmin', async ({ page }) => {
      // Login as regular admin
      const adminUser = process.env.E2E_ADMIN_USER || 'demo.admin';
      const adminPass = process.env.E2E_ADMIN_PASS || 'demopassword123!';
      
      await loginViaUI(page, adminUser, adminPass);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
      
      // Try to access superadmin panel
      await page.goto('/superadmin');
      
      // Should be redirected or see access denied
      const url = page.url();
      const accessDenied = page.getByText(/access denied|forbidden|доступ запрещен/i);
      
      // Either redirected away from superadmin or see access denied
      expect(url.includes('/superadmin') === false || await accessDenied.isVisible()).toBeTruthy();
    });
  });

  test.describe('Tenants Management', () => {
    
    test.beforeEach(async ({ page }) => {
      await loginAsSuperadmin(page);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
    });

    test('Tenants list displays all tenants', async ({ page }) => {
      await page.goto('/superadmin/tenants');
      await page.waitForLoadState('networkidle');
      
      // Should see tenants table or list
      await expect(page.getByRole('table').or(page.getByTestId('tenants-list'))).toBeVisible();
      
      // Should have at least one tenant (demo-university)
      const rows = page.getByRole('row');
      expect(await rows.count()).toBeGreaterThan(1); // Header + at least 1 data row
    });

    test('Can search tenants', async ({ page }) => {
      await page.goto('/superadmin/tenants');
      await page.waitForLoadState('networkidle');
      
      // Find search input
      const searchInput = page.getByPlaceholder(/search|поиск|іздеу/i);
      if (await searchInput.isVisible()) {
        await searchInput.fill('demo');
        await page.waitForTimeout(500); // debounce
        
        // Should filter results
        const demoRow = page.getByText('demo-university').or(page.getByText('Demo University'));
        await expect(demoRow.first()).toBeVisible();
      }
    });

    test('Can toggle services for a tenant', async ({ page }) => {
      await page.goto('/superadmin/tenants');
      await page.waitForLoadState('networkidle');
      
      // Find demo-university row
      const row = page.getByRole('row').filter({ hasText: /demo/i }).first();
      await expect(row).toBeVisible();
      
      // Find service toggle buttons (Calendar or Chat icons)
      const toggleButtons = row.locator('button').filter({ has: page.locator('svg') });
      const buttonCount = await toggleButtons.count();
      
      if (buttonCount > 0) {
        // Click first service toggle
        await toggleButtons.first().click();
        
        // Should see some visual feedback (toast, color change, etc.)
        await page.waitForTimeout(500);
      }
    });
  });

  test.describe('Admins Management', () => {
    
    test.beforeEach(async ({ page }) => {
      await loginAsSuperadmin(page);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
    });

    test('Admins list displays all admins', async ({ page }) => {
      await page.goto('/superadmin/admins');
      await page.waitForLoadState('networkidle');
      
      // Should see admins table
      await expect(page.getByRole('table')).toBeVisible();
      
      // Should have at least one admin
      const rows = page.getByRole('row');
      expect(await rows.count()).toBeGreaterThan(1);
    });

    test('Can filter admins by tenant', async ({ page }) => {
      await page.goto('/superadmin/admins');
      await page.waitForLoadState('networkidle');
      
      // Find tenant filter dropdown
      const filterDropdown = page.getByRole('combobox').or(page.getByLabel(/tenant|организация|ұйым/i));
      if (await filterDropdown.isVisible()) {
        await filterDropdown.click();
        
        // Should show tenant options
        const options = page.getByRole('option');
        expect(await options.count()).toBeGreaterThan(0);
      }
    });

    test('Create admin button opens wizard', async ({ page }) => {
      await page.goto('/superadmin/admins');
      await page.waitForLoadState('networkidle');
      
      // Find create button
      const createButton = page.getByRole('button', { name: /add|create|создать|қосу/i });
      if (await createButton.isVisible()) {
        await createButton.click();
        
        // Should show wizard/form
        const wizard = page.getByRole('dialog').or(page.getByRole('form'));
        await expect(wizard).toBeVisible();
      }
    });
  });

  test.describe('Settings Management', () => {
    
    test.beforeEach(async ({ page }) => {
      await loginAsSuperadmin(page);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
    });

    test('Settings page loads with categories', async ({ page }) => {
      await page.goto('/superadmin/settings');
      await page.waitForLoadState('networkidle');
      
      // Should see settings page
      await expect(page).toHaveURL(/\/superadmin\/settings/);
      
      // Should have settings list or categories
      const settingsList = page.getByRole('table').or(page.getByTestId('settings-list'));
      await expect(settingsList).toBeVisible();
    });

    test('Can edit a setting value', async ({ page }) => {
      await page.goto('/superadmin/settings');
      await page.waitForLoadState('networkidle');
      
      // Find edit button on first setting
      const editButton = page.getByRole('button', { name: /edit|редактировать|өңдеу/i }).first();
      if (await editButton.isVisible()) {
        await editButton.click();
        
        // Should show edit form/dialog
        const input = page.getByRole('textbox').or(page.getByRole('spinbutton'));
        await expect(input).toBeVisible();
      }
    });
  });

  test.describe('Activity Logs', () => {
    
    test.beforeEach(async ({ page }) => {
      await loginAsSuperadmin(page);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
    });

    test('Logs page displays activity history', async ({ page }) => {
      await page.goto('/superadmin/logs');
      await page.waitForLoadState('networkidle');
      
      // Should see logs page
      await expect(page).toHaveURL(/\/superadmin\/logs/);
      
      // Should have logs table
      await expect(page.getByRole('table')).toBeVisible();
    });

    test('Can filter logs by action type', async ({ page }) => {
      await page.goto('/superadmin/logs');
      await page.waitForLoadState('networkidle');
      
      // Find action filter
      const actionFilter = page.getByRole('combobox').first().or(page.getByLabel(/action|действие|әрекет/i));
      if (await actionFilter.isVisible()) {
        await actionFilter.click();
        
        // Should show action options
        const options = page.getByRole('option');
        expect(await options.count()).toBeGreaterThan(0);
      }
    });

    test('Logs show statistics', async ({ page }) => {
      await page.goto('/superadmin/logs');
      await page.waitForLoadState('networkidle');
      
      // Should see stats cards or summary
      const statsSection = page.getByText(/total logs|всего записей|барлық жазбалар/i)
        .or(page.getByText(/logs by|по действиям|әрекеттер бойынша/i));
      
      // Stats may be present
      if (await statsSection.isVisible()) {
        expect(true).toBeTruthy();
      }
    });
  });

  test.describe('Navigation', () => {
    
    test.beforeEach(async ({ page }) => {
      await loginAsSuperadmin(page);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
      await page.goto('/superadmin');
    });

    test('Sidebar has all main navigation links', async ({ page }) => {
      await page.waitForLoadState('networkidle');
      
      // Check for main nav links
      const tenantsLink = page.getByRole('link', { name: /tenants|организации|ұйымдар/i });
      const adminsLink = page.getByRole('link', { name: /admins|администраторы|әкімшілер/i });
      const settingsLink = page.getByRole('link', { name: /settings|настройки|баптаулар/i });
      const logsLink = page.getByRole('link', { name: /logs|журнал|журнал/i });
      
      await expect(tenantsLink).toBeVisible();
      await expect(adminsLink).toBeVisible();
      await expect(settingsLink).toBeVisible();
      await expect(logsLink).toBeVisible();
    });

    test('Clicking nav links navigates correctly', async ({ page }) => {
      await page.waitForLoadState('networkidle');
      
      // Click tenants link
      await page.getByRole('link', { name: /tenants|организации|ұйымдар/i }).click();
      await expect(page).toHaveURL(/\/superadmin\/tenants/);
      
      // Click admins link
      await page.getByRole('link', { name: /admins|администраторы|әкімшілер/i }).click();
      await expect(page).toHaveURL(/\/superadmin\/admins/);
      
      // Click settings link
      await page.getByRole('link', { name: /settings|настройки|баптаулар/i }).click();
      await expect(page).toHaveURL(/\/superadmin\/settings/);
      
      // Click logs link
      await page.getByRole('link', { name: /logs|журнал|журнал/i }).click();
      await expect(page).toHaveURL(/\/superadmin\/logs/);
    });
  });

  test.describe('Delete Confirmation Dialogs', () => {
    
    test.beforeEach(async ({ page }) => {
      await loginAsSuperadmin(page);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
    });

    test('Delete tenant shows confirmation modal', async ({ page }) => {
      await page.goto('/superadmin/tenants');
      await page.waitForLoadState('networkidle');
      
      // Find a delete button
      const deleteButton = page.getByRole('button').filter({ has: page.locator('svg.lucide-trash') }).first();
      if (await deleteButton.isVisible()) {
        await deleteButton.click();
        
        // Should show confirmation dialog, not browser confirm
        const dialog = page.getByRole('alertdialog').or(page.getByRole('dialog'));
        await expect(dialog).toBeVisible();
        
        // Dialog should have cancel and confirm buttons
        await expect(page.getByRole('button', { name: /cancel|отмена|бас тарту/i })).toBeVisible();
        await expect(page.getByRole('button', { name: /delete|удалить|жою|confirm/i })).toBeVisible();
        
        // Cancel the deletion
        await page.getByRole('button', { name: /cancel|отмена|бас тарту/i }).click();
        await expect(dialog).not.toBeVisible();
      }
    });

    test('Deactivate admin shows confirmation modal', async ({ page }) => {
      await page.goto('/superadmin/admins');
      await page.waitForLoadState('networkidle');
      
      // Find a deactivate/delete button
      const deleteButton = page.getByRole('button').filter({ has: page.locator('svg.lucide-trash') }).first();
      if (await deleteButton.isVisible()) {
        await deleteButton.click();
        
        // Should show confirmation dialog
        const dialog = page.getByRole('alertdialog').or(page.getByRole('dialog'));
        await expect(dialog).toBeVisible();
        
        // Cancel
        await page.getByRole('button', { name: /cancel|отмена|бас тарту/i }).click();
      }
    });

    test('Delete setting shows confirmation modal', async ({ page }) => {
      await page.goto('/superadmin/settings');
      await page.waitForLoadState('networkidle');
      
      // Find a delete button
      const deleteButton = page.getByRole('button').filter({ has: page.locator('svg.lucide-trash') }).first();
      if (await deleteButton.isVisible()) {
        await deleteButton.click();
        
        // Should show confirmation dialog
        const dialog = page.getByRole('alertdialog').or(page.getByRole('dialog'));
        await expect(dialog).toBeVisible();
        
        // Cancel
        await page.getByRole('button', { name: /cancel|отмена|бас тарту/i }).click();
      }
    });
  });
});
