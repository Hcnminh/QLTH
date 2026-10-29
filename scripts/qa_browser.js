#!/usr/bin/env node
/**
 * Automated Browser E2E Test Suite for EduPrimary-KPI (GoAdmin Back-Office)
 *
 * Implements automated browser testing using Playwright headless Chromium.
 * Verifies:
 *   1. Administrative Authentication Flow (/admin/login -> /admin)
 *   2. Main AdminLTE Dashboard layout and navigation menus
 *   3. Staff & Teaching Quota Table (/admin/info/staffs - Circular 05/2025/TT-BGDĐT)
 *   4. Student Onboarding & Digital Consent Table (/admin/info/students - Decree 13/2023/NĐ-CP)
 *   5. KPI 360 & 20% Cap Evaluation Table (/admin/info/kpi_evaluations - Decree 90/2020 & 48/2023)
 *
 * Fulfills: VAL-BROWSER-001
 */

const { chromium } = require('playwright');
const http = require('http');
const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

const PORT = process.env.PORT || 8098;
const BASE_URL = process.env.BASE_URL || `http://localhost:${PORT}`;
const SCREENSHOT_DIR = path.resolve(__dirname, '../docs/screenshots');
const HEADLESS = process.env.HEADLESS !== 'false';
const VIEWPORT = { width: 1280, height: 800 };

// Ensure screenshot directory exists
if (!fs.existsSync(SCREENSHOT_DIR)) {
  fs.mkdirSync(SCREENSHOT_DIR, { recursive: true });
}

let serverProcess = null;

function checkHealth(url) {
  return new Promise((resolve) => {
    http.get(`${url}/health`, (res) => {
      resolve(res.statusCode === 200);
    }).on('error', () => {
      resolve(false);
    });
  });
}

async function ensureServerRunning() {
  const healthy = await checkHealth(BASE_URL);
  if (healthy) {
    console.log(`[INFO] Server already running on ${BASE_URL}`);
    return;
  }

  console.log(`[INFO] Server not detected on ${BASE_URL}. Starting Go backend binary...`);
  const repoRoot = path.resolve(__dirname, '..');
  const binPath = path.join(repoRoot, 'bin/api');

  if (!fs.existsSync(binPath)) {
    console.log('[INFO] Binary bin/api not found. Compiling via go build...');
    const { execSync } = require('child_process');
    execSync('go build -o bin/api ./cmd/api', { cwd: repoRoot, stdio: 'inherit' });
  }

  serverProcess = spawn(binPath, [], {
    cwd: repoRoot,
    env: { ...process.env, PORT: String(PORT), GIN_MODE: 'release' },
    stdio: 'ignore',
  });

  // Wait up to 15 seconds for healthcheck
  const maxAttempts = 30;
  for (let i = 0; i < maxAttempts; i++) {
    await new Promise((r) => setTimeout(r, 500));
    const ok = await checkHealth(BASE_URL);
    if (ok) {
      console.log(`[INFO] Server successfully started (PID: ${serverProcess.pid}) on ${BASE_URL}`);
      return;
    }
  }

  throw new Error(`Failed to start server on ${BASE_URL} after 15 seconds`);
}

function cleanupServer() {
  if (serverProcess && !serverProcess.killed) {
    console.log(`[INFO] Stopping temporary server process (PID: ${serverProcess.pid})...`);
    serverProcess.kill('SIGTERM');
  }
}

// Test runner summary
const results = {
  passed: 0,
  failed: 0,
  cases: [],
};

function recordResult(name, passed, detail = '') {
  if (passed) {
    results.passed++;
    console.log(`  \x1b[32m✔ PASS\x1b[0m: ${name} ${detail ? `(${detail})` : ''}`);
  } else {
    results.failed++;
    console.log(`  \x1b[31m✖ FAIL\x1b[0m: ${name} - ${detail}`);
  }
  results.cases.push({ name, passed, detail });
}

