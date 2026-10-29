package boarding

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// National Institute of Nutrition (Viện Dinh dưỡng Quốc gia) standards for Primary Pupils (6-11 years old)
const (
	// Bữa trưa tiểu học (Primary Lunch) chiếm 40-45% năng lượng cả ngày
	StandardLunchMinKcal = 600.0
	StandardLunchMaxKcal = 750.0

	// Bữa phụ / Xế (Snack) chiếm 10-15% năng lượng
	StandardSnackMinKcal = 150.0
	StandardSnackMaxKcal = 250.0

	// Bữa sáng (Breakfast)
	StandardBreakfastMinKcal = 350.0
	StandardBreakfastMaxKcal = 500.0

	// Tỷ lệ phân bổ năng lượng chuẩn P:L:G theo Viện Dinh Dưỡng
	StandardProteinMinPercent = 13.0
	StandardProteinMaxPercent = 20.0
	StandardLipidMinPercent   = 20.0
	StandardLipidMaxPercent   = 30.0
	StandardGlucidMinPercent  = 50.0
	StandardGlucidMaxPercent  = 65.0

	// Quy chuẩn lưu mẫu thức ăn 24h theo Quyết định 1246/QĐ-BYT Bộ Y Tế
	MinSampleRetentionHours = 24.0
	MinStorageTempCelsius   = 2.0
	MaxStorageTempCelsius   = 8.0
)

// DishItem represents an individual dish in the boarding menu
type DishItem struct {
	Name     string  `json:"name"`
	WeightG  float64 `json:"weight_g"`
	Kcal     float64 `json:"kcal,omitempty"`
	ProteinG float64 `json:"protein_g,omitempty"`
	LipidG   float64 `json:"lipid_g,omitempty"`
	GlucidG  float64 `json:"glucid_g,omitempty"`
	Category string  `json:"category,omitempty"` // MAIN, SOUP, VEGETABLE, RICE, DESSERT
}

// NutritionalBreakdown represents computed nutrition values and ratios
type NutritionalBreakdown struct {
	TotalKcal           float64 `json:"total_kcal"`
	ProteinG            float64 `json:"protein_g"`
	LipidG              float64 `json:"lipid_g"`
	GlucidG             float64 `json:"glucid_g"`
	ProteinPercent      float64 `json:"protein_percent"`
	LipidPercent        float64 `json:"lipid_percent"`
	GlucidPercent       float64 `json:"glucid_percent"`
	EnergyBalanceStatus string  `json:"energy_balance_status"` // CHUẨN_CÂN_ĐỐI, THIẾU_NĂNG_LƯỢNG, DƯ_THỪA_NĂNG_LƯỢNG
	PLGBalanceStatus    string  `json:"plg_balance_status"`    // CHUẨN_CÂN_ĐỐI, LỆCH_TỶ_LỆ
}

// NutritionComplianceResult represents adherence to National Institute of Nutrition guidelines
type NutritionComplianceResult struct {
	IsCompliant              bool                 `json:"is_compliant"`
	KcalCompliant            bool                 `json:"kcal_compliant"`
	RatioCompliant           bool                 `json:"ratio_compliant"`
	SampleRetentionCompliant bool                 `json:"sample_retention_compliant"`
	GuidelineStandard        string               `json:"guideline_standard"`
	Breakdown                NutritionalBreakdown `json:"breakdown"`
	Recommendations          []string             `json:"recommendations"`
}

// SampleRetentionVerification represents 24h food sample check under Decree/Decision 1246/QĐ-BYT
type SampleRetentionVerification struct {
	SampledAt          time.Time  `json:"sampled_at"`
	InspectionAt       time.Time  `json:"inspection_at"`
	StorageTemperature float64    `json:"storage_temperature"` // 2 - 8 độ C
	KeeperID           *uuid.UUID `json:"keeper_id,omitempty"`
	KeeperName         string     `json:"keeper_name,omitempty"`
	IsRetained24Hours  bool       `json:"is_retained_24_hours"`
	IsTempCompliant    bool       `json:"is_temp_compliant"`
	IsVerified         bool       `json:"is_verified"`
	Remarks            string     `json:"remarks"`
}

