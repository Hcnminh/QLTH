package teaching

import "math"

// TeacherWorkloadProfile represents the parameters of a teacher in an academic year
type TeacherWorkloadProfile struct {
	TeacherID            string
	TeacherName          string
	SchoolType           string // 'PRIMARY', 'ETHNIC_BOARDING_PRIMARY' (PTDTBT), 'SPECIAL_EDUCATION'
	IsHomeroom           bool   // Giáo viên chủ nhiệm
	IsHeadSubject        bool   // Tổ trưởng chuyên môn
	IsDeputyHeadSubject  bool   // Tổ phó chuyên môn
	HasYoungChild        bool   // Nuôi con nhỏ dưới 12 tháng tuổi
	CounselingHoursWeek  int    // Tiết kiêm nhiệm tư vấn tâm lý học sinh (3-8 tiết)
	InclusiveTeachingWk  int    // Số tiết dạy hòa nhập trẻ khuyết tật
	PreGrade1EthnicHours int    // Số tiết dạy tiếng Việt cho trẻ DTTS trước lớp 1 hè
	ActualTaughtHours    int    // Tổng số tiết thực tế đã giảng dạy cả năm
	WeeksInSchoolYear    int    // Số tuần thực dạy (chuẩn là 35 tuần theo Điều 5 TT 05/2025)
}

// WorkloadQuotaCalculationResult stores detailed calculation of teaching quota
type WorkloadQuotaCalculationResult struct {
	StandardAnnualHours int     `json:"standard_annual_hours"` // 805 tiết đối với tiểu học thông thường
	WeeklyReducedHours  int     `json:"weekly_reduced_hours"`  // Số tiết giảm trừ hàng tuần
	AnnualReducedHours  int     `json:"annual_reduced_hours"`  // Tổng số tiết giảm cả năm
	NetAnnualQuota      int     `json:"net_annual_quota"`      // Định mức thực dạy cả năm phải đạt
	InclusiveBonusHours float64 `json:"inclusive_bonus_hours"` // Tiết quy đổi dạy hòa nhập
	TotalEffectiveHours float64 `json:"total_effective_hours"` // Tổng số tiết thực hiện quy đổi
	OvertimeHours       float64 `json:"overtime_hours"`        // Số tiết vượt định mức (thừa giờ tính tiền)
}

// QuotaEngine handles calculations under Circular 05/2025/TT-BGDĐT
type QuotaEngine struct{}

func NewQuotaEngine() *QuotaEngine {
	return &QuotaEngine{}
}

// CalculateQuota computes standard quota, reductions, and overtime hours
func (qe *QuotaEngine) CalculateQuota(profile TeacherWorkloadProfile) WorkloadQuotaCalculationResult {
	weeks := profile.WeeksInSchoolYear
	if weeks <= 0 {
		weeks = 35 // Chuẩn 35 tuần theo Điều 5 TT 05/2025
	}

	// 1. Xác định định mức chuẩn hàng năm (Điều 6 TT 05/2025)
	standardAnnual := 805 // 23 tiết/tuần * 35 tuần
	if profile.SchoolType == "ETHNIC_BOARDING_PRIMARY" || profile.SchoolType == "SPECIAL_EDUCATION" {
		standardAnnual = 21 * weeks // 21 tiết/tuần * 35 tuần = 735 tiết
	}

	// 2. Tính số tiết giảm trừ hàng tuần (Điều 9 & 10 TT 05/2025)
	weeklyReduction := 0
	if profile.IsHomeroom {
		weeklyReduction += 4 // GVCN tiểu học được giảm 4 tiết/tuần (Mới theo TT 05/2025)
	}
	if profile.IsHeadSubject {
		weeklyReduction += 3 // Tổ trưởng chuyên môn giảm 3 tiết/tuần
	}
	if profile.IsDeputyHeadSubject {
		weeklyReduction += 1 // Tổ phó chuyên môn giảm 1 tiết/tuần
	}
	if profile.HasYoungChild {
		weeklyReduction += 3 // Nữ giáo viên nuôi con nhỏ dưới 12 tháng giảm 3 tiết/tuần
	}
	if profile.CounselingHoursWeek > 0 {
		weeklyReduction += profile.CounselingHoursWeek // Tư vấn tâm lý
	}

	annualReduced := weeklyReduction * weeks
	netQuota := standardAnnual - annualReduced
	if netQuota < 0 {
		netQuota = 0
	}

	// 3. Quy đổi các hoạt động dạy học đặc thù (Điều 12 & 13 TT 05/2025)
	// Dạy hè tiếng Việt cho trẻ DTTS trước lớp 1: hệ số 1.5
	summerEthnicBonus := float64(profile.PreGrade1EthnicHours) * 1.5

	// Dạy hòa nhập học sinh khuyết tật: quy đổi phụ cấp giờ dạy (hệ số 0.2 trên tiết dạy lớp hòa nhập)
	inclusiveBonus := float64(profile.InclusiveTeachingWk*weeks) * 0.2

	totalEffective := float64(profile.ActualTaughtHours) + summerEthnicBonus + inclusiveBonus

	// 4. Tính toán số tiết thừa giờ
	overtime := totalEffective - float64(netQuota)
	if overtime < 0 {
		overtime = 0
	}

	return WorkloadQuotaCalculationResult{
		StandardAnnualHours: standardAnnual,
		WeeklyReducedHours:  weeklyReduction,
		AnnualReducedHours:  annualReduced,
		NetAnnualQuota:      netQuota,
		InclusiveBonusHours: math.Round(inclusiveBonus*10) / 10,
		TotalEffectiveHours: math.Round(totalEffective*10) / 10,
		OvertimeHours:       math.Round(overtime*10) / 10,
	}
}
