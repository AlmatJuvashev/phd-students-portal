import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';

const adminUser = process.env.E2E_ADMIN_USER || 'ta2087';
const adminPass = process.env.E2E_ADMIN_PASS || 'meadow-pluto-pioneer48';

test.describe('User Management (Admin)', () => {
  
  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, adminUser, adminPass);
    // Navigate to Admin Dashboard -> Users
    await page.getByRole('link', { name: /admin|админ/i }).click();
    await page.getByRole('link', { name: /users|пользователи|пайдаланушылар/i }).click();
  });

  test('Create New Student', async ({ page }) => {
    // 1. Click "Create User" or "Create Student"
    // Assuming there's a button to add a user/student
    // Based on layout.tsx/common.json, it might be "Create User" or specific "Create Student" tab
    
    // Check if we are on the users page
    // Check if we are on the users page
    await expect(page.getByRole('heading', { name: /user management|управление пользователями/i, level: 1 })).toBeVisible();
    
    // Check if Students tab is active (default) and its heading is visible
    await expect(page.getByRole('heading', { name: /manage students|управление студентами/i, level: 2 })).toBeVisible();

    // Click "Create Student" button
    await page.getByRole('button', { name: /create student|создать студента/i }).click();
    
    // Fill form
    const timestamp = Date.now();
    const newEmail = `test_student_${timestamp}@example.com`;
    const newUsername = `student_${timestamp}`;
    
    // Scope to modal (custom implementation, no dialog role)
    const modal = page.locator('.fixed.z-50 > .relative');
    await expect(modal).toBeVisible();

    await modal.getByPlaceholder(/email/i).fill(newEmail);
    await modal.getByPlaceholder(/first name|имя/i).fill('Test');
    await modal.getByPlaceholder(/last name|фамилия/i).fill('User');
    
    // Handle Selects (Program, Specialty, Department, Cohort)
    // Helper to select first option
    const selectFirstOption = async (placeholderRegex: RegExp) => {
        await modal.getByRole('combobox').filter({ hasText: placeholderRegex }).click();
        // Wait for options and click first one
        await page.getByRole('option').first().click();
    };

    // We need to select in order because they might be dependent (Program -> Specialty)
    // But in the code:
    // Program -> Specialty (filtered by program)
    // Department -> Cohort (filtered by department)
    
    // Note: If no options exist, this will fail. We might need to seed data.
    try {
        await selectFirstOption(/program|программа/i);
        await selectFirstOption(/specialty|специальность/i);
        await selectFirstOption(/department|кафедра/i);
        await selectFirstOption(/cohort|поток/i);
    } catch (e) {
        console.log('Failed to select options. Dictionaries might be empty.');
        throw e;
    }

    // Submit
    await modal.getByRole('button', { name: /create student|создать студента/i }).click();

    // Verify success message or appearance in table
    await expect(page.getByText(/created|создан/i)).toBeVisible();
    
    // Verify user in list
    await page.getByPlaceholder(/search|поиск/i).fill(newEmail);
    await expect(page.getByText(newEmail)).toBeVisible();
  });
});
