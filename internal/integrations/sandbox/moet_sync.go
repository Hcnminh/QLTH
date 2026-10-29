package sandbox

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"qlth/internal/db"

	"github.com/google/uuid"
)

// MOETSyncRequest represents a batch sync request for the MOET CSDL Outbox.
type MOETSyncRequest struct {
	SchoolID     string   `json:"school_id" binding:"required"`
	AcademicYear string   `json:"academic_year" binding:"required"`
	SyncType     string   `json:"sync_type" binding:"required"` // "STUDENT_PROFILES", "TRANSCRIPT_SNAPSHOTS", "ALL"
	RecordIDs    []string `json:"record_ids,omitempty"`         // Optional list of specific record IDs
	BatchSize    int      `json:"batch_size,omitempty"`         // Default 50
}

// MOETQueueItem represents a single outbound payload processed in the outbox.
type MOETQueueItem struct {
	ItemID       string `json:"item_id"`
	RecordType   string `json:"record_type"` // "STUDENT", "TRANSCRIPT"
	RecordID     string `json:"record_id"`
	StudentName  string `json:"student_name,omitempty"`
	Action       string `json:"action"`      // "INSERT", "UPDATE"
	SyncStatus   string `json:"sync_status"` // "SYNCED", "FAILED"
	ResponseCode string `json:"response_code"`
	Message      string `json:"message"`
}

// MOETSyncResult represents the result of processing outbound items to MOET CSDL.
type MOETSyncResult struct {
	Environment       string          `json:"environment"` // Always "SANDBOX_STUB"
	Mode              string          `json:"mode"`        // "SANDBOX_STUB"
	BatchID           string          `json:"batch_id"`
	SchoolID          string          `json:"school_id"`
	AcademicYear      string          `json:"academic_year"`
	SyncType          string          `json:"sync_type"`
	Status            string          `json:"status"` // "COMPLETED", "PARTIAL"
	TotalProcessed    int             `json:"total_processed"`
	SuccessCount      int             `json:"success_count"`
	FailedCount       int             `json:"failed_count"`
	MOETTransactionID string          `json:"moet_transaction_id"`
	ProcessedAt       string          `json:"processed_at"`
	ItemsProcessed    []MOETQueueItem `json:"items_processed"`
	LegalNotice       string          `json:"legal_notice"`
}

// MOETQueueStatus holds outbox queue statistics.
type MOETQueueStatus struct {
	Environment        string `json:"environment"`
	SchoolID           string `json:"school_id"`
	PendingStudents    int    `json:"pending_students"`
	PendingTranscripts int    `json:"pending_transcripts"`
	TotalQueued        int    `json:"total_queued"`
	LastSyncedAt       string `json:"last_synced_at"`
	Status             string `json:"status"` // "HEALTHY", "SYNCING"
}

// MOETSyncService manages synchronization of educational records with the MOET CSDL sandbox.
type MOETSyncService struct {
	mu           sync.RWMutex
	lastSyncedAt time.Time
}

// NewMOETSyncService initializes a new MOET CSDL sync service.
func NewMOETSyncService() *MOETSyncService {
	return &MOETSyncService{
		lastSyncedAt: time.Now().Add(-2 * time.Hour),
	}
}

// ProcessOutboxSync simulates sending pending students and transcript snapshots to MOET CSDL.
func (s *MOETSyncService) ProcessOutboxSync(ctx context.Context, req MOETSyncRequest) (*MOETSyncResult, error) {
	if strings.TrimSpace(req.SchoolID) == "" {
		return nil, fmt.Errorf("school_id không được để trống")
	}
	if strings.TrimSpace(req.AcademicYear) == "" {
		return nil, fmt.Errorf("academic_year không được để trống")
	}

	syncType := strings.ToUpper(strings.TrimSpace(req.SyncType))
	if syncType == "" {
		syncType = "ALL"
	}

	batchSize := req.BatchSize
	if batchSize <= 0 {
		batchSize = 50
	}

	items := s.collectQueueItems(ctx, req, syncType, batchSize)

	now := time.Now().UTC()
	s.mu.Lock()
	s.lastSyncedAt = now
	s.mu.Unlock()

	batchID := fmt.Sprintf("MOET-BATCH-%s-%06d", now.Format("20060102"), now.Unix()%1000000)
	txID := fmt.Sprintf("TX-MOET-CSDL-%s-%s", now.Format("20060102150405"), uuid.New().String()[:8])

	return &MOETSyncResult{
		Environment:       "SANDBOX_STUB",
		Mode:              "SANDBOX_STUB",
		BatchID:           batchID,
		SchoolID:          req.SchoolID,
		AcademicYear:      req.AcademicYear,
		SyncType:          syncType,
		Status:            "COMPLETED",
		TotalProcessed:    len(items),
		SuccessCount:      len(items),
		FailedCount:       0,
		MOETTransactionID: txID,
		ProcessedAt:       now.Format(time.RFC3339),
		ItemsProcessed:    items,
		LegalNotice:       "Mô phỏng đồng bộ hàng đợi Outbox lên CSDL Ngành Giáo dục và Đào tạo theo Thông tư 42/2021/TT-BGDĐT",
	}, nil
}

func (s *MOETSyncService) collectQueueItems(ctx context.Context, req MOETSyncRequest, syncType string, batchSize int) []MOETQueueItem {
	var items []MOETQueueItem
	database, err := db.GetDB()
	if err == nil && database != nil {
		schoolUUID, _ := uuid.Parse(req.SchoolID)
		if syncType == "STUDENT_PROFILES" || syncType == "ALL" {
			items = append(items, fetchStudentItems(ctx, database, schoolUUID, batchSize)...)
		}
		if syncType == "TRANSCRIPT_SNAPSHOTS" || syncType == "ALL" {
			items = append(items, fetchTranscriptItems(ctx, database, req.AcademicYear, batchSize)...)
		}
	}

	if len(items) == 0 {
		items = generateFallbackItems(req, syncType)
	}
	return items
}

