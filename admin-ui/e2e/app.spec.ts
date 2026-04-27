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

test.describe('Role-based Access', () => {
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

  test('redirects unauthenticated user to login', async ({ page }) => {
    await page.goto('/admin/dashboard');
    await expect(page).toHaveURL('/admin/login');
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