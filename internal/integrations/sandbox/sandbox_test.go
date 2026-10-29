package sandbox

import (
	"context"
	"strings"
	"testing"
)

func TestVNeID_VerifyLevel2_Success(t *testing.T) {
	service := NewVNeIDService()
	req := VNeIDVerificationRequest{
		NationalID:     "001202012345",
		FullName:       "Nguyễn Văn An",
		DateOfBirth:    "2018-05-15",
		Gender:         "NAM",
		ConsentGranted: true,
	}

	result, err := service.VerifyLevel2(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Environment != "SANDBOX_STUB" || result.Mode != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB environment marker, got %s / %s", result.Environment, result.Mode)
	}
	if result.VerificationLevel != "LEVEL_2" {
		t.Errorf("expected LEVEL_2 verification, got: %s", result.VerificationLevel)
	}
	if !result.BiometricMatched {
		t.Errorf("expected BiometricMatched to be true")
	}
	if result.Status != "VERIFIED" {
		t.Errorf("expected status VERIFIED, got: %s", result.Status)
	}
	if result.NationalID != req.NationalID {
		t.Errorf("expected NationalID %s, got %s", req.NationalID, result.NationalID)
	}
}

func TestVNeID_VerifyLevel2_InvalidCCCD(t *testing.T) {
	service := NewVNeIDService()
	req := VNeIDVerificationRequest{
		NationalID:     "12345", // Not 12 digits
		FullName:       "Nguyễn Văn An",
		ConsentGranted: true,
	}

	_, err := service.VerifyLevel2(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid 12-digit CCCD, got nil")
	}
	if !strings.Contains(err.Error(), "12 chữ số") {
		t.Errorf("expected error mentioning 12 digits, got: %v", err)
	}
}

func TestVNeID_VerifyLevel2_MissingConsent(t *testing.T) {
	service := NewVNeIDService()
	req := VNeIDVerificationRequest{
		NationalID:     "001202012345",
		FullName:       "Nguyễn Văn An",
		ConsentGranted: false, // Violates Decree 13/2023
	}

	_, err := service.VerifyLevel2(context.Background(), req)
	if err == nil {
		t.Fatal("expected error when digital consent is false, got nil")
	}
	if !strings.Contains(err.Error(), "Nghị định 13/2023/NĐ-CP") {
		t.Errorf("expected error citing Decree 13/2023, got: %v", err)
	}
}

