package curriculum_test

import (
	"strings"
	"testing"

	"qlth/internal/core/curriculum"

	"github.com/google/uuid"
)

func TestTimetableEngine_ValidateSlotBounds(t *testing.T) {
	engine := curriculum.NewTimetableEngine()
	validSlot := curriculum.TimetableSlot{
		SchoolID:       uuid.New(),
		AcademicYearID: "2026-2027",
		TermID:         "2026-2027-HK1",
		ClassID:        uuid.New(),
		TeacherID:      uuid.New(),
		SubjectID:      "TOAN",
		DayOfWeek:      2, // Thứ 2
		Period:         1, // Tiết 1
		WeekNumber:     0, // Cả kỳ
	}

	if err := engine.ValidateSlotBounds(validSlot); err != nil {
		t.Fatalf("expected valid slot, got %v", err)
	}

	// Invalid day of week (Sunday is 1 or 8, but we require 2-7)
	slotInvalidDay := validSlot
	slotInvalidDay.DayOfWeek = 1
	if err := engine.ValidateSlotBounds(slotInvalidDay); err != curriculum.ErrInvalidDayOfWeek {
		t.Errorf("expected ErrInvalidDayOfWeek, got %v", err)
	}

	// Invalid period (period 11)
	slotInvalidPeriod := validSlot
	slotInvalidPeriod.Period = 11
	if err := engine.ValidateSlotBounds(slotInvalidPeriod); err != curriculum.ErrInvalidPeriod {
		t.Errorf("expected ErrInvalidPeriod, got %v", err)
	}

	// Invalid week range (week 36)
	slotInvalidWeek := validSlot
	slotInvalidWeek.WeekNumber = 36
	if err := engine.ValidateSlotBounds(slotInvalidWeek); err != curriculum.ErrInvalidWeekRange {
		t.Errorf("expected ErrInvalidWeekRange, got %v", err)
	}

	// Missing required fields
	slotMissingTeacher := validSlot
	slotMissingTeacher.TeacherID = uuid.Nil
	if err := engine.ValidateSlotBounds(slotMissingTeacher); err != curriculum.ErrEmptyTimetableRequired {
		t.Errorf("expected ErrEmptyTimetableRequired, got %v", err)
	}
}

func TestTimetableEngine_TeacherOverlapConflict(t *testing.T) {
	engine := curriculum.NewTimetableEngine()
	teacherID := uuid.New()
	schoolID := uuid.New()
	class1ID := uuid.New()
	class2ID := uuid.New()

	existing := []curriculum.TimetableSlot{
		{
			ID:             uuid.New(),
			SchoolID:       schoolID,
			AcademicYearID: "2026-2027",
			TermID:         "2026-2027-HK1",
			ClassID:        class1ID,
			ClassName:      "5A1",
			TeacherID:      teacherID,
			TeacherName:    "Thầy Nguyễn Văn An",
			SubjectID:      "TOAN",
			DayOfWeek:      2, // Thứ 2
			Period:         1, // Tiết 1
			WeekNumber:     0,
		},
	}

	// Attempt to schedule the SAME teacher in class 5A2 on Thứ 2, Tiết 1
	collidingSlot := curriculum.TimetableSlot{
		SchoolID:       schoolID,
		AcademicYearID: "2026-2027",
		TermID:         "2026-2027-HK1",
		ClassID:        class2ID,
		ClassName:      "5A2",
		TeacherID:      teacherID,
		TeacherName:    "Thầy Nguyễn Văn An",
		SubjectID:      "TOAN",
		DayOfWeek:      2,
		Period:         1,
		WeekNumber:     0,
	}

	conflict := engine.CheckConflicts(existing, collidingSlot)
	if conflict == nil {
		t.Fatalf("expected teacher schedule conflict, got nil")
	}

	if conflict.ConflictType != "TEACHER_OVERLAP" {
		t.Errorf("expected TEACHER_OVERLAP, got %s", conflict.ConflictType)
	}
	if !strings.Contains(conflict.Message, "Thông tư 05/2025/TT-BGDĐT") {
		t.Errorf("expected conflict message to cite Thông tư 05/2025, got %s", conflict.Message)
	}
}

