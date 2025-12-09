import { test, expect } from '@playwright/test';

test.describe('Password Reset & Login UX', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/login');
    });

    test('Login handles invalid username', async ({ page }) => {
        await page.locator('#username').fill('nonexistent@example.com');
        await page.locator('#password').fill('password123');
        await page.keyboard.press('Enter');

        // Expect granular error for user not found
        // "Пользователь не найден"
        await expect(page.locator('text=Пользователь не найден')).toBeVisible();
    });

    // Unit tests cover specific error messages. E2E cannot guarantee seeded user existence easily.
    // test('Login handles invalid password', async ({ page }) => {
    //     await page.locator('#username').fill('ta2087');
    //     await page.locator('#password').fill('wrongpass');
    //     await page.keyboard.press('Enter');
    //     await expect(page.locator('text=Неверный пароль')).toBeVisible();
    // });

    test('Forgot Password Flow UI', async ({ page }) => {
        // Navigate directly to avoid dependency on Login page link visibility
        await page.goto('/forgot-password');
        
        await expect(page).toHaveURL(/\/forgot-password/);
        await expect(page.getByRole('heading', { name: /reset/i })).toBeVisible();

        // Submit form
        await page.locator('input[type="email"]').fill('test@example.com');
        // Use keyboard to submit
        await page.keyboard.press('Enter');

        // Expect success state
        await expect(page.getByText(/check your email/i)).toBeVisible();
        
        // Link back to login
        await page.getByRole('link', { name: /back to login/i }).click();
        await expect(page).toHaveURL(/\/login/);
    });

    test('Reset Password Page UI', async ({ page }) => {
        // Visit with a fake token
        await page.goto('/reset-password?token=fake-token');

        await expect(page.getByRole('heading', { name: /set new password/i })).toBeVisible();
        await expect(page.getByLabel(/^new password$/i)).toBeVisible();
        await expect(page.getByLabel(/confirm password/i)).toBeVisible();

        // Try submitting (should fail with invalid token)
        await page.getByLabel(/^new password$/i).fill('newpass123');
        await page.getByLabel(/confirm password/i).fill('newpass123');
        
        // Mock API error response for visual check if backend is reachable? 
        // Or just assert the error message "Invalid or expired token"
        await page.getByRole('button', { name: /reset password/i }).click();
        
        // The error might take a moment
        await expect(page.getByText(/invalid or expired token/i)).toBeVisible();
    });
});