// NutritionEngine encapsulates nutrition calculations and compliance checks
type NutritionEngine struct{}

// NewNutritionEngine creates a new NutritionEngine instance
func NewNutritionEngine() *NutritionEngine {
	return &NutritionEngine{}
}

// round1 rounds a float64 to 1 decimal place
func round1(val float64) float64 {
	return math.Round(val*10.0) / 10.0
}

// CalculateDishNutrition computes total Kcal and P:L:G ratio from dish components
func (e *NutritionEngine) CalculateDishNutrition(dishes []DishItem) NutritionalBreakdown {
	var totalProteinG, totalLipidG, totalGlucidG, explicitKcal float64

	for _, d := range dishes {
		totalProteinG += d.ProteinG
		totalLipidG += d.LipidG
		totalGlucidG += d.GlucidG
		explicitKcal += d.Kcal
	}

	// 1g Protein = 4 Kcal, 1g Lipid = 9 Kcal, 1g Glucid = 4 Kcal
	calcKcal := (totalProteinG * 4.0) + (totalLipidG * 9.0) + (totalGlucidG * 4.0)
	finalKcal := calcKcal
	if finalKcal <= 0 && explicitKcal > 0 {
		finalKcal = explicitKcal
	}

	var pPct, lPct, gPct float64
	if finalKcal > 0 {
		pPct = round1(((totalProteinG * 4.0) / finalKcal) * 100.0)
		lPct = round1(((totalLipidG * 9.0) / finalKcal) * 100.0)
		gPct = round1(((totalGlucidG * 4.0) / finalKcal) * 100.0)
	}

	return NutritionalBreakdown{
		TotalKcal:      round1(finalKcal),
		ProteinG:       round1(totalProteinG),
		LipidG:         round1(totalLipidG),
		GlucidG:        round1(totalGlucidG),
		ProteinPercent: pPct,
		LipidPercent:   lPct,
		GlucidPercent:  gPct,
	}
}