func fetchStudentItems(ctx context.Context, database *sql.DB, schoolUUID uuid.UUID, limit int) []MOETQueueItem {
	var list []MOETQueueItem
	query := `SELECT id, full_name, student_code FROM students WHERE school_id = $1 LIMIT $2`
	rows, err := database.QueryContext(ctx, query, schoolUUID, limit)
	if err != nil {
		return list
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		var name, code string
		if rows.Scan(&id, &name, &code) == nil {
			list = append(list, MOETQueueItem{
				ItemID:       uuid.New().String(),
				RecordType:   "STUDENT",
				RecordID:     id.String(),
				StudentName:  name,
				Action:       "UPSERT",
				SyncStatus:   "SYNCED",
				ResponseCode: "200_MOET_OK",
				Message:      fmt.Sprintf("Đồng bộ hồ sơ học sinh %s (%s) lên CSDL ngành thành công", name, code),
			})
		}
	}
	return list
}

func fetchTranscriptItems(ctx context.Context, database *sql.DB, academicYear string, limit int) []MOETQueueItem {
	var list []MOETQueueItem
	query := `
		SELECT t.id, t.student_id, COALESCE(s.full_name, 'Học sinh')
		FROM academic_transcripts_snapshots t
		LEFT JOIN students s ON s.id = t.student_id
		WHERE t.academic_year_id = $1
		LIMIT $2
	`
	rows, err := database.QueryContext(ctx, query, academicYear, limit)
	if err != nil {
		return list
	}
	defer rows.Close()

	for rows.Next() {
		var id, studentID uuid.UUID
		var sName string
		if rows.Scan(&id, &studentID, &sName) == nil {
			list = append(list, MOETQueueItem{
				ItemID:       uuid.New().String(),
				RecordType:   "TRANSCRIPT",
				RecordID:     id.String(),
				StudentName:  sName,
				Action:       "INSERT",
				SyncStatus:   "SYNCED",
				ResponseCode: "200_MOET_OK",
				Message:      fmt.Sprintf("Đồng bộ học bạ số năm học %s học sinh %s thành công", academicYear, sName),
			})
		}
	}
	return list
}

func generateFallbackItems(req MOETSyncRequest, syncType string) []MOETQueueItem {
	var items []MOETQueueItem
	if len(req.RecordIDs) > 0 {
		for _, recID := range req.RecordIDs {
			recType := "STUDENT"
			if syncType == "TRANSCRIPT_SNAPSHOTS" {
				recType = "TRANSCRIPT"
			}
			items = append(items, MOETQueueItem{
				ItemID:       uuid.New().String(),
				RecordType:   recType,
				RecordID:     recID,
				StudentName:  "Nguyễn Văn An",
				Action:       "UPSERT",
				SyncStatus:   "SYNCED",
				ResponseCode: "200_MOET_OK",
				Message:      "Đồng bộ bản ghi chỉ định lên CSDL ngành Giáo dục thành công",
			})
		}
		return items
	}

	return []MOETQueueItem{
		{
			ItemID:       uuid.New().String(),
			RecordType:   "STUDENT",
			RecordID:     "c0000000-0000-0000-0000-000000000001",
			StudentName:  "Nguyễn Hoàng Quân",
			Action:       "UPSERT",
			SyncStatus:   "SYNCED",
			ResponseCode: "200_MOET_OK",
			Message:      "Đồng bộ hồ sơ định danh học sinh lên CSDL MOET thành công",
		},
		{
			ItemID:       uuid.New().String(),
			RecordType:   "TRANSCRIPT",
			RecordID:     uuid.New().String(),
			StudentName:  "Nguyễn Hoàng Quân",
			Action:       "INSERT",
			SyncStatus:   "SYNCED",
			ResponseCode: "200_MOET_OK",
			Message:      fmt.Sprintf("Đồng bộ học bạ số năm học %s kèm xác thực chữ ký SmartCA lên MOET thành công", req.AcademicYear),
		},
	}
}

// GetQueueStatus inspects pending items in the outbound queue.
func (s *MOETSyncService) GetQueueStatus(ctx context.Context, schoolID string) (*MOETQueueStatus, error) {
	s.mu.RLock()
	lastSync := s.lastSyncedAt
	s.mu.RUnlock()

	pendingStudents := 0
	pendingTranscripts := 0

	database, err := db.GetDB()
	if err == nil && database != nil {
		var sCount, tCount int
		_ = database.QueryRowContext(ctx, "SELECT COUNT(*) FROM students").Scan(&sCount)
		_ = database.QueryRowContext(ctx, "SELECT COUNT(*) FROM academic_transcripts_snapshots WHERE teacher_signature IS NOT NULL").Scan(&tCount)
		pendingStudents = sCount
		pendingTranscripts = tCount
	} else {
		pendingStudents = 3
		pendingTranscripts = 2
	}

	return &MOETQueueStatus{
		Environment:        "SANDBOX_STUB",
		SchoolID:           schoolID,
		PendingStudents:    pendingStudents,
		PendingTranscripts: pendingTranscripts,
		TotalQueued:        pendingStudents + pendingTranscripts,
		LastSyncedAt:       lastSync.Format(time.RFC3339),
		Status:             "HEALTHY",
	}, nil
}