func TestTimetableEngine_ClassOverlapConflict(t *testing.T) {
	engine := curriculum.NewTimetableEngine()
	schoolID := uuid.New()
	classID := uuid.New()
	teacher1ID := uuid.New()
	teacher2ID := uuid.New()

	existing := []curriculum.TimetableSlot{
		{
			ID:             uuid.New(),
			SchoolID:       schoolID,
			AcademicYearID: "2026-2027",
			TermID:         "2026-2027-HK1",
			ClassID:        classID,
			ClassName:      "1A1",
			TeacherID:      teacher1ID,
			SubjectID:      "TIENG_VIET",
			DayOfWeek:      3, // Thứ 3
			Period:         2, // Tiết 2
			WeekNumber:     0,
		},
	}

	// Attempt to schedule another subject for class 1A1 at the exact same period
	collidingSlot := curriculum.TimetableSlot{
		SchoolID:       schoolID,
		AcademicYearID: "2026-2027",
		TermID:         "2026-2027-HK1",
		ClassID:        classID,
		ClassName:      "1A1",
		TeacherID:      teacher2ID,
		SubjectID:      "TOAN",
		DayOfWeek:      3,
		Period:         2,
		WeekNumber:     0,
	}

	conflict := engine.CheckConflicts(existing, collidingSlot)
	if conflict == nil {
		t.Fatalf("expected class schedule conflict, got nil")
	}
	if conflict.ConflictType != "CLASS_OVERLAP" {
		t.Errorf("expected CLASS_OVERLAP, got %s", conflict.ConflictType)
	}
}

func TestTimetableEngine_RoomOverlapConflict(t *testing.T) {
	engine := curriculum.NewTimetableEngine()
	schoolID := uuid.New()

	existing := []curriculum.TimetableSlot{
		{
			ID:             uuid.New(),
			SchoolID:       schoolID,
			AcademicYearID: "2026-2027",
			TermID:         "2026-2027-HK1",
			ClassID:        uuid.New(),
			ClassName:      "5A1",
			TeacherID:      uuid.New(),
			SubjectID:      "TIN_HOC",
			DayOfWeek:      4, // Thứ 4
			Period:         3, // Tiết 3
			RoomNumber:     "Phòng Tin Học 1",
			WeekNumber:     0,
		},
	}

	// Different teacher and different class, but same specialized room at same time
	collidingSlot := curriculum.TimetableSlot{
		SchoolID:       schoolID,
		AcademicYearID: "2026-2027",
		TermID:         "2026-2027-HK1",
		ClassID:        uuid.New(),
		ClassName:      "5A2",
		TeacherID:      uuid.New(),
		SubjectID:      "TIN_HOC",
		DayOfWeek:      4,
		Period:         3,
		RoomNumber:     "Phòng Tin Học 1",
		WeekNumber:     0,
	}

	conflict := engine.CheckConflicts(existing, collidingSlot)
	if conflict == nil {
		t.Fatalf("expected room conflict, got nil")
	}
	if conflict.ConflictType != "ROOM_OVERLAP" {
		t.Errorf("expected ROOM_OVERLAP, got %s", conflict.ConflictType)
	}
}

func TestTimetableEngine_NonConflictingSlot_Success(t *testing.T) {
	engine := curriculum.NewTimetableEngine()
	teacherID := uuid.New()
	schoolID := uuid.New()

	existing := []curriculum.TimetableSlot{
		{
			ID:             uuid.New(),
			SchoolID:       schoolID,
			AcademicYearID: "2026-2027",
			TermID:         "2026-2027-HK1",
			ClassID:        uuid.New(),
			TeacherID:      teacherID,
			SubjectID:      "TOAN",
			DayOfWeek:      2,
			Period:         1, // Tiết 1
			WeekNumber:     0,
		},
	}

	// Same teacher, but Tiết 2 (different period) -> NO conflict
	newSlot := curriculum.TimetableSlot{
		SchoolID:       schoolID,
		AcademicYearID: "2026-2027",
		TermID:         "2026-2027-HK1",
		ClassID:        uuid.New(),
		TeacherID:      teacherID,
		SubjectID:      "TOAN",
		DayOfWeek:      2,
		Period:         2, // Tiết 2
		WeekNumber:     0,
	}

	saved, err := engine.ValidateAndAddSlot(existing, newSlot)
	if err != nil {
		t.Fatalf("expected non-conflicting slot to succeed, got %v", err)
	}
	if saved.ID == uuid.Nil {
		t.Errorf("expected saved slot to have ID")
	}

	// Calculate weekly periods for teacher
	allSlots := append(existing, saved)
	weeklyHours := engine.CalculateTeacherWeeklyPeriods(allSlots, teacherID)
	if weeklyHours != 2 {
		t.Errorf("expected 2 periods, got %d", weeklyHours)
	}
}
