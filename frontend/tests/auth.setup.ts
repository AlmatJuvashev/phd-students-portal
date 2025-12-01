import { test as setup, expect } from '@playwright/test';

const authFile = 'playwright/.auth/user.json';

setup('authenticate', async ({ page }) => {
  // Perform authentication steps. Replace these actions with your own.
  await page.goto('/login');
  await page.getByLabel('Username').fill('tu6260');
  await page.locator('#password').fill('thunder-pluto-river72');
  await page.locator('button[type="submit"]').click();
  // Wait until the page receives the cookies.
  //
  // Sometimes login flow sets cookies in the process of several redirects.
  // Wait for the final URL to ensure that the cookies are actually set.
  await page.waitForLoadState('networkidle');
  // Alternatively, you can wait until the page reaches a state where all cookies are set.
  await expect(page.getByText('Student Dashboard')).toBeVisible();

  // End of authentication steps.

  await page.context().storageState({ path: authFile });
});
