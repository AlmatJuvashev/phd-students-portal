import { test, expect } from '@playwright/test';

test.describe('Leaderboard Feature', () => {
  test.beforeEach(async ({ page }) => {
    // Login as student
    // Assuming standard login flow or mock
    await page.goto('/login');
    await page.fill('input[name="username"]', 'demo.student1');
    await page.fill('input[type="password"]', 'demopassword123!');
    await page.click('button[type="submit"]');
    await page.waitForURL((url) => url.pathname === '/' || url.pathname === '/journey');
    
    // Navigate to journey for leaderboard test
    await page.goto('/journey');
  });

  test('should open leaderboard modal and display data', async ({ page }) => {
    // Navigate to Map if not already there
    // Check if we are on map
    await expect(page.locator('[data-testid="journey-map"]')).toBeVisible();

    // Click Leaderboard Trigger (Cup Icon/Total XP info)
    // The trigger has title="View Leaderboard"
    await page.click('button[title="View Leaderboard"]');

    // Verify Modal Opens
    const modal = page.locator('text=Leaderboard'); // "Scoreboard" or "Leaderboard" based on translations
    // Title in ScoreboardModal is t('scoreboard.title') -> "LEADERBOARD" or similar.
    // Let's look for "h2" with text content. To be robust, we look for key elements.
    await expect(page.locator('h2')).toContainText(/Leaderboard|Рейтинг/i);

    // Verify User Stats Summary exists
    await expect(page.locator('text=Your Position')).toBeVisible();
    await expect(page.locator('text=Total XP')).toBeVisible(); // Might be "Your XP" based on code

    // Verify "The Chase" list is present
    await expect(page.locator('text=THE CHASE')).toBeVisible();

    // Verify at least one entry (Me) is there
    // We look for the "You" badge
    await expect(page.getByText('YOU', { exact: true })).toBeVisible();

    // Close Modal
    await page.click('button:has(svg.lucide-x)'); // specific close button
    await expect(page.locator('h2')).not.toBeVisible();
  });
});
