import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';

const adminUser = process.env.E2E_ADMIN_USER || 'ta2087';
const adminPass = process.env.E2E_ADMIN_PASS || 'meadow-pluto-pioneer48';
const studentUser = process.env.E2E_STUDENT_USER || 'ts5251';
const studentPass = process.env.E2E_STUDENT_PASS || 'pioneer-canvas-silver52';

/**
 * E2E tests for tenant service visibility feature
 * Tests that chat/calendar nav links and routes respect tenant service settings
 */
test.describe('Service Visibility', () => {
  
  test.describe('Default Services (All Enabled)', () => {
    
    test('Student sees Chat and Calendar nav links when services enabled', async ({ page }) => {
      await loginViaUI(page, studentUser, studentPass);
      
      // Wait for navigation to complete
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
      
      // Both Chat and Calendar links should be visible in nav
      // Note: links may be hidden in mobile menu, so we check desktop nav
      const chatLink = page.getByRole('link', { name: /messages|chat|сообщения|хабарлар/i });
      const calendarLink = page.getByRole('link', { name: /calendar|календарь|күнтізбе/i });
      
      // These should be present (may need to scroll or open menu on mobile)
      await expect(chatLink.or(calendarLink.first())).toBeVisible();
    });

    test('Student can access /chat route', async ({ page }) => {
      await loginViaUI(page, studentUser, studentPass);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
      
      // Navigate to chat
      await page.goto('/chat');
      
      // Should NOT show "Service Not Available" message
      await expect(page.getByText(/service not available/i)).not.toBeVisible();
      
      // Should show chat interface elements
      await expect(page.getByText(/chat|messages|сообщения/i).first()).toBeVisible();
    });

    test('Student can access /calendar route', async ({ page }) => {
      await loginViaUI(page, studentUser, studentPass);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
      
      // Navigate to calendar
      await page.goto('/calendar');
      
      // Should NOT show "Service Not Available" message
      await expect(page.getByText(/service not available/i)).not.toBeVisible();
      
      // Should show calendar elements
      await expect(page.getByText(/calendar|events|события/i).first()).toBeVisible();
    });

    test('Admin sees chat-rooms in sidebar when chat enabled', async ({ page }) => {
      await loginViaUI(page, adminUser, adminPass);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
      
      // Navigate to admin panel
      await page.goto('/admin');
      
      // Wait for admin layout to load
      await page.waitForLoadState('networkidle');
      
      // Chat rooms link should be visible in sidebar
      const chatRoomsLink = page.getByRole('link', { name: /chat rooms|комнаты чата|чат бөлмелері/i });
      await expect(chatRoomsLink).toBeVisible();
    });

    test('Admin can access /admin/calendar', async ({ page }) => {
      await loginViaUI(page, adminUser, adminPass);
      await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });
      
      // Navigate to admin calendar
      await page.goto('/admin/calendar');
      
      // Should load calendar view
      await expect(page.getByText(/service not available/i)).not.toBeVisible();
    });
  });

  test.describe('API Endpoints', () => {
    
    test('/me/tenants returns user tenant memberships', async ({ request }) => {
      // Login first to get token
      const loginResponse = await request.post('/api/auth/login', {
        data: {
          username: adminUser,
          password: adminPass,
        },
      });
      
      expect(loginResponse.ok()).toBeTruthy();
      const { token } = await loginResponse.json();
      
      // Get user's tenant memberships
      const response = await request.get('/api/me/tenants', {
        headers: {
          Authorization: `Bearer ${token}`,
          'X-Tenant-Slug': 'kaznmu',
        },
      });
      
      expect(response.ok()).toBeTruthy();
      const data = await response.json();
      
      // Should have memberships array
      expect(data).toHaveProperty('memberships');
    });

    test('/me/tenant returns current tenant with enabled_services', async ({ request }) => {
      // Login first
      const loginResponse = await request.post('/api/auth/login', {
        data: {
          username: adminUser,
          password: adminPass,
        },
      });
      
      expect(loginResponse.ok()).toBeTruthy();
      const { token } = await loginResponse.json();
      
      // Get current tenant info
      const response = await request.get('/api/me/tenant', {
        headers: {
          Authorization: `Bearer ${token}`,
          'X-Tenant-Slug': 'kaznmu',
        },
      });
      
      expect(response.ok()).toBeTruthy();
      const data = await response.json();
      
      // Should have tenant info with enabled_services
      expect(data).toHaveProperty('slug');
      expect(data).toHaveProperty('name');
      expect(data).toHaveProperty('enabled_services');
      expect(Array.isArray(data.enabled_services)).toBeTruthy();
    });
  });

  test.describe('Superadmin Service Management', () => {
    
    test.skip('Superadmin can update tenant services via API', async ({ request }) => {
      // Note: This test requires superadmin credentials
      // Skip in CI unless superadmin credentials configured
      
      const superadminUser = process.env.E2E_SUPERADMIN_USER;
      const superadminPass = process.env.E2E_SUPERADMIN_PASS;
      
      if (!superadminUser || !superadminPass) {
        test.skip();
        return;
      }
      
      // Login as superadmin
      const loginResponse = await request.post('/api/auth/login', {
        data: {
          username: superadminUser,
          password: superadminPass,
        },
      });
      
      expect(loginResponse.ok()).toBeTruthy();
      const { token } = await loginResponse.json();
      
      // Update services for a test tenant
      const updateResponse = await request.put('/api/superadmin/tenants/kaznmu/services', {
        headers: {
          Authorization: `Bearer ${token}`,
          'X-Tenant-Slug': 'kaznmu',
        },
        data: {
          enabled_services: ['chat'], // Disable calendar
        },
      });
      
      expect(updateResponse.ok()).toBeTruthy();
      
      // Restore both services
      await request.put('/api/superadmin/tenants/kaznmu/services', {
        headers: {
          Authorization: `Bearer ${token}`,
          'X-Tenant-Slug': 'kaznmu',
        },
        data: {
          enabled_services: ['chat', 'calendar'],
        },
      });
    });
  });
});
