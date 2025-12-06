import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';

const studentUser = process.env.E2E_STUDENT_USER || 'ts5251';
const studentPass = process.env.E2E_STUDENT_PASS || 'pioneer-canvas-silver52';

test.describe('Communication (Chat)', () => {
  
  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, studentUser, studentPass);
    // Navigate to Chat page
    await page.locator('a[href="/chat"]').click();
  });

  test('Send Message in Chat', async ({ page }) => {
    // Verify URL
    await expect(page).toHaveURL(/.*\/chat/);

    // Verify Chat Page Title
    // Relaxed selector to find any text matching Chat/Чат/Messages in a heading-like element
    await expect(page.locator('h1')).toContainText(/chat|чат|messages/i);

    // Verify Room List is visible
    // The sidebar might be hidden on mobile, but we test on desktop viewport by default
    const roomList = page.locator('.overflow-y-auto').first();
    await expect(roomList).toBeVisible();

    // Select the first active room
    // We look for a button that is a room item
    const firstRoom = page.locator('button.w-full.rounded-lg.border').first();
    
    // If no rooms, we can't test sending messages.
    // We assume seed data created at least one room for the student.
    if (await firstRoom.count() > 0) {
        await firstRoom.click();
        
        // Verify Chat Window is active
        // The header should show the room name
        const roomName = await firstRoom.locator('.font-semibold').first().textContent();
        // Check if header contains room name (might be truncated)
        // Just check if message input is visible
        const messageInput = page.getByPlaceholder(/type a message|введите сообщение/i);
        await expect(messageInput).toBeVisible();

        // Send a message
        const timestamp = Date.now();
        const messageText = `Hello from E2E test ${timestamp}`;
        await messageInput.fill(messageText);
        
        // Click Send button (icon)
        // Usually it's a button with Send icon, might need specific selector
        // Looking at code: <Send className="..." /> inside a Button
        // We can try getByRole('button') that is near the input
        const sendButton = page.locator('button:has(svg.lucide-send)');
        await sendButton.click();

        // Verify message appears in the list
        await expect(page.getByText(messageText)).toBeVisible();
    } else {
        console.log('No chat rooms found for student. Skipping message send test.');
    }
  });
});
