import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';
import { createUserViaAPI, loginViaAPI } from './utils/api_auth';

const adminUser = process.env.E2E_ADMIN_USER || 'ta2087';
const adminPass = process.env.E2E_ADMIN_PASS || 'meadow-pluto-pioneer48';

test.describe('Chat Flow (Admin & Student)', () => {
  
  const roomName = `E2E Chat Room ${Date.now()}`;
  const newStudentEmail = `chatuser${Date.now()}@example.com`;
  const newStudentName = `Chat User ${Date.now()}`;

  test('Admin creates room and Student sends message', async ({ page, request }) => {
    // Increase timeout
    test.setTimeout(60000);

    // Login as Admin via API to get token
    const token = await loginViaAPI(request, adminUser, adminPass);

    // Create a fresh student
    // API returns { username, temp_password, ... }
    const createdUser = await createUserViaAPI(token, {
        email: newStudentEmail,
        first_name: newStudentName,
        last_name: 'Test',
        role: 'student'
    });

    // --- Admin Flow ---
    await loginViaUI(page, adminUser, adminPass);
    
    // Navigate to Admin Chat Rooms
    await page.goto('/admin/chat-rooms');
    
    // Create Room
    await expect(page.getByRole('heading', { name: /chat rooms|чаты/i })).toBeVisible();
    
    await page.getByPlaceholder(/PhD 2025|Public Health/i).fill(roomName);
    
    await page.getByRole('button', { name: /create|создать/i, exact: true }).click();
    
    // Verify room created
    await expect(page.getByText(roomName)).toBeVisible();
    
    // Add Student to Room
    const roomRow = page.locator('div.border.rounded-md').filter({ hasText: roomName });
    await roomRow.getByRole('button', { name: /members|участники/i }).click();
    
    // Modal opens
    const modal = page.locator('.fixed.z-50 > .relative');
    await expect(modal).toBeVisible();
    
    // Search for student by email (more reliable)
    await modal.getByPlaceholder(/search by name|поиск по имени/i).fill(newStudentEmail);
    
    // Wait for search results
    const userResult = modal.locator('div.border.rounded-md').filter({ hasText: newStudentEmail }).first();
    await expect(userResult).toBeVisible();
    
    // Click Add
    await userResult.getByRole('button', { name: /add|добавить/i }).click();
    
    // Close Modal
    // Target the button with X icon
    const closeButton = modal.locator('button').filter({ has: page.locator('svg.lucide-x') }).first();
    if (await closeButton.isVisible()) {
        await closeButton.click();
    } else {
        await page.keyboard.press('Escape');
    }
    await expect(modal).not.toBeVisible();
    
    // Logout Admin
    // Logout Admin
    // User Menu is the button with rounded-full avatar in the sticky header
    await page.locator('.sticky.top-0 button:has(.rounded-full)').first().click(); // Avatar
    await page.getByRole('menuitem', { name: /logout|выйти/i }).click();
    await expect(page).toHaveURL('/login');

    // --- Student Flow ---
    // Use the credentials returned by the API
    await loginViaUI(page, createdUser.username, createdUser.temp_password);
    
    // Navigate to Chat
    await page.locator('a[href="/chat"]').click();
    
    // Verify Room is visible
    await expect(page.getByText(roomName)).toBeVisible();
    
    // Select Room
    await page.getByText(roomName).click();
    
    // Send Message
    const messageText = `Hello from Student in ${roomName}`;
    await page.getByPlaceholder(/type a message|введите сообщение/i).fill(messageText);
    await page.locator('button:has(svg.lucide-send)').click();
    
    // Verify Message
    await expect(page.getByText(messageText)).toBeVisible();
  });
});
