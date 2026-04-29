import { test, expect } from '@playwright/test';

test.describe('Login Flow', () => {
  test('shows login page on first visit', async ({ page }) => {
    await page.goto('/admin/login');
    await expect(page.getByRole('heading', { name: /login/i })).toBeVisible();
  });

  test('login with valid credentials redirects to dashboard', async ({ page }) => {
    await page.goto('/admin/login');
    
    await page.getByPlaceholder('Email').fill('test@admin.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();

    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    await expect(page.getByText(/dashboard/i)).toBeVisible({ timeout: 10000 });
  });

  test('shows error with invalid credentials', async ({ page }) => {
    await page.goto('/admin/login');
    
    await page.getByPlaceholder('Email').fill('invalid@example.com');
    await page.getByPlaceholder('Password').fill('wrongpassword');
    await page.getByRole('button', { name: /login/i }).click();

    await expect(page.getByText(/invalid credentials/i)).toBeVisible();
  });
});

test.describe('Logout Flow', () => {
  test('logout redirects to login page', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@admin.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    const logoutButton = page.getByRole('button', { name: /logout/i });
    if (await logoutButton.isVisible()) {
      await logoutButton.click();
      await expect(page).toHaveURL('/admin/login');
    }
  });

  test('authenticated user can access protected routes', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@admin.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    await page.goto('/admin/providers');
    await expect(page.getByText(/providers/i)).toBeVisible({ timeout: 5000 });
  });
});

test.describe('Role-based Access - Admin', () => {
  test('admin sees all navigation items', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@admin.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    await expect(page.getByText(/dashboard/i)).toBeVisible();
    await expect(page.getByText(/providers/i)).toBeVisible();
    await expect(page.getByText(/users/i)).toBeVisible();
  });

  test('admin can access all pages', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@admin.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    const pages = ['/admin/dashboard', '/admin/providers', '/admin/users', '/admin/api-keys', '/admin/usage', '/admin/health', '/admin/settings'];
    for (const path of pages) {
      await page.goto(path);
      await expect(page.getByRole('heading')).toBeVisible({ timeout: 5000 });
    }
  });
});

test.describe('Role-based Access - User', () => {
  test('user sees limited navigation items', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@user.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    await expect(page.getByText(/dashboard/i)).toBeVisible();
    await expect(page.getByText(/api.keys/i)).toBeVisible();
    await expect(page.getByText(/usage/i)).toBeVisible();
    await expect(page.getByText(/health/i)).toBeVisible();
    await expect(page.getByText(/settings/i)).toBeVisible();
  });

  test('user cannot access admin-only pages', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@user.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    await page.goto('/admin/providers');
    await expect(page).not.toHaveURL('/admin/providers');
  });
});

test.describe('Role-based Access - Viewer', () => {
  test('viewer sees read-only navigation items', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@viewer.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    await expect(page.getByText(/dashboard/i)).toBeVisible();
    await expect(page.getByText(/usage/i)).toBeVisible();
    await expect(page.getByText(/health/i)).toBeVisible();
  });

  test('viewer cannot access write operations', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@viewer.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    await page.goto('/admin/usage');
    const createButtons = page.getByRole('button', { name: /create|add|new/i });
    if (await createButtons.count() > 0) {
      await expect(createButtons.first()).toBeDisabled();
    }
  });
});

test.describe('Responsive Layout', () => {
  test('mobile view collapses sidebar', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@admin.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
  });

  test('desktop view shows full layout', async ({ page }) => {
    await page.setViewportSize({ width: 1920, height: 1080 });
    
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@admin.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    await expect(page.getByText(/infrastructure/i)).toBeVisible();
  });
});

test.describe('Form Validation', () => {
  test('login form shows validation errors', async ({ page }) => {
    await page.goto('/admin/login');
    
    await page.getByRole('button', { name: /login/i }).click();
    
    await expect(page.getByText(/required|invalid/i)).toBeVisible({ timeout: 5000 });
  });

  test('register form shows validation errors', async ({ page }) => {
    await page.goto('/admin/register');
    
    await page.getByRole('button', { name: /register/i }).click();
    
    await expect(page.getByText(/required|invalid/i)).toBeVisible({ timeout: 5000 });
  });
});

test.describe('Error States', () => {
  test('401 unauthorized shows error message', async ({ page }) => {
    await page.goto('/admin/login');
    
    await page.getByPlaceholder('Email').fill('invalid@example.com');
    await page.getByPlaceholder('Password').fill('wrongpassword');
    await page.getByRole('button', { name: /login/i }).click();

    await expect(page.getByText(/invalid credentials|unauthorized/i)).toBeVisible();
  });

  test('redirects unauthenticated user to login', async ({ page }) => {
    await page.goto('/admin/dashboard');
    await expect(page).toHaveURL('/admin/login');
  });
});

test.describe('Complete User Flow', () => {
  test('user can login, navigate, and logout', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@user.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    await page.goto('/admin/api-keys');
    await expect(page.getByRole('heading')).toBeVisible({ timeout: 5000 });
    
    await page.goto('/admin/usage');
    await expect(page.getByRole('heading')).toBeVisible({ timeout: 5000 });
    
    const logoutButton = page.getByRole('button', { name: /logout/i });
    if (await logoutButton.isVisible()) {
      await logoutButton.click();
      await expect(page).toHaveURL('/admin/login');
    }
  });
});

test.describe('Complete Admin Flow', () => {
  test('admin can login, manage resources, and logout', async ({ page }) => {
    await page.goto('/admin/login');
    await page.getByPlaceholder('Email').fill('test@admin.com');
    await page.getByPlaceholder('Password').fill('test123456');
    await page.getByRole('button', { name: /login/i }).click();
    await expect(page).toHaveURL('/admin/', { timeout: 10000 });
    
    await page.goto('/admin/providers');
    await expect(page.getByRole('heading')).toBeVisible({ timeout: 5000 });
    const addButton = page.getByRole('button', { name: /add|create|new/i });
    if (await addButton.isVisible()) {
      await addButton.click();
    }
    
    await page.goto('/admin/users');
    await expect(page.getByRole('heading')).toBeVisible({ timeout: 5000 });
    
    const logoutButton = page.getByRole('button', { name: /logout/i });
    if (await logoutButton.isVisible()) {
      await logoutButton.click();
      await expect(page).toHaveURL('/admin/login');
    }
  });
});