func TestSmartCA_SignDocument_Success(t *testing.T) {
	service := NewSmartCAService()
	req := SmartCASignRequest{
		DocumentID:   "d1000000-0000-0000-0000-000000000001",
		DocumentType: "TRANSCRIPT_SNAPSHOT",
		SignerID:     "b0000000-0000-0000-0000-000000000001",
		SignerName:   "Trần Thị Bình",
		SignerRole:   "GVCN",
	}

	result, err := service.SignDocument(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Environment != "SANDBOX_STUB" {
		t.Errorf("expected environment SANDBOX_STUB, got %s", result.Environment)
	}
	if result.SignatureFormat != "CMS/PKCS#7" {
		t.Errorf("expected CMS/PKCS#7, got %s", result.SignatureFormat)
	}
	if !strings.HasPrefix(result.CMSPKCS7Signature, "MIAGCSqGSIb3DQEHAqCAMIACAQEx") {
		t.Errorf("expected CMS/PKCS#7 standard header in signature")
	}
	if result.TSAToken.TSAProvider == "" || result.TSAToken.TokenBase64 == "" {
		t.Errorf("expected valid TSA token structure")
	}
	if result.Certificate.SerialNumber == "" {
		t.Errorf("expected certificate serial number")
	}
}

func TestSmartCA_SignDocument_MissingFields(t *testing.T) {
	service := NewSmartCAService()
	req := SmartCASignRequest{
		DocumentID: "", // Missing
		SignerID:   "b1",
		SignerName: "Teacher",
		SignerRole: "GVCN",
	}

	_, err := service.SignDocument(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty document_id, got nil")
	}
}

func TestVietQR_Generate_Success(t *testing.T) {
	service := NewVietQRService()
	req := VietQRGenerateRequest{
		BankBIN:       "970422",
		AccountNumber: "0123456789",
		AccountName:   "Truong Tieu Hoc EduPrimary",
		Amount:        550000,
		Description:   "TIEN AN BAN TRU THANG 11",
	}

	result, err := service.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Environment != "SANDBOX_STUB" {
		t.Errorf("expected environment SANDBOX_STUB, got: %s", result.Environment)
	}
	if !strings.HasPrefix(result.QRString, "000201010212") {
		t.Errorf("expected dynamic EMVCo prefix 000201010212, got: %s", result.QRString)
	}
	if !strings.Contains(result.QRString, "6304") {
		t.Errorf("expected CRC16 Tag 6304 in QR string")
	}
	if result.CRC == "" || len(result.CRC) != 4 {
		t.Errorf("expected 4-character hex CRC, got: %s", result.CRC)
	}
	if !strings.Contains(result.ImageURL, "https://img.vietqr.io/image/970422-0123456789-compact2.png") {
		t.Errorf("expected standard VietQR image URL, got: %s", result.ImageURL)
	}
}

func TestVietQR_Generate_InvalidBIN(t *testing.T) {
	service := NewVietQRService()
	req := VietQRGenerateRequest{
		BankBIN:       "9704", // Not 6 digits
		AccountNumber: "123456",
		Amount:        100000,
		Description:   "Test",
	}

	_, err := service.Generate(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid Bank BIN, got nil")
	}
}

func TestVietQR_CRC16Calculation(t *testing.T) {
	// Standard EMVCo test vector or check that calculateCRC16CCITT returns non-zero
	data := []byte("00020101021238540010A00000072701200006970422011001234567890208QRIBFTTA530370454065500005802VN62280824TIEN AN BAN TRU THANG 116304")
	crc := calculateCRC16CCITT(data)
	if crc == 0 {
		t.Errorf("calculated CRC should not be 0 for this vector")
	}
}

func TestMOETSync_ProcessOutboxSync_Success(t *testing.T) {
	service := NewMOETSyncService()
	req := MOETSyncRequest{
		SchoolID:     "00000000-0000-0000-0000-000000000001",
		AcademicYear: "2026-2027",
		SyncType:     "ALL",
		BatchSize:    10,
	}

	result, err := service.ProcessOutboxSync(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Environment != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB, got: %s", result.Environment)
	}
	if result.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got: %s", result.Status)
	}
	if result.TotalProcessed == 0 {
		t.Errorf("expected processed items count > 0")
	}
	if result.MOETTransactionID == "" {
		t.Errorf("expected valid MOETTransactionID")
	}
}

func TestMOETSync_GetQueueStatus(t *testing.T) {
	service := NewMOETSyncService()
	status, err := service.GetQueueStatus(context.Background(), "00000000-0000-0000-0000-000000000001")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if status.Environment != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB, got: %s", status.Environment)
	}
	if status.Status != "HEALTHY" {
		t.Errorf("expected HEALTHY status, got: %s", status.Status)
	}
}

func TestEInvoice_IssueInvoice_Success(t *testing.T) {
	service := NewEInvoiceService()
	req := EInvoiceIssueRequest{
		SchoolID:     "00000000-0000-0000-0000-000000000001",
		BuyerName:    "Nguyễn Thị Hoa",
		StudentName:  "Nguyễn Hoàng Quân",
		StudentClass: "5A1",
		Items: []InvoiceItem{
			{
				Name:      "Tiền ăn bán trú tháng 11/2026",
				Unit:      "Tháng",
				Quantity:  1,
				UnitPrice: 660000,
				VATRate:   0,
			},
			{
				Name:      "Phí dịch vụ bán trú",
				Unit:      "Tháng",
				Quantity:  1,
				UnitPrice: 150000,
				VATRate:   0,
			},
		},
	}

	result, err := service.IssueInvoice(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Environment != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB, got: %s", result.Environment)
	}
	if result.InvoiceSeries != "1C26TAA" {
		t.Errorf("expected Decree 123 series 1C26TAA, got: %s", result.InvoiceSeries)
	}
	if result.TotalAmount != 810000 {
		t.Errorf("expected total 810,000 VND, got: %f", result.TotalAmount)
	}
	if !strings.Contains(result.AmountInWords, "đồng chẵn") {
		t.Errorf("expected currency in words, got: %s", result.AmountInWords)
	}
	if result.TaxAuthorityCode == "" || result.LookupCode == "" {
		t.Errorf("expected TaxAuthorityCode and LookupCode")
	}
}

func TestEInvoice_IssueInvoice_EmptyItems(t *testing.T) {
	service := NewEInvoiceService()
	req := EInvoiceIssueRequest{
		SchoolID:  "00000000-0000-0000-0000-000000000001",
		BuyerName: "Nguyễn Thị Hoa",
		Items:     nil, // Empty items
	}

	_, err := service.IssueInvoice(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty items, got nil")
	}
}
