#!/usr/bin/env node
/**
 * Automated Comprehensive Browser E2E Dogfooding Suite for EduPrimary-KPI
 *
 * Implements Playwright headless browser testing and multi-role dogfooding across:
 *   1. Complete Set of 10 Subsystems:
 *      - Phân hệ 1: Quản lý Hồ sơ Học sinh & Consent NĐ 13/2023
 *      - Phân hệ 2: Đánh giá Học sinh TT 27/2020 & TT 89/2026 (Strategy Versioning 01/12/2026)
 *      - Phân hệ 3: Kế hoạch Bài dạy CV 2345 / STEM TT 86 & Định mức Tiết dạy TT 05/2025
 *      - Phân hệ 4: Bán trú Dinh dưỡng (Kcal Viện Dinh Dưỡng, Lưu mẫu 24h QĐ 1246/QĐ-BYT, Hoàn tiền ăn)
 *      - Phân hệ 5: Y tế Học đường & Sức khỏe (BMI, Thị lực, Sổ y tế, Đơn thuốc phụ huynh)
 *      - Phân hệ 6: Quản lý Lớp học & Danh mục Môn học (CT GDPT 2018)
 *      - Phân hệ 7: Tài chính, VietQR Động EMVCo & Hóa đơn Điện tử NĐ 123/2020/NĐ-CP
 *      - Phân hệ 8: Đánh giá KPI 360 độ (NĐ 90/2020, NĐ 48/2023, Điểm liệt TT 29/2024, Trần 20%)
 *      - Phân hệ 9: Học bạ số Bất biến & Chữ ký số SmartCA (CMS/PKCS#7, TSA Timestamp)
 *      - Phân hệ 10: Cổng Tích hợp Ngoại vi Sandbox Stubs (VNeID Level 2, MOET CSDL Outbox)
 *
 *   2. Multi-Role Workflows across all 7 Personas:
 *      - BGH (Ban Giám Hiệu / Hiệu Trưởng)
 *      - GVCN (Giáo Viên Chủ Nhiệm)
 *      - GVBM (Giáo Viên Bộ Môn)
 *      - Kế toán (KE_TOAN)
 *      - Y tế (Y_TE)
 *      - Văn thư (VAN_THU)
 *      - Phụ huynh (PHHS)
 *
 *   3. All 11 GoAdmin UI Table Views & Detail Panels Render Verification
 *   4. Collection of Annotated Step-by-Step Screenshots in docs/screenshots/
 *
 * Fulfills: VAL-BROWSER-003, VAL-BROWSER-004
 */

