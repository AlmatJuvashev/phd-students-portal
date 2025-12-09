import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';

const studentUser = process.env.E2E_STUDENT_USER || 'ts5251';
const studentPass = process.env.E2E_STUDENT_PASS || 'pioneer-canvas-silver52';

test.describe('Journey Map Conditions', () => {
  
  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, studentUser, studentPass);
    console.log('Current URL:', page.url());
    console.log('Page Title:', await page.title());
    await page.screenshot({ path: 'debug_after_login.png' });
    
    // Listen for console errors
    page.on('console', msg => {
      if (msg.type() === 'error') console.log(`Console Error: "${msg.text()}"`);
    });

    await page.getByRole('link', { name: /journey|путь|жол/i }).click();
  });

  test('RP Required when graduation > 3 years', async ({ page }) => {
    // 1. Open Profile Node
    const profileNode = page.getByText(/doctoral profile|профиль докторанта/i);
    await expect(profileNode).toBeVisible();
    await profileNode.click();
    
    // Wait for sheet to open
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    
    // Wait for loading to finish (text "Loading..." should disappear)
    await expect(dialog.getByText(/loading|загрузка/i)).not.toBeVisible({ timeout: 10000 });
    
    // 2. Set Graduation Date to 4 years ago
    // 2. Set Graduation Date to 4 years ago
    // Click the date picker trigger
    await page.getByRole('button', { name: /select a date|выберите дату/i }).click();
    
    // Navigate back 4 years. Assuming the calendar has a "Previous Year" button or we can navigate by months.
    // If it's standard shadcn/ui calendar, it might be hard to navigate years quickly without a year dropdown.
    // Let's try to just go back enough months or find a year selector.
    // Alternatively, we can type the date if the input is accessible.
    // But usually shadcn date picker is read-only input.
    
    // Let's try to find the year selector if it exists, or just click "Previous Year" 4 times if available?
    // Standard DayPicker usually has "Previous Month". 4 years = 48 clicks. That's too many.
    // Maybe there is a year dropdown?
    // Let's assume we can just pick *any* date in the past for now to see if it works, 
    // but we need > 3 years.
    
    // Hack: If we can't easily navigate years, we might need to use `evaluate` to set the value in the form state directly if possible.
    // But this is E2E.
    
    // Let's try to click the "Year" in the calendar header if it allows switching views.
    // Or check if there's a year dropdown.
    
    // If we can't easily select a date 4 years ago, let's try to set the system time? No.
    
    // Let's try to type into the input if possible.
    // The previous error said "Input Name: null", so maybe the input is hidden.
    
    // Let's try to use `page.evaluate` to find the hidden input and set value?
    // Or better, assume the calendar has a year navigation.
    
    // Let's try to just select the *current* date first to verify interaction, then figure out the year.
    // Actually, for "RP Not Required", < 3 years. Today is fine.
    // For "RP Required", > 3 years.
    
    // Let's try to click the "Previous Month" button 48 times? It's slow but might work.
    // Selector for prev month: `button[name="previous-month"]` or `aria-label="Go to previous month"`.
    
    // Let's try a loop for 4 years (approx 48 months).
    // But first, let's just try to pick a date in the past.
    
    // Let's try to find if there is a year select.
    // If not, I'll use a loop.
    
    const prevMonthBtn = page.getByLabel('Go to previous month');
    if (await prevMonthBtn.isVisible()) {
        for (let i = 0; i < 40; i++) { // 3+ years
            await prevMonthBtn.click();
        }
    }
    
    // Click a day (e.g., 15th)
    await page.getByRole('gridcell', { name: '15' }).first().click();
    
    // 3. Save
    await page.getByTestId('node-submit-button').click();
    
    // Close the sheet (click outside or close button)
    // If save closes it, we don't need to.
    // Usually "Save and Submit" closes it.
    // Wait for potential success message or just wait a bit
    await page.waitForTimeout(1000);

    // Close the sheet explicitly
    await page.keyboard.press('Escape');

    // Wait for dialog to disappear.
    await expect(page.getByTestId('node-details-sheet')).not.toBeVisible();
 
    // 4. Verify RP Node is Visible
    // "Research Proposal actualization" / "Актуализация Research Proposal"
    await expect(page.getByText(/research proposal actualization|актуализация research proposal/i)).toBeVisible();
  });

  test('RP Not Required when graduation < 3 years', async ({ page }) => {
    // 1. Open Profile Node
    await page.getByText(/doctoral profile|профиль докторанта/i).click();
    
    // Wait for sheet
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    await expect(dialog.getByText(/loading|загрузка/i)).not.toBeVisible({ timeout: 10000 });
    
    // 2. Set Graduation Date to 1 year ago (or just today/recent)
    await page.getByRole('button', { name: /select a date|выберите дату/i }).click();
    
    // Just go back a few months
    const prevMonthBtn = page.getByLabel('Go to previous month');
    if (await prevMonthBtn.isVisible()) {
        for (let i = 0; i < 5; i++) {
            await prevMonthBtn.click();
        }
    }
    
    // Click a day
    await page.getByRole('gridcell', { name: '15' }).first().click();
    
    // 3. Save
    await page.getByTestId('node-submit-button').click();
    
    // Wait for potential success message or just wait a bit
    await page.waitForTimeout(1000);
    
    // Close the sheet explicitly if it doesn't close automatically
    await page.keyboard.press('Escape');
    
    // Wait for dialog to close
    await expect(page.getByTestId('node-details-sheet')).not.toBeVisible();

    // 4. Verify RP Node is Hidden
    await expect(page.getByText(/research proposal actualization|актуализация research proposal/i)).not.toBeVisible();
    
    // 5. Verify Next Stage (Normokontrol) is Visible (if it's the alternative)
    // "Normokontrol" / "Нормоконтроль"
    // D1_normokontrol_ncste
    // await expect(page.getByText(/normokontrol|нормоконтроль/i)).toBeVisible(); 
  });
});
