package curriculum

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDayOfWeek       = errors.New("thứ trong tuần không hợp lệ (phải từ Thứ 2 đến Thứ 7)")
	ErrInvalidPeriod          = errors.New("tiết học không hợp lệ (phải từ tiết 1 đến tiết 10)")
	ErrInvalidWeekRange       = errors.New("tuần áp dụng không hợp lệ (0 cho cả kỳ, hoặc từ tuần 1 đến tuần 35)")
	ErrTeacherConflict        = errors.New("xung đột thời khóa biểu giáo viên theo Thông tư 05/2025/TT-BGDĐT")
	ErrClassConflict          = errors.New("xung đột thời khóa biểu lớp học")
	ErrRoomConflict           = errors.New("xung đột phòng học")
	ErrEmptyTimetableRequired = errors.New("thiếu thông tin bắt buộc của tiết học (trường, năm học, học kỳ, lớp, giáo viên, môn học)")
)

// TimetableSlot represents a scheduled teaching slot
type TimetableSlot struct {
	ID             uuid.UUID `json:"id"`
	SchoolID       uuid.UUID `json:"school_id"`
	AcademicYearID string    `json:"academic_year_id"` // '2026-2027'
	TermID         string    `json:"term_id"`          // '2026-2027-HK1'
	ClassID        uuid.UUID `json:"class_id"`
	ClassName      string    `json:"class_name,omitempty"`
	TeacherID      uuid.UUID `json:"teacher_id"`
	TeacherName    string    `json:"teacher_name,omitempty"`
	SubjectID      string    `json:"subject_id"`
	DayOfWeek      int       `json:"day_of_week"` // 2 (Thứ 2) to 7 (Thứ 7)
	Period         int       `json:"period"`      // 1-5 (Sáng), 6-10 (Chiều)
	RoomNumber     string    `json:"room_number,omitempty"`
	WeekNumber     int       `json:"week_number"` // 0 = toan ky, 1-35 = tuan cu the
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TimetableConflict represents a detected schedule collision
type TimetableConflict struct {
	ConflictType string        `json:"conflict_type"` // 'TEACHER_OVERLAP', 'CLASS_OVERLAP', 'ROOM_OVERLAP'
	Message      string        `json:"message"`
	Conflicting  TimetableSlot `json:"conflicting_slot"`
	Attempted    TimetableSlot `json:"attempted_slot"`
}

func (c TimetableConflict) Error() string {
	return c.Message
}

// TimetableEngine handles schedule validation, conflict detection, and Circular 05/2025 quota adherence
type TimetableEngine struct{}

// NewTimetableEngine creates a new timetable engine instance
func NewTimetableEngine() *TimetableEngine {
	return &TimetableEngine{}
}

// ValidateSlotBounds verifies structural boundary requirements of a single slot
func (e *TimetableEngine) ValidateSlotBounds(slot TimetableSlot) error {
	if slot.SchoolID == uuid.Nil || slot.ClassID == uuid.Nil || slot.TeacherID == uuid.Nil ||
		strings.TrimSpace(slot.AcademicYearID) == "" || strings.TrimSpace(slot.TermID) == "" ||
		strings.TrimSpace(slot.SubjectID) == "" {
		return ErrEmptyTimetableRequired
	}
	if slot.DayOfWeek < 2 || slot.DayOfWeek > 7 {
		return ErrInvalidDayOfWeek
	}
	if slot.Period < 1 || slot.Period > 10 {
		return ErrInvalidPeriod
	}
	if slot.WeekNumber < 0 || slot.WeekNumber > 35 {
		return ErrInvalidWeekRange
	}
	return nil
}

// IsWeekOverlapping checks if two week specifications collide (0 matches any week)
func IsWeekOverlapping(w1, w2 int) bool {
	if w1 == 0 || w2 == 0 {
		return true
	}
	return w1 == w2
}

// CheckConflicts evaluates an attempted slot against all existing schedule slots
func (e *TimetableEngine) CheckConflicts(existing []TimetableSlot, newSlot TimetableSlot) *TimetableConflict {
	for _, slot := range existing {
		// Ignore self-comparison on update
		if slot.ID != uuid.Nil && newSlot.ID != uuid.Nil && slot.ID == newSlot.ID {
			continue
		}

		// Only slots in the same academic year and term can conflict
		if slot.AcademicYearID != newSlot.AcademicYearID || slot.TermID != newSlot.TermID {
			continue
		}

		// Must be on the exact same day of week and period
		if slot.DayOfWeek != newSlot.DayOfWeek || slot.Period != newSlot.Period {
			continue
		}

		// Must overlap in week applicability
		if !IsWeekOverlapping(slot.WeekNumber, newSlot.WeekNumber) {
			continue
		}

		// 1. Teacher schedule conflict (Circular 05/2025/TT-BGDĐT: no teacher overlapping hours)
		if slot.TeacherID == newSlot.TeacherID {
			teacherName := slot.TeacherName
			if teacherName == "" {
				teacherName = newSlot.TeacherName
			}
			if teacherName == "" {
				teacherName = slot.TeacherID.String()
			}

			clsName := slot.ClassName
			if clsName == "" {
				clsName = slot.ClassID.String()
			}

			msg := fmt.Sprintf(
				"Xung đột lịch giảng dạy: Giáo viên %s đã được phân công dạy tại lớp %s (Thứ %d, Tiết %d) theo Thông tư 05/2025/TT-BGDĐT",
				teacherName, clsName, slot.DayOfWeek, slot.Period,
			)
			return &TimetableConflict{
				ConflictType: "TEACHER_OVERLAP",
				Message:      msg,
				Conflicting:  slot,
				Attempted:    newSlot,
			}
		}

		// 2. Class schedule conflict (A class cannot have two concurrent lessons)
		if slot.ClassID == newSlot.ClassID {
			clsName := slot.ClassName
			if clsName == "" {
				clsName = newSlot.ClassName
			}
			if clsName == "" {
				clsName = slot.ClassID.String()
			}

			msg := fmt.Sprintf(
				"Xung đột lịch lớp học: Lớp %s đã được xếp tiết môn %s (Thứ %d, Tiết %d)",
				clsName, slot.SubjectID, slot.DayOfWeek, slot.Period,
			)
			return &TimetableConflict{
				ConflictType: "CLASS_OVERLAP",
				Message:      msg,
				Conflicting:  slot,
				Attempted:    newSlot,
			}
		}

		// 3. Room conflict (If both assign a room number and they match)
		if strings.TrimSpace(slot.RoomNumber) != "" &&
			strings.TrimSpace(newSlot.RoomNumber) != "" &&
			strings.EqualFold(strings.TrimSpace(slot.RoomNumber), strings.TrimSpace(newSlot.RoomNumber)) {
			msg := fmt.Sprintf(
				"Xung đột phòng học: Phòng %s đã được sử dụng bởi lớp %s (Thứ %d, Tiết %d)",
				slot.RoomNumber, slot.ClassName, slot.DayOfWeek, slot.Period,
			)
			return &TimetableConflict{
				ConflictType: "ROOM_OVERLAP",
				Message:      msg,
				Conflicting:  slot,
				Attempted:    newSlot,
			}
		}
	}

	return nil
}

// ValidateAndAddSlot validates slot bounds and ensures zero scheduling conflicts
func (e *TimetableEngine) ValidateAndAddSlot(existing []TimetableSlot, newSlot TimetableSlot) (TimetableSlot, error) {
	if err := e.ValidateSlotBounds(newSlot); err != nil {
		return newSlot, err
	}

	if conflict := e.CheckConflicts(existing, newSlot); conflict != nil {
		return newSlot, conflict
	}

	if newSlot.ID == uuid.Nil {
		newSlot.ID = uuid.New()
	}
	now := time.Now()
	if newSlot.CreatedAt.IsZero() {
		newSlot.CreatedAt = now
	}
	newSlot.UpdatedAt = now

	return newSlot, nil
}

// CalculateTeacherWeeklyPeriods counts weekly periods assigned to a teacher
func (e *TimetableEngine) CalculateTeacherWeeklyPeriods(slots []TimetableSlot, teacherID uuid.UUID) int {
	count := 0
	for _, slot := range slots {
		if slot.TeacherID == teacherID {
			count++
		}
	}
	return count
}