const { chromium, request } = require('playwright');
const http = require('http');
const { spawn, execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const PORT = process.env.PORT || 8098;
const BASE_URL = process.env.BASE_URL || `http://localhost:${PORT}`;
const SCREENSHOT_DIR = path.resolve(__dirname, '../docs/screenshots');
const HEADLESS = process.env.HEADLESS !== 'false';
const VIEWPORT = { width: 1366, height: 860 };
const SCHOOL_ID = '00000000-0000-0000-0000-000000000001';

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
    execSync('go build -o bin/api ./cmd/api', { cwd: repoRoot, stdio: 'inherit' });
  }

  serverProcess = spawn(binPath, [], {
    cwd: repoRoot,
    env: { ...process.env, PORT: String(PORT), GIN_MODE: 'release' },
    stdio: 'ignore',
  });

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
  console.log(' EduPrimary-KPI Full Browser Dogfooding E2E Suite (Playwright)');
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
    userAgent: 'EduPrimary-Playwright-E2E/2.0',
  });

  const page = await context.newPage();

  // Create API request context for multi-role workflow simulation
  const apiContext = await request.newContext({
    baseURL: BASE_URL,
    extraHTTPHeaders: {
      'Content-Type': 'application/json',
      'X-School-ID': SCHOOL_ID,
    },
  });

  try {
    // =================================================================
    // SECTION 1: PUBLIC HEALTH & READINESS
    // =================================================================
    console.log('\n--- SECTION 1: Public Health & Service Discovery ---');
    const healthRes = await page.goto(`${BASE_URL}/health`);
    const healthText = await page.textContent('body');
    const isHealthy = healthRes.status() === 200 && healthText.includes('"service":"eduprimary-kpi-core"');
    recordResult('Health Endpoint Check', isHealthy, `Status ${healthRes.status()}`);

    // =================================================================
    // SECTION 2: ADMINISTRATIVE AUTHENTICATION & LOGIN FLOW
    // =================================================================
    console.log('\n--- SECTION 2: Administrative Authentication Flow (/admin/login) ---');
    await page.goto(`${BASE_URL}/admin/login`);
    await page.waitForSelector('#sign-up-form', { timeout: 10000 });

    const pageTitle = await page.title();
    const hasLoginTitle = pageTitle.includes('EduPrimary-KPI');
    const hasUsernameInput = await page.locator('#username').isVisible();
    const hasPasswordInput = await page.locator('#password').isVisible();
    const hasSubmitButton = await page.locator('button.btn-primary').isVisible();

    recordResult('Admin Login Screen Elements',
      hasLoginTitle && hasUsernameInput && hasPasswordInput && hasSubmitButton,
      `Form controls visible, Title: "${pageTitle}"`);

    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '01_login_page.png') });
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-login.png') });

    console.log('  -> Submitting credentials: admin / admin');
    await page.fill('#username', 'admin');
    await page.fill('#password', 'admin');
    await page.click('button.btn-primary');

    await page.waitForURL('**/admin', { timeout: 10000 });
    const currentUrl = page.url();
    const loginRedirectPassed = currentUrl.endsWith('/admin') || currentUrl.includes('/admin?');
    recordResult('Admin Signin Redirect', loginRedirectPassed, `Redirected to ${currentUrl}`);

    // =================================================================
    // SECTION 3: CENTRAL DASHBOARD & SIDEBAR NAVIGATION
    // =================================================================
    console.log('\n--- SECTION 3: GoAdmin Central Dashboard & Subsystem Widgets ---');
    await page.waitForSelector('.logo-lg', { timeout: 8000 });
    const dashTitle = await page.title();
    const hasDashTitle = dashTitle.includes('EduPrimary-KPI');
    const brandLogo = await page.locator('.logo-lg').innerText();
    const hasBrandLogo = brandLogo.includes('EduPrimary');

    // Verify all 11 Subsystem Sidebar Menus exist
    const expectedSidebarUris = [
      '/admin/info/staffs',
      '/admin/info/classes',
      '/admin/info/students',
      '/admin/info/lesson_plans',
      '/admin/info/subjects',
      '/admin/info/student_subject_evaluations',
      '/admin/info/student_behavior_evaluations',
      '/admin/info/academic_transcripts_snapshots',
      '/admin/info/boarding_menus',
      '/admin/info/student_health_records',
      '/admin/info/kpi_evaluations',
    ];

    let allMenusPresent = true;
    for (const uri of expectedSidebarUris) {
      const isVisible = await page.locator(`.sidebar a[href="${uri}"]`).isVisible();
      if (!isVisible) {
        allMenusPresent = false;
        console.log(`    [WARN] Sidebar link missing for: ${uri}`);
      }
    }

    recordResult('Dashboard Brand & 11 Subsystem Menus',
      hasDashTitle && hasBrandLogo && allMenusPresent,
      `Brand: "${brandLogo}", All 11 subsystem sidebar links visible`);

    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '02_dashboard.png') });
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '02_dashboard_overview.png') });
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-dashboard.png') });

    // =================================================================
    // SECTION 4: ALL 11 GOADMIN UI TABLE VIEWS & DETAIL PANELS
    // =================================================================
    console.log('\n--- SECTION 4: Exercising All 11 GoAdmin Subsystem UI Tables ---');

    async function verifyTable({
      pathUri,
      tableName,
      expectedColumns,
      expectedText,
      screenshotName,
      detailScreenshotName,
    }) {
      console.log(`\n  [TABLE] Navigating to ${tableName} (${pathUri})`);
      const res = await page.goto(`${BASE_URL}${pathUri}`);
      await page.waitForSelector('table.table', { timeout: 10000 });

      const headers = (await page.locator('table.table th').allInnerTexts()).join(' ');
      const bodyText = await page.locator('table.table tbody').innerText();

      const columnsMatched = expectedColumns.every((col) => headers.includes(col));
      const textMatched = expectedText.every((txt) => bodyText.includes(txt));

      recordResult(`${tableName} UI Table Render`,
        res.status() === 200 && columnsMatched && textMatched,
        `Status 200, Columns: [${expectedColumns.join(', ')}] matched`);

      if (screenshotName) {
        await page.screenshot({ path: path.join(SCREENSHOT_DIR, screenshotName) });
      }

      // Check detail view if available
      if (detailScreenshotName) {
        const detailBtn = page.locator('.grid-row-view').first();
        if (await detailBtn.isVisible().catch(() => false)) {
          await detailBtn.click();
          await page.waitForSelector('.box', { timeout: 8000 });
          await page.screenshot({ path: path.join(SCREENSHOT_DIR, detailScreenshotName) });
          // Go back to the list table
          await page.goto(`${BASE_URL}${pathUri}`);
          await page.waitForSelector('table.table', { timeout: 10000 });
        }
      }
    }

    // 1. Staffs & Teaching Quotas Table (TT 05/2025/TT-BGDĐT)
    await verifyTable({
      pathUri: '/admin/info/staffs',
      tableName: 'Phân hệ 3 - Giáo viên & Định mức Tiết dạy',
      expectedColumns: ['ID', 'Full Name', 'Homeroom', 'Quota'],
      expectedText: ['Nguyễn Văn An', '805 tiết/năm', 'GVCN'],
      screenshotName: '03_staff_table.png',
      detailScreenshotName: '18_staff_detail.png',
    });
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '03_staffs_teaching_quota.png') });
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-staff-table.png') });

    // 2. Classes Management Table
    await verifyTable({
      pathUri: '/admin/info/classes',
      tableName: 'Phân hệ 6 - Quản lý Danh sách Lớp học',
      expectedColumns: ['Tên Lớp', 'Khối Lớp', 'Phòng học', 'Sĩ số'],
      expectedText: ['1A1', '5A1', 'Phòng'],
      screenshotName: '06_classes_management.png',
    });

    // 3. Students & Digital Consent Table (Nghị định 13/2023/NĐ-CP Điều 20)
    await verifyTable({
      pathUri: '/admin/info/students',
      tableName: 'Phân hệ 1 - Hồ sơ Học sinh & Digital Consent NĐ 13',
      expectedColumns: ['ID', 'Full Name', 'Student Code', 'Digital Consent'],
      expectedText: ['Nguyễn Hoàng Nam', 'HS2026001', 'NĐ 13'],
      screenshotName: '04_student_table.png',
      detailScreenshotName: '19_student_detail.png',
    });
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '05_students_digital_consent.png') });
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-student-table.png') });

    // 4. Lesson Plans Table (Công văn 2345/BGDĐT-GDTH & STEM TT 86)
    await verifyTable({
      pathUri: '/admin/info/lesson_plans',
      tableName: 'Phân hệ 3 - Kế hoạch Bài dạy CV 2345 / STEM',
      expectedColumns: ['Tên Bài dạy', 'Môn học', 'Khối', 'Tuần', 'Chuẩn Khung', 'Trạng thái Duyệt'],
      expectedText: ['CV 2345', 'Khối', 'tiết'],
      screenshotName: '07_lesson_plans_cv2345.png',
      detailScreenshotName: '14_lesson_plan_detail.png',
    });

    // 5. Subjects Catalog Table (CT GDPT 2018)
    await verifyTable({
      pathUri: '/admin/info/subjects',
      tableName: 'Phân hệ 2 - Danh mục Môn học & Hình thức Đánh giá',
      expectedColumns: ['Mã Môn', 'Tên Môn học', 'Hình thức Đánh giá'],
      expectedText: ['Toán', 'Tiếng Việt', 'Tin học'],
      screenshotName: '08_subjects_catalog.png',
    });

    // 6. Student Subject Evaluations Table (TT 27/2020 & TT 89/2026)
    await verifyTable({
      pathUri: '/admin/info/student_subject_evaluations',
      tableName: 'Phân hệ 2 - Sổ Đánh giá Môn học Học sinh',
      expectedColumns: ['Mã Học sinh', 'Môn học', 'Học kỳ', 'Mức Đạt', 'Điểm KT'],
      expectedText: ['HTT', '10.0'],
      screenshotName: '09_subject_evaluations_tt27_89.png',
    });

    // 7. Student Behavior Evaluations Table (5 Phẩm chất, 10 Năng lực T/D/C)
    await verifyTable({
      pathUri: '/admin/info/student_behavior_evaluations',
      tableName: 'Phân hệ 2 - Đánh giá Năng lực & Phẩm chất',
      expectedColumns: ['Mã Học sinh', 'Học kỳ', 'Yêu nước', 'Nhân ái', 'Chăm chỉ', 'Trung thực', 'Trách nhiệm'],
      expectedText: ['T (Tốt)'],
      screenshotName: '10_behavior_competency_evaluations.png',
    });

    // 8. Academic Transcripts Snapshots Table (Học bạ số & Chữ ký số SmartCA TSA)
    await verifyTable({
      pathUri: '/admin/info/academic_transcripts_snapshots',
      tableName: 'Phân hệ 9 - Snapshot Học bạ số Bất biến & SmartCA',
      expectedColumns: ['Mã Học sinh', 'Năm học', 'Khối Lớp', 'Thông tư Áp dụng', 'Kết quả Lên lớp'],
      expectedText: ['Thông tư', 'Được lên lớp'],
      screenshotName: '11_transcripts_smartca_snapshot.png',
      detailScreenshotName: '17_transcript_snapshot_detail.png',
    });

    // 9. Boarding Menus Table (Kcal Viện Dinh Dưỡng & QĐ 1246/QĐ-BYT Lưu mẫu 24h)
    await verifyTable({
      pathUri: '/admin/info/boarding_menus',
      tableName: 'Phân hệ 4 - Thực đơn Dinh dưỡng Bán trú & Lưu mẫu 24h',
      expectedColumns: ['Ngày thực đơn', 'Bữa ăn', 'Năng lượng (Kcal)', 'Lưu mẫu 24h'],
      expectedText: ['Bữa trưa', 'Lưu mẫu 24h'],
      screenshotName: '12_boarding_nutrition_menu.png',
      detailScreenshotName: '15_boarding_menu_detail.png',
    });

    // 10. Student Health Records Table (BMI, Thị lực & Nha học đường)
    await verifyTable({
      pathUri: '/admin/info/student_health_records',
      tableName: 'Phân hệ 5 - Hồ sơ Sức khỏe Học sinh & Chỉ số BMI',
      expectedColumns: ['Student ID', 'Ngày khám', 'Chiều cao', 'Cân nặng', 'BMI', 'Phân loại thể trạng'],
      expectedText: ['Bình thường'],
      screenshotName: '13_student_health_bmi.png',
      detailScreenshotName: '16_student_health_detail.png',
    });

    // 11. KPI 360 Evaluations Table (Nghị định 90/2020 & 48/2023, Khống chế trần 20%)
    await verifyTable({
      pathUri: '/admin/info/kpi_evaluations',
      tableName: 'Phân hệ 8 - Đánh giá KPI 360 & Khống chế Trần 20%',
      expectedColumns: ['Staff Name', 'Final Score'],
      expectedText: ['Nguyễn Văn An', '93.20', 'Xuất sắc', '20%'],
      screenshotName: '05_kpi_table.png',
      detailScreenshotName: '20_kpi_evaluation_detail.png',
    });
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '13b_kpi_evaluations_20pct_cap.png') });
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'val-browser-001-kpi-table.png') });

    // =================================================================
    // SECTION 5: MULTI-ROLE WORKFLOWS ACROSS ALL 7 PERSONAS
    // =================================================================
    console.log('\n--- SECTION 5: Multi-Role Workflow Dogfooding Across 7 Personas ---');

    // -----------------------------------------------------------------
    // ROLE 1: BAN GIÁM HIỆU (BGH)
    // -----------------------------------------------------------------
    console.log('\n[ROLE 1: BGH - Ban Giám Hiệu Workflows]');

    // 1.1 Quota calculation verification (TT 05/2025: 805 chuẩn, GVCN giảm 4*35=140, net 665, thừa 55)
    const quotaRes = await apiContext.post('/api/v1/teaching/calculate-quota', {
      headers: { 'X-User-Role': 'BGH' },
      data: {
        teacher_id: 'GV001',
        teacher_name: 'Lê Thị Mai',
        school_type: 'PRIMARY',
        is_homeroom: true,
        weeks_in_school_year: 35,
        actual_taught_hours: 720,
      },
    });
    const quotaData = await quotaRes.json();
    const quotaOk = quotaRes.status() === 200
      && quotaData.data.standard_annual_hours === 805
      && quotaData.data.net_annual_quota === 665
      && quotaData.data.overtime_hours === 55;
    recordResult('BGH: Quota Calculation (TT 05/2025 805h, GVCN -4h/w)',
      quotaOk,
      `Standard 805h, Net: ${quotaData.data?.net_annual_quota}h, Overtime: ${quotaData.data?.overtime_hours}h`);

    // 1.2 Timetable slot configuration & conflict enforcement (TT 05/2025)
    const timetableRes = await apiContext.post('/api/v1/timetable/slots', {
      headers: { 'X-User-Role': 'BGH' },
      data: {
        academic_year_id: '2026-2027',
        term_id: '2026-2027-HK1',
        class_id: 'f0000000-0000-0000-0000-000000000002',
        class_name: '5A1',
        teacher_id: 'b0000000-0000-0000-0000-000000000001',
        teacher_name: 'Nguyễn Văn An',
        subject_id: 'TOAN',
        day_of_week: 4,
        period: 1,
        room_number: 'Phòng 5A1',
        week_number: 0,
      },
    });
    recordResult('BGH: Timetable Slot Configuration',
      timetableRes.status() === 201 || timetableRes.status() === 200,
      `Status ${timetableRes.status()}`);

    const timetableConflictRes = await apiContext.post('/api/v1/timetable/slots', {
      headers: { 'X-User-Role': 'BGH' },
      data: {
        academic_year_id: '2026-2027',
        term_id: '2026-2027-HK1',
        class_id: 'f0000000-0000-0000-0000-000000000001',
        class_name: '1A1',
        teacher_id: 'b0000000-0000-0000-0000-000000000001',
        teacher_name: 'Nguyễn Văn An',
        subject_id: 'TOAN',
        day_of_week: 4,
        period: 1,
        room_number: 'Phòng 1A1',
        week_number: 0,
      },
    });
    const timetableConflictText = await timetableConflictRes.text();
    recordResult('BGH: Timetable Teacher Collision Detection (TT 05/2025 Conflict)',
      timetableConflictRes.status() === 409 && timetableConflictText.includes('Thông tư 05/2025/TT-BGDĐT'),
      'HTTP 409 Conflict with legal citation');

    // 1.3 Academic Transcript Snapshot Strategy Resolution (TT27 vs TT89 on 01/12/2026)
    const snapTT27Res = await apiContext.post('/api/v1/transcripts/generate', {
      headers: { 'X-User-Role': 'BGH' },
      data: {
        student_id: 'c0000000-0000-0000-0000-000000000001',
        academic_year_id: '2026-2027',
        grade_level: 5,
        evaluation_date: '2026-10-25',
      },
    });
    const snapTT27Data = await snapTT27Res.json();
    recordResult('BGH: Transcript Strategy Resolution before 01/12/2026 -> TT27_2020',
      snapTT27Res.status() === 200 && snapTT27Data.data?.regulation_applied === 'TT27_2020',
      `Applied: ${snapTT27Data.data?.regulation_applied}`);

    const snapTT89Res = await apiContext.post('/api/v1/transcripts/generate', {
      headers: { 'X-User-Role': 'BGH' },
      data: {
        student_id: 'c0000000-0000-0000-0000-000000000001',
        academic_year_id: '2026-2027',
        grade_level: 5,
        evaluation_date: '2026-12-05',
      },
    });
    const snapTT89Data = await snapTT89Res.json();
    recordResult('BGH: Transcript Strategy Resolution from 01/12/2026 -> TT89_2026',
      snapTT89Res.status() === 200 && snapTT89Data.data?.regulation_applied === 'TT89_2026',
      `Applied: ${snapTT89Data.data?.regulation_applied}`);

    // 1.4 MOET CSDL Outbox Sync
    const moetRes = await apiContext.post('/api/v1/integrations/moet/sync', {
      headers: { 'X-User-Role': 'BGH' },
      data: {
        school_id: SCHOOL_ID,
        academic_year: '2026-2027',
        sync_type: 'ALL',
        batch_size: 20,
      },
    });
    const moetData = await moetRes.json();
    recordResult('BGH: MOET CSDL Outbox Sync Queue Execution',
      moetRes.status() === 200 && moetData.data?.status === 'COMPLETED',
      `Transaction ID: ${moetData.data?.moet_transaction_id}`);

    // 1.5 Sandbox Stubs Health Overview
    const sandboxRes = await apiContext.get('/api/v1/integrations/sandbox/status', {
      headers: { 'X-User-Role': 'BGH' },
    });
    const sandboxData = await sandboxRes.json();
    recordResult('BGH: External Integrations Sandbox Status Check',
      sandboxRes.status() === 200 && sandboxData.data?.environment === 'SANDBOX_STUB' && sandboxData.data?.integrations?.length >= 5,
      `Environment: ${sandboxData.data?.environment}, Active modules: ${sandboxData.data?.integrations?.length}`);

    // -----------------------------------------------------------------
    // ROLE 2: GIÁO VIÊN CHỦ NHIỆM (GVCN)
    // -----------------------------------------------------------------
    console.log('\n[ROLE 2: GVCN - Giáo Viên Chủ Nhiệm Workflows]');

    // 2.1 Student Behavior Evaluation (5 qualities, 10 competencies)
    const behRes = await apiContext.post('/api/v1/evaluations/behaviors', {
      headers: { 'X-User-Role': 'GVCN' },
      data: {
        student_id: 'c0000000-0000-0000-0000-000000000001',
        term_id: '2026-2027-HK1',
        patriotism: 'T',
        compassion: 'T',
        diligence: 'T',
        honesty: 'T',
        responsibility: 'T',
        qualities_comment: 'Học sinh gương mẫu, chấp hành tốt nội quy',
        autonomy_and_self_learning: 'T',
        communication_and_collaboration: 'T',
        problem_solving_and_creativity: 'T',
        language_competence: 'T',
        mathematical_competence: 'T',
        scientific_competence: 'T',
        technological_competence: 'T',
        computer_competence: 'T',
        aesthetic_competence: 'T',
        physical_competence: 'T',
      },
    });
    recordResult('GVCN: Student Behavior & 10 Competencies Evaluation (T/D/C)',
      behRes.status() === 200,
      `Status ${behRes.status()}`);

    // 2.2 AI Copilot Personalized Feedback Generation
    const aiCopilotRes = await apiContext.post('/api/v1/evaluations/ai-copilot', {
      headers: { 'X-User-Role': 'GVCN' },
      data: {
        subject_name: 'Toán',
        achievement_level: 'HTT',
        exam_score: 10.0,
        diligence_level: 'T',
        student_name: 'Minh Khôi',
      },
    });
    const aiCopilotData = await aiCopilotRes.json();
    const hasStudentName = (aiCopilotData.suggestions || []).some((s) => s.includes('Minh Khôi'));
    recordResult('GVCN: AI Copilot Evaluation Assistant Feedback',
      aiCopilotRes.status() === 200 && hasStudentName,
      'Generated personalized remark containing student name');

    // 2.3 Student Health Profile Inspection
    const healthProfileRes = await apiContext.get('/api/v1/students/c0000000-0000-0000-0000-000000000001/health', {
      headers: { 'X-User-Role': 'GVCN' },
    });
    const healthProfileData = await healthProfileRes.json();
    recordResult('GVCN: Review Student Health Profile & BMI Classification',
      healthProfileRes.status() === 200 && healthProfileData.data?.bmi > 0,
      `BMI: ${healthProfileData.data?.bmi}, Classification: ${healthProfileData.data?.bmi_category}`);

    // 2.4 SmartCA Digital Signature Execution
    const smartcaRes = await apiContext.post('/api/v1/integrations/smartca/sign', {
      headers: { 'X-User-Role': 'GVCN' },
      data: {
        document_id: 'd1000000-0000-0000-0000-000000000001',
        document_type: 'TRANSCRIPT_SNAPSHOT',
        signer_id: 'b0000000-0000-0000-0000-000000000001',
        signer_name: 'Trần Thị Bình',
        signer_role: 'GVCN',
      },
    });
    const smartcaData = await smartcaRes.json();
    recordResult('GVCN: SmartCA Digital Signing with TSA Timestamp',
      smartcaRes.status() === 200 && smartcaData.data?.signature_format === 'CMS/PKCS#7' && smartcaData.data?.tsa_token?.tsa_provider,
      `Signature Format: ${smartcaData.data?.signature_format}, TSA: ${smartcaData.data?.tsa_token?.tsa_provider}`);

    // -----------------------------------------------------------------
    // ROLE 3: GIÁO VIÊN BỘ MÔN (GVBM)
    // -----------------------------------------------------------------
    console.log('\n[ROLE 3: GVBM - Giáo Viên Bộ Môn Workflows]');

    // 3.1 Submit CV 2345 Lesson Plan
    const lessonPlanRes = await apiContext.post('/api/v1/lesson-plans/submit', {
      headers: { 'X-User-Role': 'GVBM' },
      data: {
        teacher_id: 'b0000000-0000-0000-0000-000000000001',
        teacher_name: 'Nguyễn Văn An',
        grade_level: 5,
        subject_id: 'TOAN',
        title: 'Bài 15: Phép cộng phân số thập phân (Chuẩn CV 2345)',
        week_number: 11,
        duration_slots: 2,
        curriculum_standard: 'CV_2345',
        objectives: {
          knowledge: 'Nắm vững quy tắc cộng phân số thập phân',
          competencies: 'Năng lực tư duy và lập luận toán học',
          qualities: 'Chăm chỉ, trách nhiệm',
        },
        equipment: 'Thước kẻ, bảng phụ, máy chiếu',
      },
    });
    const lessonPlanData = await lessonPlanRes.json();
    recordResult('GVBM: Submit Lesson Plan CV 2345 (4 Phases & Objectives)',
      lessonPlanRes.status() === 201 && lessonPlanData.data?.current_status === 'SUBMITTED',
      `Status: ${lessonPlanData.data?.current_status}, ID: ${lessonPlanData.data?.id}`);

    // 3.2 Submit Subject Evaluation Score
    const subjectEvalRes = await apiContext.post('/api/v1/evaluations/subject-scores', {
      headers: { 'X-User-Role': 'GVBM' },
      data: {
        student_id: 'c0000000-0000-0000-0000-000000000001',
        term_id: '2026-2027-HK1',
        subject_id: 'TOAN',
        achievement_level: 'HTT',
        exam_score: 10.0,
        teacher_comment: 'Hoàn thành xuất sắc các nội dung môn Toán',
      },
    });
    recordResult('GVBM: Submit Subject Score & Achievement Level (HTT/10.0)',
      subjectEvalRes.status() === 200,
      `Status ${subjectEvalRes.status()}`);

    // 3.3 Negative validation: Reject invalid achievement level
    const invalidEvalRes = await apiContext.post('/api/v1/evaluations/subject-scores', {
      headers: { 'X-User-Role': 'GVBM' },
      data: {
        student_id: 'c0000000-0000-0000-0000-000000000001',
        term_id: '2026-2027-HK1',
        subject_id: 'TOAN',
        achievement_level: 'GIOI', // invalid under TT 27/89
      },
    });
    recordResult('GVBM: Reject Invalid Achievement Level Under TT 27/TT 89',
      invalidEvalRes.status() === 400,
      'HTTP 400 Bad Request');

    // -----------------------------------------------------------------
    // ROLE 4: KẾ TOÁN (KE_TOAN)
    // -----------------------------------------------------------------
    console.log('\n[ROLE 4: Kế toán - School Accountant Workflows]');

    // 4.1 Plan balanced boarding menu with Kcal/P:L:G target
    const menuCreateRes = await apiContext.post('/api/v1/boarding/menus', {
      headers: { 'X-User-Role': 'KE_TOAN' },
      data: {
        menu_date: '2026-11-05',
        meal_type: 'LUNCH',
        dishes: [
          { name: 'Thịt bò xào rau củ', weight_g: 130, protein_g: 24, lipid_g: 16, glucid_g: 6 },
          { name: 'Canh bí đỏ hầm xương', weight_g: 160, protein_g: 7, lipid_g: 3, glucid_g: 12 },
          { name: 'Cơm trắng', weight_g: 200, protein_g: 6, lipid_g: 0.8, glucid_g: 80 },
        ],
        sample_retention_verified: true,
      },
    });
    const menuCreateData = await menuCreateRes.json();
    const createdMenuId = menuCreateData.data?.id;
    recordResult('Kế toán: Create Balanced Boarding Menu (Kcal & P:L:G Targets)',
      menuCreateRes.status() === 201 && menuCreateData.data?.total_kcal > 500,
      `Calories: ${menuCreateData.data?.total_kcal} Kcal, Menu ID: ${createdMenuId}`);

    // 4.2 Generate dynamic VietQR (EMVCo standard)
    const vietqrRes = await apiContext.post('/api/v1/integrations/vietqr/generate', {
      headers: { 'X-User-Role': 'KE_TOAN' },
      data: {
        bank_bin: '970422',
        account_number: '0123456789',
        account_name: 'Truong Tieu Hoc EduPrimary',
        amount: 660000,
        description: 'TIEN AN BAN TRU THANG 11',
      },
    });
    const vietqrData = await vietqrRes.json();
    recordResult('Kế toán: Generate Dynamic VietQR EMVCo Payment Payload',
      vietqrRes.status() === 200 && vietqrData.data?.qr_string.startsWith('000201010212'),
      `EMVCo String: ${vietqrData.data?.qr_string.substring(0, 24)}...`);

    // 4.3 Issue electronic invoice under Decree 123/2020/NĐ-CP
    const einvoiceRes = await apiContext.post('/api/v1/integrations/einvoice/issue', {
      headers: { 'X-User-Role': 'KE_TOAN' },
      data: {
        school_id: SCHOOL_ID,
        buyer_name: 'Trần Văn Hùng',
        student_name: 'Trần Bảo Nam',
        student_class: '5A1',
        items: [
          { name: 'Tiền ăn bán trú tháng 11/2026', unit: 'Tháng', quantity: 1, unit_price: 660000, vat_rate: 0 },
        ],
      },
    });
    const einvoiceData = await einvoiceRes.json();
    recordResult('Kế toán: Issue Electronic Invoice (Decree 123 Serial 1C26TAA)',
      einvoiceRes.status() === 200 && einvoiceData.data?.invoice_series === '1C26TAA',
      `Series: ${einvoiceData.data?.invoice_series}, Status: ${einvoiceData.data?.status}`);

    // -----------------------------------------------------------------
    // ROLE 5: Y TẾ HỌC ĐƯỜNG (Y_TE)
    // -----------------------------------------------------------------
    console.log('\n[ROLE 5: Y tế - School Nurse Workflows]');

    // 5.1 Verify 24h food sample retention compliance (QĐ 1246/QĐ-BYT at 4°C)
    const targetMenuId = createdMenuId || 'm0000000-0000-0000-0000-000000000001';
    const sampleRetentionRes = await apiContext.post(
      `/api/v1/boarding/menus/${targetMenuId}/sample-retention`,
      {
        headers: { 'X-User-Role': 'Y_TE' },
        data: {
          sampled_at: '2026-11-05T11:00:00Z',
          storage_temperature: 4.0,
          keeper_name: 'Nguyễn Thị Hoa (Y tế học đường)',
        },
      }
    );
    const sampleRetentionData = await sampleRetentionRes.json();
    const isRetentionVerified = sampleRetentionRes.status() === 200
      && (sampleRetentionData.message || '').includes('1246/QĐ-BYT')
      && sampleRetentionData.data?.storage_temperature === 4.0;
    recordResult('Y tế: Verify 24h Food Sample Retention (QĐ 1246/QĐ-BYT @ 4°C)',
      isRetentionVerified,
      `Storage: ${sampleRetentionData.data?.storage_temperature}°C, Message: ${sampleRetentionData.message}`);

    // 5.2 Record school clinic first aid treatment
    const medicalEventRes = await apiContext.post('/api/v1/students/c0000000-0000-0000-0000-000000000001/health/medical-events', {
      headers: { 'X-User-Role': 'Y_TE' },
      data: {
        event_type: 'SƠ_CẤP_CỨU',
        symptoms: 'Trầy xước nhẹ khuỷu tay trong giờ thể dục',
        treatment: 'Sát trùng dung dịch Povidine và băng gạc vô trùng',
        parent_notified: true,
      },
    });
    recordResult('Y tế: Record First Aid Clinic Treatment & Parent Alert',
      medicalEventRes.status() === 201,
      `Status ${medicalEventRes.status()}`);

    // -----------------------------------------------------------------
    // ROLE 6: VĂN THƯ (VAN_THU)
    // -----------------------------------------------------------------
    console.log('\n[ROLE 6: Văn thư - School Registrar Workflows]');

    // 6.1 Decree 13/2023 Child Consent Enforcement: Reject onboarding without consent (422)
    const noConsentRes = await apiContext.post('/api/v1/students/onboard', {
      headers: { 'X-User-Role': 'VAN_THU' },
      data: {
        student_code: 'HS9099',
        full_name: 'Lê Bảo Châu',
        birth_date: '2019-03-12T00:00:00Z',
        guardian_name: 'Lê Văn Trọng',
        guardian_phone: '0912999888',
        has_digital_consent: false,
      },
    });
    recordResult('Văn thư: Decree 13 Article 20 Enforcement (Reject Missing Consent)',
      noConsentRes.status() === 422,
      'HTTP 422 Unprocessable Entity');

    // 6.2 Accept onboarding with verified digital consent
    const consentRes = await apiContext.post('/api/v1/students/onboard', {
      headers: { 'X-User-Role': 'VAN_THU' },
      data: {
        student_code: 'HS9099',
        full_name: 'Lê Bảo Châu',
        national_id: '001206099888',
        birth_date: '2019-03-12T00:00:00Z',
        guardian_name: 'Lê Văn Trọng',
        guardian_national_id: '001075088777',
        guardian_phone: '0912999888',
        guardian_relationship: 'FATHER',
        has_digital_consent: true,
        consent_method: 'VNEID_SSO',
      },
    });
    recordResult('Văn thư: Onboard Student With Digital Consent (Decree 13 Compliant)',
      consentRes.status() === 201,
      'HTTP 201 Created');

    // 6.3 VNeID Level 2 biometric verification
    const vneidRes = await apiContext.post('/api/v1/integrations/vneid/verify', {
      headers: { 'X-User-Role': 'VAN_THU' },
      data: {
        national_id: '001202012345',
        full_name: 'Nguyễn Văn An',
        date_of_birth: '2018-05-15',
        gender: 'MALE',
        consent_granted: true,
      },
    });
    const vneidData = await vneidRes.json();
    recordResult('Văn thư: VNeID Level 2 Biometric Identity Verification',
      vneidRes.status() === 200 && vneidData.data?.verification_level === 'LEVEL_2' && vneidData.data?.biometric_matched === true,
      `Level: ${vneidData.data?.verification_level}, Biometric: ${vneidData.data?.biometric_matched}`);

    // -----------------------------------------------------------------
    // ROLE 7: PHỤ HUYNH HỌC SINH (PHHS)
    // -----------------------------------------------------------------
    console.log('\n[ROLE 7: Phụ huynh - Parents & Guardians Workflows]');

    // 7.1 Cancel lunch meal before 8:00 AM -> Refund eligible
    const cancelBefore8Res = await apiContext.post('/api/v1/boarding/cancel-meal', {
      headers: { 'X-User-Role': 'PHHS' },
      data: {
        student_id: 'c0000000-0000-0000-0000-000000000001',
        attendance_date: new Date().toISOString().split('T')[0],
        cancellation_hour: 7,
        meal_fee_per_day: 35000,
      },
    });
    const cancelBefore8Data = await cancelBefore8Res.json();
    recordResult('PHHS: Cancel Lunch Before 8:00 AM (Refund Eligible 35,000 VND)',
      cancelBefore8Res.status() === 200 && cancelBefore8Data.data?.is_refund_eligible === true && cancelBefore8Data.data?.refund_amount === 35000,
      `Refund Eligible: ${cancelBefore8Data.data?.is_refund_eligible}, Amount: ${cancelBefore8Data.data?.refund_amount}đ`);

    // 7.2 Cancel lunch meal after 8:00 AM -> No refund under kitchen preparation policy
    const cancelAfter8Res = await apiContext.post('/api/v1/boarding/cancel-meal', {
      headers: { 'X-User-Role': 'PHHS' },
      data: {
        student_id: 'c0000000-0000-0000-0000-000000000001',
        attendance_date: new Date().toISOString().split('T')[0],
        cancellation_hour: 9,
      },
    });
    const cancelAfter8Data = await cancelAfter8Res.json();
    recordResult('PHHS: Cancel Lunch After 8:00 AM (Kitchen Policy Enforced - No Refund)',
      cancelAfter8Res.status() === 200 && cancelAfter8Data.data?.is_refund_eligible === false,
      `Refund Eligible: ${cancelAfter8Data.data?.is_refund_eligible}`);

    // 7.3 Submit child prescription medication request for school nurse
    const prescriptionRes = await apiContext.post('/api/v1/students/c0000000-0000-0000-0000-000000000001/health/prescriptions', {
      headers: { 'X-User-Role': 'PHHS' },
      data: {
        medication_name: 'Siro Ho Prospan',
        dosage: '5ml sau ăn trưa',
        administration_time: '12:30',
        instructions: 'Pha với 20ml nước ấm',
        parent_notes: 'Cháu dễ nôn khi uống thuốc đặc',
      },
    });
    recordResult('PHHS: Submit Prescription Medication Request for Nurse',
      prescriptionRes.status() === 201,
      `Status ${prescriptionRes.status()}`);

    // 7.4 RBAC Security Constraint: PHHS forbidden from administrative KPI batch
    const forbiddenRes = await apiContext.post('/api/v1/kpi/evaluate-school-batch', {
      headers: { 'X-User-Role': 'PHHS' },
      data: [],
    });
    recordResult('PHHS: RBAC Access Control Restriction (Cannot Batch KPI)',
      forbiddenRes.status() === 403,
      'HTTP 403 Forbidden');

    // =================================================================
    // SECTION 6: SCREENSHOT CATALOG AUDIT
    // =================================================================
    console.log('\n--- SECTION 6: Screenshot Artifacts Audit ---');
    const screenshotFiles = fs.readdirSync(SCREENSHOT_DIR).filter((f) => f.endsWith('.png'));
    console.log(`  Total screenshots captured in ${SCREENSHOT_DIR}: ${screenshotFiles.length}`);
    screenshotFiles.sort().forEach((file, idx) => {
      console.log(`    [${idx + 1}] ${file}`);
    });

    const hasSufficientScreenshots = screenshotFiles.length >= 15;
    recordResult('Visual Screenshots Artifacts Catalog (>= 15 Screenshots)',
      hasSufficientScreenshots,
      `Cataloged ${screenshotFiles.length} screenshots (VAL-BROWSER-004 compliant)`);

  } finally {
    await apiContext.dispose();
    await browser.close();
    cleanupServer();
  }

  // =================================================================
  // SUMMARY REPORT
  // =================================================================
  console.log('\n================================================================');
  console.log(` Test Summary: ${results.passed} Passed, ${results.failed} Failed (Total ${results.cases.length})`);
  console.log('================================================================');

  if (results.failed > 0) {
    console.error('\x1b[31m[ERROR] Some browser E2E test cases failed.\x1b[0m');
    process.exit(1);
  } else {
    console.log('\x1b[32m[SUCCESS] All Playwright browser E2E dogfooding test cases passed! (VAL-BROWSER-003 & VAL-BROWSER-004 satisfied)\x1b[0m\n');
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
