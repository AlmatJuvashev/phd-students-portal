import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';

const tenantBAdmin = { user: 'admin.b', pass: 'demopassword123!' };
const tenantBStudent = { user: 'student.b', pass: 'demopassword123!', name: 'Student B' };
const tenantAStudent = { user: 'demo.student1', name: 'Emma Brown' }; // From demo seed

test.describe('Multitenancy Isolation', () => {
  test('Admin B should see Student B but NOT Tenant A students', async ({ page }) => {
    // 1. Login as Admin B
    await loginViaUI(page, tenantBAdmin.user, tenantBAdmin.pass);

    // 2. Navigation
    // Ensure we are on dashboard then go to Students
    // Or direct navigation if allowed
    await page.goto('/admin/students');
    // Ensure table/list loaded
    // Wait for at least one row or empty state?
    // Student B should be there, so wait for it.
    
    // 3. Verify Student B is visible
    await expect(page.getByText(tenantBStudent.name)).toBeVisible({ timeout: 10000 });

    // 4. Verify Tenant A Student is NOT visible
    await expect(page.getByText(tenantAStudent.name)).not.toBeVisible();

    // 5. Verify Total Count if displayed (Optional)
    // "Found X students"
  });
});
