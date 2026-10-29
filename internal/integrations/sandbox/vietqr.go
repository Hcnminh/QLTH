package sandbox

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// VietQRGenerateRequest represents the input parameters to generate a dynamic VietQR.
type VietQRGenerateRequest struct {
	BankBIN       string  `json:"bank_bin" binding:"required"`       // 6-digit BIN e.g. "970422" (MBBank), "970415" (Vietinbank)
	AccountNumber string  `json:"account_number" binding:"required"` // Bank account number
	AccountName   string  `json:"account_name,omitempty"`            // Account holder name
	Amount        float64 `json:"amount" binding:"required"`         // Amount in VND
	Description   string  `json:"description" binding:"required"`    // Payment description (e.g. "TIEN AN BAN TRU THANG 11")
	StudentID     string  `json:"student_id,omitempty"`              // Optional student ID reference
	BillID        string  `json:"bill_id,omitempty"`                 // Optional bill/invoice ID
}

// VietQRGenerateResult represents the EMVCo-compliant QR payload.
type VietQRGenerateResult struct {
	Environment string  `json:"environment"` // Always "SANDBOX_STUB"
	Mode        string  `json:"mode"`        // "SANDBOX_STUB"
	QRString    string  `json:"qr_string"`   // EMVCo TLV formatted string with valid CRC16
	ImageURL    string  `json:"image_url"`   // Direct image link for rendering
	BankBIN     string  `json:"bank_bin"`
	AccountNo   string  `json:"account_no"`
	AccountName string  `json:"account_name"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"` // "VND" (code 704)
	Description string  `json:"description"`
	CRC         string  `json:"crc"`
	GeneratedAt string  `json:"generated_at"`
	LegalNotice string  `json:"legal_notice"`
}

// VietQRService generates dynamic EMVCo-compliant VietQR codes in sandbox mode.
type VietQRService struct{}

// NewVietQRService initializes a VietQR sandbox generator.
func NewVietQRService() *VietQRService {
	return &VietQRService{}
}

var binRegex = regexp.MustCompile(`^\d{6}$`)

// Generate produces a dynamic EMVCo-compliant VietQR string with mathematical CRC-16 checksum
// and a standardized image URL.
func (s *VietQRService) Generate(ctx context.Context, req VietQRGenerateRequest) (*VietQRGenerateResult, error) {
	if !binRegex.MatchString(req.BankBIN) {
		return nil, fmt.Errorf("mã ngân hàng (Bank BIN) phải đúng 6 chữ số theo chuẩn NAPAS (nhận được: %s)", req.BankBIN)
	}
	if strings.TrimSpace(req.AccountNumber) == "" {
		return nil, fmt.Errorf("số tài khoản không được để trống")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("số tiền thanh toán phải lớn hơn 0")
	}

	accName := strings.ToUpper(strings.TrimSpace(req.AccountName))
	if accName == "" {
		accName = "TRUONG TIEU HOC EDUPRIMARY"
	}
	desc := strings.TrimSpace(req.Description)

	// Format EMVCo TLV string (NAPAS 247 / Circular 19/2016/TT-NHNN & Decision 1927/QD-NHNN)
	// Tag 00: Payload Format Indicator ("01")
	tag00 := "000201"
	// Tag 01: Point of Initiation Method ("12" = Dynamic QR with fixed amount)
	tag01 := "010212"

	// Tag 38: Merchant Account Information (NAPAS 247)
	// Subtag 00: AID ("A000000727")
	sub00 := "0010A000000727"
	// Subtag 01: Beneficiary info -> Tag 00 (BIN, len 6) + Tag 01 (Account number)
	binSub := fmt.Sprintf("0006%s", req.BankBIN)
	accSub := fmt.Sprintf("01%02d%s", len(req.AccountNumber), req.AccountNumber)
	beneficiaryContent := binSub + accSub
	sub01 := fmt.Sprintf("01%02d%s", len(beneficiaryContent), beneficiaryContent)
	// Subtag 02: Service code ("QRIBFTTA" - Transfer to Account)
	sub02 := "0208QRIBFTTA"
	tag38Content := sub00 + sub01 + sub02
	tag38 := fmt.Sprintf("38%02d%s", len(tag38Content), tag38Content)

	// Tag 53: Transaction Currency (VND = 704)
	tag53 := "5303704"

	// Tag 54: Transaction Amount
	amountStr := fmt.Sprintf("%.0f", req.Amount)
	tag54 := fmt.Sprintf("54%02d%s", len(amountStr), amountStr)

	// Tag 58: Country Code ("VN")
	tag58 := "5802VN"

	// Tag 62: Additional Data Field (Purpose of transaction)
	subTag08 := fmt.Sprintf("08%02d%s", len(desc), desc)
	tag62 := fmt.Sprintf("62%02d%s", len(subTag08), subTag08)

	// Assemble data before CRC
	rawBeforeCRC := tag00 + tag01 + tag38 + tag53 + tag54 + tag58 + tag62 + "6304"

	// Calculate CRC-16/CCITT-FALSE
	crcVal := calculateCRC16CCITT([]byte(rawBeforeCRC))
	crcHex := fmt.Sprintf("%04X", crcVal)
	fullQRString := rawBeforeCRC + crcHex

	// Build VietQR image URL
	imgURL := fmt.Sprintf("https://img.vietqr.io/image/%s-%s-compact2.png?amount=%s&addInfo=%s&accountName=%s",
		req.BankBIN,
		req.AccountNumber,
		amountStr,
		url.QueryEscape(desc),
		url.QueryEscape(accName),
	)

	return &VietQRGenerateResult{
		Environment: "SANDBOX_STUB",
		Mode:        "SANDBOX_STUB",
		QRString:    fullQRString,
		ImageURL:    imgURL,
		BankBIN:     req.BankBIN,
		AccountNo:   req.AccountNumber,
		AccountName: accName,
		Amount:      req.Amount,
		Currency:    "VND",
		Description: desc,
		CRC:         crcHex,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		LegalNotice: "Mã VietQR động tuân thủ chuẩn EMVCo và Quyết định 1927/QĐ-NHNN (NAPAS 247 Sandbox Stub)",
	}, nil
}

// calculateCRC16CCITT calculates the 16-bit CRC checksum using polynomial 0x1021 and initial value 0xFFFF
// as specified in EMVCo QR Code Specification for Payment Systems.
func calculateCRC16CCITT(data []byte) uint16 {
	crc := uint16(0xFFFF)
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if (crc & 0x8000) != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}
