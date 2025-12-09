import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';

const adminUser = process.env.E2E_ADMIN_USER || 'ta2087';
const adminPass = process.env.E2E_ADMIN_PASS || 'meadow-pluto-pioneer48';

test.describe('Calendar Management', () => {
  
  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, adminUser, adminPass);
    // Navigate to Calendar page
    await page.getByRole('link', { name: /calendar|календарь|күнтізбе/i }).click();
  });

  test('Create New Event', async ({ page }) => {
    // Verify Calendar Page Title
    await expect(page.getByRole('heading', { name: /calendar|календарь/i, level: 1 })).toBeVisible();

    // Click "New Event" button
    await page.getByRole('button', { name: /new event|новое событие|жаңа оқиға/i }).click();

    // Verify Modal opens
    const modal = page.getByRole('dialog');
    await expect(modal).toBeVisible();

    // Fill Form
    const timestamp = Date.now();
    const eventTitle = `Test Event ${timestamp}`;
    
    await modal.getByLabel(/title|название/i).fill(eventTitle);
    
    // Dates are pre-filled, but we can modify them if needed.
    // For simplicity, we keep default dates (next hour).
    
    // Select Event Type (optional, defaults to Meeting)
    // Select Meeting Type (Online/Offline) - defaults to Offline
    
    // Fill Description
    await modal.getByLabel(/description|описание/i).fill('This is an E2E test event');

    // Save
    await modal.getByRole('button', { name: /save|сохранить/i }).click();

    // Verify Modal closes
    await expect(modal).not.toBeVisible();

    // Verify Event appears on Calendar
    // react-big-calendar renders events as divs with text
    await expect(page.getByText(eventTitle)).toBeVisible();
  });
});