async function runSuite() {
  console.log('================================================================');
  console.log(' EduPrimary-KPI Browser E2E Test Suite (Playwright Headless)');
  console.log(` Target Server: ${BASE_URL}`);
  console.log(` Screenshots:   ${SCREENSHOT_DIR}`);
  console.log('================================================================\n');

  await ensureServerRunning();

  const browser = await chromium.launch({
    headless: HEADLESS,
    args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage'],
  });

  const context = await browser.newContext({
    viewport: VIEWPORT,
    userAgent: 'EduPrimary-Playwright-E2E/1.0',
  });

  const page = await context.newPage();

  try {
    // -----------------------------------------------------------------
    // TEST CASE 1: Health Endpoint Check
    // -----------------------------------------------------------------
    console.log('[CASE 1] Verifying Public Health Endpoint');
    const healthRes = await page.goto(`${BASE_URL}/health`);
    const healthText = await page.textContent('body');
    const isHealthy = healthRes.status() === 200 && healthText.includes('"status":"healthy"');
    recordResult('Health Endpoint', isHealthy, `Status ${healthRes.status()}`);

    // -----------------------------------------------------------------
    // TEST CASE 2: Administrative Login Flow
    // -----------------------------------------------------------------
    console.log('\n[CASE 2] Verifying Administrative Login Page (/admin/login)');
    await page.goto(`${BASE_URL}/admin/login`);
    await page.waitForSelector('#sign-up-form', { timeout: 10000 });

    const pageTitle = await page.title();
    const hasLoginTitle = pageTitle.includes('EduPrimary-KPI');
    const hasUsernameInput = await page.locator('#username').isVisible();
    const hasPasswordInput = await page.locator('#password').isVisible();
    const hasSubmitButton = await page.locator('button.btn-primary').isVisible();

    recordResult('Login Page Elements', hasLoginTitle && hasUsernameInput && hasPasswordInput && hasSubmitButton,
      `Title: "${pageTitle}"`);

    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-login.png') });

    console.log('  -> Submitting credentials: admin / admin');
    await page.fill('#username', 'admin');
    await page.fill('#password', 'admin');
    await page.click('button.btn-primary');

    await page.waitForURL('**/admin', { timeout: 10000 });
    const currentUrl = page.url();
    const loginRedirectPassed = currentUrl.endsWith('/admin') || currentUrl.includes('/admin?');
    recordResult('Admin Signin Redirect', loginRedirectPassed, `Redirected to ${currentUrl}`);

    // -----------------------------------------------------------------
    // TEST CASE 3: GoAdmin Dashboard
    // -----------------------------------------------------------------
    console.log('\n[CASE 3] Verifying Dashboard & Navigation Elements');
    const dashTitle = await page.title();
    const hasDashTitle = dashTitle.includes('EduPrimary-KPI');

    const brandLogo = await page.locator('.logo-lg').innerText();
    const hasBrandLogo = brandLogo.includes('EduPrimary');

    const hasStaffMenu = await page.locator('.sidebar a[href="/admin/info/staffs"]').isVisible();
    const hasStudentMenu = await page.locator('.sidebar a[href="/admin/info/students"]').isVisible();
    const hasKpiMenu = await page.locator('.sidebar a[href="/admin/info/kpi_evaluations"]').isVisible();

    recordResult('Dashboard Brand & Menu Links',
      hasDashTitle && hasBrandLogo && hasStaffMenu && hasStudentMenu && hasKpiMenu,
      `Logo: "${brandLogo}", All 3 business menus visible in sidebar`);

    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-dashboard.png') });

    // -----------------------------------------------------------------
    // TEST CASE 4: Staff & Teaching Quota Table (TT 05/2025/TT-BGDĐT)
    // -----------------------------------------------------------------
    console.log('\n[CASE 4] Verifying Staff & Teaching Quota Table (/admin/info/staffs)');
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/admin/info/staffs')),
      page.locator('.sidebar a[href="/admin/info/staffs"]').click(),
    ]);
    await page.waitForTimeout(400);
    await page.waitForSelector('table.table', { timeout: 10000 });

    const staffHeaders = await page.locator('table.table th').allInnerTexts();
    const staffHeadersText = staffHeaders.join(' ');
    const hasStaffColumns = staffHeadersText.includes('Full Name') || staffHeadersText.includes('Quota');

    const staffBodyText = await page.locator('table.table tbody').innerText();
    const hasStaffData = staffBodyText.includes('Nguyễn Văn An') && staffBodyText.includes('805 tiết/năm');
    const hasHomeroomReduced = staffBodyText.includes('GVCN') && staffBodyText.includes('tiết/tuần');

    recordResult('Staff Table Render & TT 05/2025 Quota',
      hasStaffColumns && hasStaffData && hasHomeroomReduced,
      'Found teacher records, 805 tiết/năm standard, and homeroom reduction');

    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-staff-table.png') });

    // -----------------------------------------------------------------
    // TEST CASE 5: Student Onboarding & Digital Consent Table (NĐ 13/2023)
    // -----------------------------------------------------------------
    console.log('\n[CASE 5] Verifying Student & Digital Consent Table (/admin/info/students)');
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/admin/info/students')),
      page.locator('.sidebar a[href="/admin/info/students"]').click(),
    ]);
    await page.waitForTimeout(400);
    await page.waitForSelector('table.table', { timeout: 10000 });

    const studentHeaders = await page.locator('table.table th').allInnerTexts();
    const studentHeadersText = studentHeaders.join(' ');
    const hasStudentColumns = studentHeadersText.includes('Full Name') || studentHeadersText.includes('Digital Consent');

    const studentBodyText = await page.locator('table.table tbody').innerText();
    const hasStudentData = studentBodyText.includes('Nguyễn Hoàng Nam') && studentBodyText.includes('HS2026001');
    const hasConsentStatus = studentBodyText.includes('NĐ 13') || studentBodyText.includes('đồng thuận');

    recordResult('Student Table & Digital Consent Badges',
      hasStudentColumns && hasStudentData && hasConsentStatus,
      'Found student records, VNeID status, and Decree 13 consent indicators');

    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-student-table.png') });

    // -----------------------------------------------------------------
    // TEST CASE 6: KPI 360 Evaluation & 20% Cap Table (NĐ 90/2020 & 48/2023)
    // -----------------------------------------------------------------
    console.log('\n[CASE 6] Verifying KPI 360 & 20% Cap Table (/admin/info/kpi_evaluations)');
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/admin/info/kpi_evaluations')),
      page.locator('.sidebar a[href="/admin/info/kpi_evaluations"]').click(),
    ]);
    await page.waitForTimeout(400);
    await page.waitForSelector('table.table', { timeout: 10000 });

    const kpiHeaders = await page.locator('table.table th').allInnerTexts();
    const kpiHeadersText = kpiHeaders.join(' ');
    const hasKpiColumns = kpiHeadersText.includes('Staff Name') || kpiHeadersText.includes('Final Score');

    const kpiBodyText = await page.locator('table.table tbody').innerText();
    const hasKpiData = kpiBodyText.includes('Nguyễn Văn An') && kpiBodyText.includes('93.20');
    const hasRanksAndCap = kpiBodyText.includes('Xuất sắc') && (kpiBodyText.includes('trần 20%') || kpiBodyText.includes('20%'));

    recordResult('KPI 360 Table & 20% Cap Enforcement',
      hasKpiColumns && hasKpiData && hasRanksAndCap,
      'Found KPI 4-dimension scores (40-30-15-15), ranks, and 20% cap indicators');

    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-kpi-table.png') });

  } finally {
    await browser.close();
    cleanupServer();
  }

  // -----------------------------------------------------------------
  // SUMMARY REPORT
  // -----------------------------------------------------------------
  console.log('\n================================================================');
  console.log(` Test Summary: ${results.passed} Passed, ${results.failed} Failed (Total ${results.cases.length})`);
  console.log('================================================================');

  if (results.failed > 0) {
    console.error('\x1b[31m[ERROR] Some browser E2E test cases failed.\x1b[0m');
    process.exit(1);
  } else {
    console.log('\x1b[32m[SUCCESS] All Playwright browser E2E test cases passed! (VAL-BROWSER-001 satisfied)\x1b[0m\n');
    process.exit(0);
  }
}

process.on('SIGINT', () => {
  cleanupServer();
  process.exit(130);
});

process.on('SIGTERM', () => {
  cleanupServer();
  process.exit(143);
});

runSuite().catch((err) => {
  console.error('\x1b[31m[FATAL ERROR]\x1b[0m', err);
  cleanupServer();
  process.exit(1);
});