// EvaluateNutritionalCompliance checks if the menu conforms to National Institute of Nutrition guidelines
func (e *NutritionEngine) EvaluateNutritionalCompliance(
	mealType string,
	kcal, pPct, lPct, gPct float64,
	sampleVerified bool,
) NutritionComplianceResult {
	var minKcal, maxKcal float64
	switch mealType {
	case "SNACK":
		minKcal, maxKcal = StandardSnackMinKcal, StandardSnackMaxKcal
	case "BREAKFAST":
		minKcal, maxKcal = StandardBreakfastMinKcal, StandardBreakfastMaxKcal
	default: // LUNCH
		minKcal, maxKcal = StandardLunchMinKcal, StandardLunchMaxKcal
	}

	var recommendations []string

	// 1. Kcal check
	kcalCompliant := kcal >= minKcal && kcal <= maxKcal
	energyStatus := "CHUẨN_CÂN_ĐỐI"
	if kcal < minKcal {
		energyStatus = "THIẾU_NĂNG_LƯỢNG"
		recommendations = append(recommendations,
			fmt.Sprintf("Năng lượng bữa ăn (%.1f Kcal) thấp hơn mức khuyến nghị tối thiểu của Viện Dinh Dưỡng (%.0f Kcal). Cần tăng khẩu phần thực phẩm chính hoặc ngũ cốc.", kcal, minKcal),
		)
	} else if kcal > maxKcal {
		energyStatus = "DƯ_THỪA_NĂNG_LƯỢNG"
		recommendations = append(recommendations,
			fmt.Sprintf("Năng lượng bữa ăn (%.1f Kcal) vượt mức khuyến nghị tối đa (%.0f Kcal). Cần giảm lượng dầu mỡ hoặc tinh bột để phòng ngừa thừa cân béo phì.", kcal, maxKcal),
		)
	}

	// 2. P:L:G ratio check
	pCompliant := pPct >= StandardProteinMinPercent && pPct <= StandardProteinMaxPercent
	lCompliant := lPct >= StandardLipidMinPercent && lPct <= StandardLipidMaxPercent
	gCompliant := gPct >= StandardGlucidMinPercent && gPct <= StandardGlucidMaxPercent
	ratioCompliant := pCompliant && lCompliant && gCompliant

	plgStatus := "CHUẨN_CÂN_ĐỐI"
	if !ratioCompliant {
		plgStatus = "LỆCH_TỶ_LỆ"
		if !pCompliant {
			recommendations = append(recommendations,
				fmt.Sprintf("Tỷ lệ Protein (%.1f%%) chưa đạt dải chuẩn [%.0f%% - %.0f%%] năng lượng. Điều chỉnh lượng thịt, cá, trứng, đậu đỗ.", pPct, StandardProteinMinPercent, StandardProteinMaxPercent),
			)
		}
		if !lCompliant {
			recommendations = append(recommendations,
				fmt.Sprintf("Tỷ lệ Lipid (%.1f%%) chưa đạt dải chuẩn [%.0f%% - %.0f%%] năng lượng. Điều chỉnh lượng dầu thực vật, mỡ hoặc món xào rán.", lPct, StandardLipidMinPercent, StandardLipidMaxPercent),
			)
		}
		if !gCompliant {
			recommendations = append(recommendations,
				fmt.Sprintf("Tỷ lệ Glucid (%.1f%%) chưa đạt dải chuẩn [%.0f%% - %.0f%%] năng lượng. Điều chỉnh lượng cơm gạo, củ hoặc tráng miệng ngọt.", gPct, StandardGlucidMinPercent, StandardGlucidMaxPercent),
			)
		}
	}

	// 3. Sample retention check (Quyết định 1246/QĐ-BYT)
	if !sampleVerified {
		recommendations = append(recommendations, "Chưa xác nhận lưu mẫu thức ăn 24h theo Quyết định 1246/QĐ-BYT của Bộ Y Tế.")
	}

	isCompliant := kcalCompliant && ratioCompliant && sampleVerified

	return NutritionComplianceResult{
		IsCompliant:              isCompliant,
		KcalCompliant:            kcalCompliant,
		RatioCompliant:           ratioCompliant,
		SampleRetentionCompliant: sampleVerified,
		GuidelineStandard:        "Viện Dinh dưỡng Quốc gia & Quyết định 1246/QĐ-BYT",
		Breakdown: NutritionalBreakdown{
			TotalKcal:           round1(kcal),
			ProteinPercent:      round1(pPct),
			LipidPercent:        round1(lPct),
			GlucidPercent:       round1(gPct),
			EnergyBalanceStatus: energyStatus,
			PLGBalanceStatus:    plgStatus,
		},
		Recommendations: recommendations,
	}
}

// VerifySampleRetention checks food sample 24-hour retention protocol
func (e *NutritionEngine) VerifySampleRetention(
	sampledAt, inspectionAt time.Time,
	temp float64,
	keeperID *uuid.UUID,
	keeperName string,
) SampleRetentionVerification {
	hoursElapsed := inspectionAt.Sub(sampledAt).Hours()
	is24Hours := hoursElapsed >= MinSampleRetentionHours
	isTempOk := temp >= MinStorageTempCelsius && temp <= MaxStorageTempCelsius
	isKeeperOk := keeperID != nil && *keeperID != uuid.Nil

	isVerified := is24Hours && isTempOk && isKeeperOk

	remarks := "Mẫu lưu thức ăn đạt chuẩn 24h và bảo quản đúng dải nhiệt độ 2-8°C theo QĐ 1246/QĐ-BYT."
	if !is24Hours {
		remarks = fmt.Sprintf("Thời gian lưu mẫu chưa đủ 24 giờ (mới đạt %.1f giờ).", hoursElapsed)
	} else if !isTempOk {
		remarks = fmt.Sprintf("Nhiệt độ bảo quản mẫu (%.1f°C) nằm ngoài dải quy chuẩn 2°C - 8°C.", temp)
	} else if !isKeeperOk {
		remarks = "Chưa chỉ định nhân viên phụ trách lưu trữ mẫu thức ăn."
	}

	return SampleRetentionVerification{
		SampledAt:          sampledAt,
		InspectionAt:       inspectionAt,
		StorageTemperature: temp,
		KeeperID:           keeperID,
		KeeperName:         keeperName,
		IsRetained24Hours:  is24Hours,
		IsTempCompliant:    isTempOk,
		IsVerified:         isVerified,
		Remarks:            remarks,
	}
}
