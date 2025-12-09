import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';

const studentUser = process.env.E2E_STUDENT_USER || 'ts5251';
const studentPass = process.env.E2E_STUDENT_PASS || 'pioneer-canvas-silver52';

test.describe('Student Journey', () => {
  
  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, studentUser, studentPass);
    // Navigate to Journey page
    await page.getByRole('link', { name: /journey|путь|жол/i }).click();
  });

  test('View Journey Steps', async ({ page }) => {
    // Verify we are on the Journey page (WorldMap title)
    await expect(page.getByRole('heading', { name: /my dissertation map|карта диссертации/i })).toBeVisible();
    
    // Verify steps (Worlds) are listed
    // "I — Preparation" or "I — Подготовка"
    await expect(page.getByText(/preparation|подготовка/i)).toBeVisible();
    
    // Verify first node is visible (Doctoral profile)
    await expect(page.getByText(/doctoral profile|профиль докторанта/i)).toBeVisible();
  });

  test('Interact with a Step', async ({ page }) => {
    // Click on the first node (Doctoral profile)
    // It's a button or clickable element.
    // In WorldMap.tsx, NodeToken is clickable.
    
    const profileNode = page.getByText(/doctoral profile|профиль докторанта/i);
    await profileNode.click();
    
    // Verify details sheet opens
    // NodeDetailsSheet uses SheetContent which usually has a title matching the node title
    await expect(page.getByRole('heading', { name: /doctoral profile|профиль докторанта/i, level: 2 })).toBeVisible();
    
    // Check if form fields are visible (e.g., Full Name)
    // "Full Name" / "ФИО"
    // Note: The sheet might be rendering the form fields.
    // S1_profile has "full_name" field.
    // We can check for label "Full Name" or "ФИО"
    // But wait, the sheet title might be enough for now.
  });
});
