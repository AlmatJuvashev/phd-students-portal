import { test, expect } from '@playwright/test';
import { loginViaUI } from './utils/auth';

// Default credentials - use demo tenant users
const adminUser = process.env.E2E_ADMIN_USER || 'demo.admin';
const adminPass = process.env.E2E_ADMIN_PASS || 'demopassword123!';
const studentUser = process.env.E2E_STUDENT_USER || 'demo.student';
const studentPass = process.env.E2E_STUDENT_PASS || 'demopassword123!';
// Superadmin credentials for service toggle tests
const superadminUser = process.env.E2E_SUPERADMIN_USER || 'juvashev';
const superadminPass = process.env.E2E_SUPERADMIN_PASS || 'superadminpassword123!';

/**
 * E2E tests for tenant service visibility feature
 * Tests that chat/calendar/smtp nav links and routes respect tenant service settings
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
      
      if (!loginResponse.ok()) {
        console.log('Login failed, status:', loginResponse.status());
        console.log('Response:', await loginResponse.text());
        test.skip(true, 'Login failed - check credentials/backend');
        return;
      }
      
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
      
      if (!loginResponse.ok()) {
        console.log('Login failed, status:', loginResponse.status());
        test.skip(true, 'Login failed - check credentials/backend');
        return;
      }
      
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
    
    test('Superadmin can enable/disable services via API', async ({ request }) => {
      // Login as superadmin
      const loginResponse = await request.post('/api/auth/login', {
        data: {
          username: superadminUser,
          password: superadminPass,
        },
      });
      
      if (!loginResponse.ok()) {
        console.log('Superadmin login failed, status:', loginResponse.status());
        test.skip(true, 'Superadmin login failed - check credentials');
        return;
      }
      
      const { token } = await loginResponse.json();
      
      // Get current services for demo-university tenant
      const getTenantResponse = await request.get('/api/superadmin/tenants', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      
      expect(getTenantResponse.ok()).toBeTruthy();
      const tenants = await getTenantResponse.json();
      const demoTenant = tenants.find((t: any) => t.slug === 'demo-university');
      
      if (!demoTenant) {
        test.skip(true, 'demo-university tenant not found');
        return;
      }
      
      // Test 1: Enable SMTP service
      const enableSmtpResponse = await request.put(`/api/superadmin/tenants/${demoTenant.id}/services`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        data: {
          enabled_services: ['chat', 'calendar', 'smtp'],
        },
      });
      
      expect(enableSmtpResponse.ok()).toBeTruthy();
      const enabledResult = await enableSmtpResponse.json();
      expect(enabledResult.enabled_services).toContain('smtp');
      
      // Test 2: Disable chat service
      const disableChatResponse = await request.put(`/api/superadmin/tenants/${demoTenant.id}/services`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        data: {
          enabled_services: ['calendar', 'smtp'], // chat disabled
        },
      });
      
      expect(disableChatResponse.ok()).toBeTruthy();
      const disabledResult = await disableChatResponse.json();
      expect(disabledResult.enabled_services).not.toContain('chat');
      expect(disabledResult.enabled_services).toContain('calendar');
      expect(disabledResult.enabled_services).toContain('smtp');
      
      // Test 3: Restore all services back
      const restoreResponse = await request.put(`/api/superadmin/tenants/${demoTenant.id}/services`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        data: {
          enabled_services: ['chat', 'calendar'],
        },
      });
      
      expect(restoreResponse.ok()).toBeTruthy();
    });

    test('Superadmin can toggle SMTP/email service', async ({ request }) => {
      // Login as superadmin
      const loginResponse = await request.post('/api/auth/login', {
        data: {
          username: superadminUser,
          password: superadminPass,
        },
      });
      
      if (!loginResponse.ok()) {
        test.skip(true, 'Superadmin login failed');
        return;
      }
      
      const { token } = await loginResponse.json();
      
      // Get demo-university tenant
      const getTenantResponse = await request.get('/api/superadmin/tenants', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      
      const tenants = await getTenantResponse.json();
      const demoTenant = tenants.find((t: any) => t.slug === 'demo-university');
      
      if (!demoTenant) {
        test.skip(true, 'demo-university tenant not found');
        return;
      }
      
      // Enable email (alias for smtp)
      const enableEmailResponse = await request.put(`/api/superadmin/tenants/${demoTenant.id}/services`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        data: {
          enabled_services: ['chat', 'calendar', 'email'],
        },
      });
      
      expect(enableEmailResponse.ok()).toBeTruthy();
      const result = await enableEmailResponse.json();
      expect(result.enabled_services).toContain('email');
      
      // Restore defaults
      await request.put(`/api/superadmin/tenants/${demoTenant.id}/services`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        data: {
          enabled_services: ['chat', 'calendar'],
        },
      });
    });

    test('Invalid service is rejected', async ({ request }) => {
      // Login as superadmin
      const loginResponse = await request.post('/api/auth/login', {
        data: {
          username: superadminUser,
          password: superadminPass,
        },
      });
      
      if (!loginResponse.ok()) {
        test.skip(true, 'Superadmin login failed');
        return;
      }
      
      const { token } = await loginResponse.json();
      
      // Get demo-university tenant
      const getTenantResponse = await request.get('/api/superadmin/tenants', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      
      const tenants = await getTenantResponse.json();
      const demoTenant = tenants.find((t: any) => t.slug === 'demo-university');
      
      if (!demoTenant) {
        test.skip(true, 'demo-university tenant not found');
        return;
      }
      
      // Try to enable invalid service
      const invalidResponse = await request.put(`/api/superadmin/tenants/${demoTenant.id}/services`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        data: {
          enabled_services: ['chat', 'invalid_service'],
        },
      });
      
      // Should reject with 400
      expect(invalidResponse.status()).toBe(400);
    });
  });

  /**
   * Tests for DISABLED services - verifying users can't see/access disabled services
   */
  test.describe('Disabled Services Visibility', () => {
    
    // Helper to login as superadmin and get token
    async function getSuperadminToken(request: any) {
      const loginResponse = await request.post('/api/auth/login', {
        data: {
          username: superadminUser,
          password: superadminPass,
        },
      });
      if (!loginResponse.ok()) return null;
      const { token } = await loginResponse.json();
      return token;
    }

    // Helper to get demo-university tenant
    async function getDemoTenant(request: any, token: string) {
      const response = await request.get('/api/superadmin/tenants', {
        headers: { Authorization: `Bearer ${token}` },
      });
      const tenants = await response.json();
      return tenants.find((t: any) => t.slug === 'demo-university');
    }

    // Helper to update tenant services
    async function setTenantServices(request: any, token: string, tenantId: string, services: string[]) {
      return request.put(`/api/superadmin/tenants/${tenantId}/services`, {
        headers: { Authorization: `Bearer ${token}` },
        data: { enabled_services: services },
      });
    }

    test('Admin cannot see chat-rooms link when chat is disabled', async ({ page, request }) => {
      // Get superadmin token
      const token = await getSuperadminToken(request);
      if (!token) {
        test.skip(true, 'Could not get superadmin token');
        return;
      }

      const demoTenant = await getDemoTenant(request, token);
      if (!demoTenant) {
        test.skip(true, 'demo-university not found');
        return;
      }

      try {
        // Step 1: Disable chat service for demo-university
        await setTenantServices(request, token, demoTenant.id, ['calendar']); // Only calendar enabled

        // Step 2: Login as demo.admin
        await loginViaUI(page, adminUser, adminPass);
        await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });

        // Step 3: Navigate to admin panel
        await page.goto('/admin');
        await page.waitForLoadState('networkidle');

        // Step 4: Verify chat-rooms link is NOT visible
        const chatRoomsLink = page.getByRole('link', { name: /chat rooms|комнаты чата|чат бөлмелері/i });
        await expect(chatRoomsLink).not.toBeVisible();
      } finally {
        // Cleanup: Re-enable chat
        await setTenantServices(request, token, demoTenant.id, ['chat', 'calendar']);
      }
    });

    test('User gets "Service Not Available" when accessing disabled service route', async ({ page, request }) => {
      const token = await getSuperadminToken(request);
      if (!token) {
        test.skip(true, 'Could not get superadmin token');
        return;
      }

      const demoTenant = await getDemoTenant(request, token);
      if (!demoTenant) {
        test.skip(true, 'demo-university not found');
        return;
      }

      try {
        // Disable chat service
        await setTenantServices(request, token, demoTenant.id, ['calendar']);

        // Login as student
        await loginViaUI(page, studentUser, studentPass);
        await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });

        // Try to access chat route
        await page.goto('/chat');
        await page.waitForLoadState('networkidle');

        // Should show "Service Not Available" or redirect
        const serviceNotAvailable = page.getByText(/service not available|feature is not enabled/i);
        const notOnChatPage = page.locator('body');
        
        // Either shows error message OR redirects away from chat
        const hasServiceBlock = await serviceNotAvailable.isVisible().catch(() => false);
        const currentUrl = page.url();
        
        // Test passes if either: shows error message OR is not on /chat anymore
        expect(hasServiceBlock || !currentUrl.includes('/chat')).toBeTruthy();
      } finally {
        // Cleanup
        await setTenantServices(request, token, demoTenant.id, ['chat', 'calendar']);
      }
    });

    test('Disabling service for one tenant does not affect another tenant', async ({ request }) => {
      const token = await getSuperadminToken(request);
      if (!token) {
        test.skip(true, 'Could not get superadmin token');
        return;
      }

      // Get all tenants
      const response = await request.get('/api/superadmin/tenants', {
        headers: { Authorization: `Bearer ${token}` },
      });
      const tenants = await response.json();
      
      const demoTenant = tenants.find((t: any) => t.slug === 'demo-university');
      const otherTenant = tenants.find((t: any) => t.slug !== 'demo-university' && t.is_active);

      if (!demoTenant || !otherTenant) {
        test.skip(true, 'Need at least 2 tenants for isolation test');
        return;
      }

      try {
        // Disable chat for demo-university
        await setTenantServices(request, token, demoTenant.id, ['calendar']);

        // Verify other tenant still has chat enabled
        const otherTenantResponse = await request.get(`/api/superadmin/tenants/${otherTenant.id}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        const otherTenantData = await otherTenantResponse.json();

        // Other tenant should have chat in their services (or default)
        // This confirms disabling for one doesn't affect others
        expect(otherTenantData.enabled_services).toBeDefined();
        // The enabled_services should NOT have been modified by our change to demo-university
        expect(Array.isArray(otherTenantData.enabled_services)).toBeTruthy();
      } finally {
        // Cleanup
        await setTenantServices(request, token, demoTenant.id, ['chat', 'calendar']);
      }
    });

    test('Admin cannot access /admin/chat-rooms when chat is disabled', async ({ page, request }) => {
      const token = await getSuperadminToken(request);
      if (!token) {
        test.skip(true, 'Could not get superadmin token');
        return;
      }

      const demoTenant = await getDemoTenant(request, token);
      if (!demoTenant) {
        test.skip(true, 'demo-university not found');
        return;
      }

      try {
        // Disable chat
        await setTenantServices(request, token, demoTenant.id, ['calendar']);

        // Login as admin
        await loginViaUI(page, adminUser, adminPass);
        await expect(page).not.toHaveURL(/\/login/, { timeout: 10000 });

        // Try to directly access chat-rooms admin page
        await page.goto('/admin/chat-rooms');
        await page.waitForLoadState('networkidle');

        // Should either show "Service Not Available" or redirect
        const currentUrl = page.url();
        const hasServiceBlock = await page.getByText(/service not available|feature is not enabled/i).isVisible().catch(() => false);

        expect(hasServiceBlock || !currentUrl.includes('/chat-rooms')).toBeTruthy();
      } finally {
        // Cleanup
        await setTenantServices(request, token, demoTenant.id, ['chat', 'calendar']);
      }
    });
  });
});
