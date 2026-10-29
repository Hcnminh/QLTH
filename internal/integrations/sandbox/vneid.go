package sandbox

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// VNeIDVerificationRequest represents a citizen Level-2 verification request
// compliant with Decree 69/2024/ND-CP and Decree 13/2023/ND-CP.
type VNeIDVerificationRequest struct {
	NationalID     string `json:"national_id" binding:"required"` // 12-digit citizen ID (CCCD)
	FullName       string `json:"full_name" binding:"required"`
	DateOfBirth    string `json:"date_of_birth,omitempty"` // YYYY-MM-DD
	Gender         string `json:"gender,omitempty"`        // MALE, FEMALE
	DeviceID       string `json:"device_id,omitempty"`
	ConsentGranted bool   `json:"consent_granted"` // Decree 13/2023 personal data consent
}

// VNeIDVerificationResult represents the realistic mock Level 2 payload
// marked clearly as SANDBOX_STUB.
type VNeIDVerificationResult struct {
	Environment        string  `json:"environment"` // Always "SANDBOX_STUB"
	Mode               string  `json:"mode"`        // "SANDBOX_STUB"
	VerificationID     string  `json:"verification_id"`
	NationalID         string  `json:"national_id"`
	FullName           string  `json:"full_name"`
	DateOfBirth        string  `json:"date_of_birth"`
	Gender             string  `json:"gender"`
	VerificationLevel  string  `json:"verification_level"` // "LEVEL_2"
	Status             string  `json:"status"`             // "VERIFIED"
	BiometricMatched   bool    `json:"biometric_matched"`  // true
	FaceMatchScore     float64 `json:"face_match_score"`   // e.g. 99.4
	IdentityCardNumber string  `json:"identity_card_number"`
	IssuedDate         string  `json:"issued_date"`
	IssuedPlace        string  `json:"issued_place"`
	PermanentAddress   string  `json:"permanent_address"`
	Hometown           string  `json:"hometown"`
	AccessToken        string  `json:"access_token"`
	IssuedAt           string  `json:"issued_at"`
	ExpiresAt          string  `json:"expires_at"`
	LegalNotice        string  `json:"legal_notice"`
}

// VNeIDService provides Level 2 identity verification in sandbox mode.
type VNeIDService struct{}

// NewVNeIDService initializes a new VNeID sandbox service.
func NewVNeIDService() *VNeIDService {
	return &VNeIDService{}
}

var cccdRegex = regexp.MustCompile(`^\d{12}$`)

// VerifyLevel2 simulates authenticating citizen information against the
// Ministry of Public Security (C06) Level-2 electronic identity database.
func (s *VNeIDService) VerifyLevel2(ctx context.Context, req VNeIDVerificationRequest) (*VNeIDVerificationResult, error) {
	// 1. Validate 12-digit CCCD format (Nghị định 69/2024/NĐ-CP)
	if !cccdRegex.MatchString(req.NationalID) {
		return nil, fmt.Errorf("số định danh cá nhân VNeID phải đúng 12 chữ số theo quy định tại Nghị định 69/2024/NĐ-CP (nhận được: %s)", req.NationalID)
	}

	// 2. Validate Digital Consent (Điều 20 Nghị định 13/2023/NĐ-CP)
	if !req.ConsentGranted {
		return nil, fmt.Errorf("không có sự đồng thuận số (Digital Consent) theo Điều 20 Nghị định 13/2023/NĐ-CP về bảo vệ dữ liệu cá nhân")
	}

	now := time.Now()
	dob := req.DateOfBirth
	if dob == "" {
		dob = "2018-09-02"
	}
	gender := strings.ToUpper(req.Gender)
	if gender == "" || (gender != "MALE" && gender != "FEMALE" && gender != "NAM" && gender != "NỮ") {
		gender = "NAM"
	}

	result := &VNeIDVerificationResult{
		Environment:        "SANDBOX_STUB",
		Mode:               "SANDBOX_STUB",
		VerificationID:     uuid.New().String(),
		NationalID:         req.NationalID,
		FullName:           strings.ToUpper(strings.TrimSpace(req.FullName)),
		DateOfBirth:        dob,
		Gender:             gender,
		VerificationLevel:  "LEVEL_2",
		Status:             "VERIFIED",
		BiometricMatched:   true,
		FaceMatchScore:     99.4,
		IdentityCardNumber: req.NationalID,
		IssuedDate:         "2023-01-10",
		IssuedPlace:        "Cục Cảnh sát Quản lý hành chính về trật tự xã hội (C06)",
		PermanentAddress:   "Số 45 Phố Tràng Thi, Phường Hàng Trống, Quận Hoàn Kiếm, Thành phố Hà Nội",
		Hometown:           "Hà Nội",
		AccessToken:        fmt.Sprintf("vneid_lvl2_jwt_sandbox_%s_%d", uuid.New().String()[:8], now.Unix()),
		IssuedAt:           now.Format(time.RFC3339),
		ExpiresAt:          now.Add(2 * time.Hour).Format(time.RFC3339),
		LegalNotice:        "Mô phỏng xác thực định danh điện tử VNeID Mức độ 2 môi trường Sandbox Stub (Bộ Công an C06 / NĐ 69/2024/NĐ-CP & NĐ 13/2023/NĐ-CP)",
	}

	return result, nil
}
