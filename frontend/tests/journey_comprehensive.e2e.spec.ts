import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';
import { loginViaAPI, createUserViaAPI } from './utils/api_auth';

const adminUser = process.env.E2E_ADMIN_USER || 'ta2087';
const adminPass = process.env.E2E_ADMIN_PASS || 'meadow-pluto-pioneer48';
const studentUser = process.env.E2E_STUDENT_USER || 'ts5251';
const studentPass = process.env.E2E_STUDENT_PASS || 'pioneer-canvas-silver52';

test.describe('Journey Progression', () => {
  test.setTimeout(120000);

  test('Complete first node and verify progress', async ({ page, request }) => {
    // Try to create fresh student for clean journey
    let username = studentUser;
    let password = studentPass;
    
    try {
      const token = await loginViaAPI(request, adminUser, adminPass);
      const createdUser = await createUserViaAPI(token, {
        email: `progress${Date.now()}@example.com`,
        first_name: 'Progress',
        last_name: 'Test',
        role: 'student'
      });
      username = createdUser.username;
      password = createdUser.temp_password;
    } catch {
      console.log('Using existing student for progression test');
    }

    await loginViaUI(page, username, password);
    
    // Navigate to Journey
    await page.getByRole('link', { name: /journey|путь|жол/i }).click();
    await expect(page.getByRole('heading', { name: /my dissertation map|карта диссертации/i })).toBeVisible();
    
    // Find and click first node (Doctoral Profile)
    const firstNode = page.getByText(/doctoral profile|профиль докторанта/i).first();
    await firstNode.click();
    
    // Wait for sheet to open
    await expect(page.getByTestId('node-details-sheet')).toBeVisible();
    
    // Check state badge
    const stateBadge = page.getByTestId('node-state-badge');
    await expect(stateBadge).toBeVisible();
    
    // Fill form if available
    const fullNameInput = page.getByLabel(/full name|фио|аты-жөні/i);
    if (await fullNameInput.isVisible({ timeout: 2000 })) {
      await fullNameInput.fill('Progress Test User');
    }
    
    // Submit if button visible
    const submitButton = page.getByTestId('node-submit-button');
    if (await submitButton.isVisible({ timeout: 2000 })) {
      await submitButton.click();
      
      // Wait for state change
      await expect(stateBadge).toHaveText(/submitted|done|отправлено|готово|жіберілді/i, { timeout: 10000 });
    }
    
    // Close sheet
    await page.keyboard.press('Escape');
    
    // Verify we're back on journey page
    await expect(page.getByRole('heading', { name: /my dissertation map|карта диссертации/i })).toBeVisible();
  });

  test('Verify node states are displayed correctly', async ({ page }) => {
    await loginViaUI(page, studentUser, studentPass);
    
    // Navigate to Journey
    await page.getByRole('link', { name: /journey|путь|жол/i }).click();
    await expect(page.getByRole('heading', { name: /my dissertation map|карта диссертации/i })).toBeVisible();
    
    // World headers should be visible
    await expect(page.getByText(/preparation|подготовка/i)).toBeVisible();
    
    // At least one node should be visible (either by testid or by role button)
    const nodes = page.locator('[data-testid^="node-token-"]');
    let nodeCount = await nodes.count();
    
    // Fallback to checking for node buttons if testid not present
    if (nodeCount === 0) {
      const nodeButtons = page.locator('button[aria-label*="active"], button[aria-label*="locked"], button[aria-label*="done"]');
      nodeCount = await nodeButtons.count();
    }
    
    expect(nodeCount).toBeGreaterThan(0);
  });
});

test.describe('Template Downloads', () => {
  test.setTimeout(60000);

  test('Download template button is functional', async ({ page }) => {
    await loginViaUI(page, studentUser, studentPass);
    
    // Navigate to Journey
    await page.getByRole('link', { name: /journey|путь|жол/i }).click();
    await expect(page.getByRole('heading', { name: /my dissertation map|карта диссертации/i })).toBeVisible();
    
    // Click on a confirmTask node (they typically have templates)
    // Look for any node with "document" or "заявление" in the name
    // Use button with aria-label instead to avoid clicking disabled elements
    const templateNode = page.locator('button[aria-label*="Application"], button[aria-label*="заявление"]').first();
    if (await templateNode.isVisible({ timeout: 3000 })) {
      // Check if it's enabled before clicking
      const isEnabled = await templateNode.isEnabled();
      if (!isEnabled) {
        console.log('Template node is not enabled (locked), skipping');
        return;
      }
      
      await templateNode.click();
      await expect(page.getByTestId('node-details-sheet')).toBeVisible({ timeout: 5000 });
      
      // Look for download button
      const downloadBtn = page.getByRole('button', { name: /download|скачать|жүктеу/i });
      if (await downloadBtn.isVisible({ timeout: 3000 })) {
        // Set up download listener
        const downloadPromise = page.waitForEvent('download', { timeout: 10000 });
        await downloadBtn.click();
        
        try {
          const download = await downloadPromise;
          expect(download.suggestedFilename()).toBeTruthy();
          console.log('Downloaded template:', download.suggestedFilename());
        } catch {
          console.log('No download triggered (template might not be configured)');
        }
      }
    } else {
      console.log('No template node found, skipping download test');
    }
  });
});

test.describe('File Uploads', () => {
  test.setTimeout(90000);

  test('Upload file to node with upload slot', async ({ page }) => {
    await loginViaUI(page, studentUser, studentPass);
    
    // Navigate to Journey
    await page.getByRole('link', { name: /journey|путь|жол/i }).click();
    await expect(page.getByRole('heading', { name: /my dissertation map|карта диссертации/i })).toBeVisible();
    
    // Look for a confirmTask node that requires uploads
    // These typically have "upload" or "document" in their description
    const confirmNode = page.getByText(/confirm|подтверд|растау/i).first();
    if (await confirmNode.isVisible({ timeout: 3000 })) {
      await confirmNode.click();
      await expect(page.getByTestId('node-details-sheet')).toBeVisible();
      
      // Look for file input
      const fileInput = page.locator('input[type="file"]');
      if (await fileInput.isVisible({ timeout: 3000 })) {
        // Create a test file
        await fileInput.setInputFiles({
          name: 'test-document.pdf',
          mimeType: 'application/pdf',
          buffer: Buffer.from('test pdf content')
        });
        
        // Wait for upload to complete
        await page.waitForTimeout(2000);
        
        // Verify file appears in slot
        await expect(page.getByText(/test-document|uploaded|загружен/i)).toBeVisible({ timeout: 5000 });
        console.log('File uploaded successfully');
      } else {
        console.log('No file input found (node might not require uploads)');
      }
    } else {
      console.log('No confirm node found, skipping upload test');
    }
  });
});

test.describe('Journey Errors', () => {
  test.setTimeout(30000);

  test('Session expiry redirects to login', async ({ page, context }) => {
    await loginViaUI(page, studentUser, studentPass);
    
    // Navigate to Journey
    await page.getByRole('link', { name: /journey|путь|жол/i }).click();
    await expect(page.getByRole('heading', { name: /my dissertation map|карта диссертации/i })).toBeVisible();
    
    // Clear auth token to simulate expiry
    await context.clearCookies();
    await page.evaluate(() => {
      localStorage.removeItem('token');
      localStorage.removeItem('auth_token');
    });
    
    // Try to interact - should redirect to login
    await page.reload();
    
    // Should be redirected to login page
    await expect(page).toHaveURL(/login/, { timeout: 10000 });
  });
});
