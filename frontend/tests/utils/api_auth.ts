import { APIRequestContext, request as playwrightRequest } from '@playwright/test';

export async function loginViaAPI(request: APIRequestContext, email: string, pass: string) {
  // Create a new context for backend API
  const apiContext = await playwrightRequest.newContext({
    baseURL: 'http://localhost:8280'
  });

  const res = await apiContext.post('/api/auth/login', {
    data: { username: email, password: pass }
  });
  if (!res.ok()) {
    throw new Error(`API Login failed: ${await res.text()}`);
  }
  const body = await res.json();
  return body.token;
}

export async function createUserViaAPI(token: string, userData: any) {
  const apiContext = await playwrightRequest.newContext({
    baseURL: 'http://localhost:8280'
  });
  
  const res = await apiContext.post('/api/admin/users', {
    headers: { Authorization: `Bearer ${token}` },
    data: userData
  });
  
  if (!res.ok()) {
    throw new Error(`API Create User failed: ${await res.text()}`);
  }
  return await res.json();
}
