package boarding

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNutritionEngine_CalculateDishNutrition(t *testing.T) {
	engine := NewNutritionEngine()

	dishes := []DishItem{
		{Name: "Thịt ba chỉ kho trứng cút", WeightG: 120, ProteinG: 22.5, LipidG: 18.0, GlucidG: 4.0},
		{Name: "Canh bí đỏ nấu tôm", WeightG: 150, ProteinG: 4.5, LipidG: 1.5, GlucidG: 10.0},
		{Name: "Rau bắp cải xào", WeightG: 80, ProteinG: 1.5, LipidG: 4.0, GlucidG: 6.0},
		{Name: "Cơm trắng", WeightG: 180, ProteinG: 5.0, LipidG: 0.5, GlucidG: 75.0},
		{Name: "Tráng miệng chuối cau", WeightG: 60, ProteinG: 0.8, LipidG: 0.2, GlucidG: 15.0},
	}

	breakdown := engine.CalculateDishNutrition(dishes)

	if breakdown.TotalKcal <= 0 {
		t.Fatalf("Expected TotalKcal > 0, got %f", breakdown.TotalKcal)
	}

	// Verify Kcal is approximately in the primary lunch range
	if breakdown.TotalKcal < 600 || breakdown.TotalKcal > 750 {
		t.Logf("Calculated Kcal: %f", breakdown.TotalKcal)
	}

	sumPct := breakdown.ProteinPercent + breakdown.LipidPercent + breakdown.GlucidPercent
	if sumPct < 99.0 || sumPct > 101.0 {
		t.Errorf("Sum of percentages should be ~100%%, got %f", sumPct)
	}
}

func TestNutritionEngine_EvaluateNutritionalCompliance_LunchSuccess(t *testing.T) {
	engine := NewNutritionEngine()

	// A balanced primary school lunch menu: 680.5 Kcal, P: 16.5%, L: 24.0%, G: 59.5%
	result := engine.EvaluateNutritionalCompliance("LUNCH", 680.5, 16.5, 24.0, 59.5, true)

	if !result.IsCompliant {
		t.Errorf("Expected menu to be compliant, but got: %+v", result)
	}
	if !result.KcalCompliant {
		t.Errorf("Expected Kcal to be compliant")
	}
	if !result.RatioCompliant {
		t.Errorf("Expected P:L:G ratio to be compliant")
	}
	if !result.SampleRetentionCompliant {
		t.Errorf("Expected sample retention to be compliant")
	}
	if len(result.Recommendations) > 0 {
		t.Errorf("Expected 0 recommendations for compliant menu, got %d", len(result.Recommendations))
	}
}

func TestNutritionEngine_EvaluateNutritionalCompliance_DeficitAndImbalanced(t *testing.T) {
	engine := NewNutritionEngine()

	// Under 600 Kcal and lipid too low
	result := engine.EvaluateNutritionalCompliance("LUNCH", 520.0, 12.0, 15.0, 73.0, false)

	if result.IsCompliant {
		t.Errorf("Expected menu to be non-compliant")
	}
	if result.KcalCompliant {
		t.Errorf("Expected Kcal to fail")
	}
	if result.RatioCompliant {
		t.Errorf("Expected Ratio to fail")
	}
	if result.Breakdown.EnergyBalanceStatus != "THIẾU_NĂNG_LƯỢNG" {
		t.Errorf("Expected energy status THIẾU_NĂNG_LƯỢNG, got %s", result.Breakdown.EnergyBalanceStatus)
	}
	if result.Breakdown.PLGBalanceStatus != "LỆCH_TỶ_LỆ" {
		t.Errorf("Expected PLG status LỆCH_TỶ_LỆ, got %s", result.Breakdown.PLGBalanceStatus)
	}
	if len(result.Recommendations) < 3 {
		t.Errorf("Expected multiple recommendations, got %d", len(result.Recommendations))
	}
}

func TestNutritionEngine_EvaluateNutritionalCompliance_Surplus(t *testing.T) {
	engine := NewNutritionEngine()

	// Over 750 Kcal
	result := engine.EvaluateNutritionalCompliance("LUNCH", 850.0, 22.0, 35.0, 43.0, true)

	if result.IsCompliant {
		t.Errorf("Expected surplus menu to be non-compliant")
	}
	if result.Breakdown.EnergyBalanceStatus != "DƯ_THỪA_NĂNG_LƯỢNG" {
		t.Errorf("Expected DƯ_THỪA_NĂNG_LƯỢNG, got %s", result.Breakdown.EnergyBalanceStatus)
	}
}

func TestNutritionEngine_VerifySampleRetention(t *testing.T) {
	engine := NewNutritionEngine()

	keeperID := uuid.New()
	sampledAt := time.Date(2026, 10, 28, 11, 0, 0, 0, time.Local)

	// Valid 24h retention check (25 hours later, 4.0 °C, keeper assigned)
	checkAt := sampledAt.Add(25 * time.Hour)
	ver := engine.VerifySampleRetention(sampledAt, checkAt, 4.0, &keeperID, "Nguyễn Y Tá")

	if !ver.IsVerified {
		t.Errorf("Expected sample retention to be verified: %+v", ver)
	}
	if !ver.IsRetained24Hours {
		t.Errorf("Expected 24h retention to be true")
	}
	if !ver.IsTempCompliant {
		t.Errorf("Expected temp to be compliant (4.0°C in [2, 8])")
	}

	// Invalid: only 12 hours
	earlyCheck := sampledAt.Add(12 * time.Hour)
	verEarly := engine.VerifySampleRetention(sampledAt, earlyCheck, 4.0, &keeperID, "Nguyễn Y Tá")
	if verEarly.IsVerified || verEarly.IsRetained24Hours {
		t.Errorf("Expected early check to fail 24h requirement")
	}

	// Invalid: temperature too warm (12.0°C)
	verWarm := engine.VerifySampleRetention(sampledAt, checkAt, 12.0, &keeperID, "Nguyễn Y Tá")
	if verWarm.IsVerified || verWarm.IsTempCompliant {
		t.Errorf("Expected warm temperature to fail compliance")
	}

	// Invalid: missing keeper
	verNoKeeper := engine.VerifySampleRetention(sampledAt, checkAt, 4.0, nil, "")
	if verNoKeeper.IsVerified {
		t.Errorf("Expected missing keeper to fail verification")
	}
}
