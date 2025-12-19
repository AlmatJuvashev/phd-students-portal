import { test, expect } from '@playwright/test';

const adminUser = process.env.E2E_ADMIN_USER || 'ta2087';
const adminPass = process.env.E2E_ADMIN_PASS || 'meadow-pluto-pioneer48';

test.describe('Authentication Security', () => {

  test('Login sets HttpOnly cookie and NOT localStorage', async ({ page, context }) => {
    await page.goto('/login');
    
    // Perform Login
    await page.getByLabel(/username|name|email/i).fill(adminUser);
    await page.locator('input[name="password"]').fill(adminPass);
    await page.getByRole('button', { name: /login|войти|кіру/i }).click();

    // Verify Redirect
    await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });

    // 1. Verify Cookie
    const cookies = await context.cookies();
    const jwtCookie = cookies.find(c => c.name === 'jwt_token');
    
    expect(jwtCookie).toBeDefined();
    expect(jwtCookie?.httpOnly).toBe(true);
    expect(jwtCookie?.secure).toBe(true); // Should be true if running on https, or localhost might default? Chrome treats localhost as secure context.
    // Note: Secure flag depends on server config and connection, might be false on http://localhost

    // 2. Verify LocalStorage is empty of token
    const localStorageToken = await page.evaluate(() => localStorage.getItem('token'));
    expect(localStorageToken).toBeNull();
  });

  test('Logout clears cookie', async ({ page, context }) => {
    // Login first
    await page.goto('/login');
    await page.getByLabel(/username|name|email/i).fill(adminUser);
    await page.locator('input[name="password"]').fill(adminPass);
    await page.getByRole('button', { name: /login|войти|кіру/i }).click();
    await expect(page).not.toHaveURL(/\/login/);

    // Logout via UI
    // Assuming UI has a logout button in header/profile menu
    // Open profile menu
    const profileBtn = page.locator('header button').last(); 
    await profileBtn.click(); // heuristic, might need update if UI changes
    await page.getByText(/logout|выйти|шығу/i).click();

    await expect(page).toHaveURL(/\/login/);

    // Verify Cookie is gone or expired
    const cookies = await context.cookies();
    const jwtCookie = cookies.find(c => c.name === 'jwt_token');
    // It might be removed entirely OR set to expire in past.
    // If removed: undefined. If expired: Playwright context.cookies() might not show it?
    // Usually backend sets it to empty string and MaxAge -1, so browser deletes it.
    expect(jwtCookie).toBeUndefined();
  });

});
