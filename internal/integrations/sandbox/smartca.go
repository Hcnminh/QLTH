package sandbox

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"qlth/internal/db"

	"github.com/google/uuid"
)

// SmartCASignRequest represents a digital signing request to SmartCA.
type SmartCASignRequest struct {
	DocumentID   string `json:"document_id" binding:"required"`
	DocumentType string `json:"document_type,omitempty"` // "TRANSCRIPT_SNAPSHOT", "LESSON_PLAN", "DIPLOMA"
	DocumentHash string `json:"document_hash,omitempty"` // SHA256 hex if provided
	SignerID     string `json:"signer_id" binding:"required"`
	SignerName   string `json:"signer_name" binding:"required"`
	SignerRole   string `json:"signer_role" binding:"required"` // "GVCN", "BGH", "TO_TRUONG"
	PinCode      string `json:"pin_code,omitempty"`
}

// TSATimestamp holds RFC 3161 compliant timestamp authority token metadata.
type TSATimestamp struct {
	TSAProvider   string `json:"tsa_provider"`
	SerialNumber  string `json:"serial_number"`
	Timestamp     string `json:"timestamp"` // RFC 3339
	PolicyOID     string `json:"policy_oid"`
	TokenBase64   string `json:"token_base64"`
	HashAlgorithm string `json:"hash_algorithm"`
}

// SignerCertificate holds the digital certificate information of the signer.
type SignerCertificate struct {
	SubjectDN    string `json:"subject_dn"`
	IssuerDN     string `json:"issuer_dn"`
	SerialNumber string `json:"serial_number"`
	ValidFrom    string `json:"valid_from"`
	ValidTo      string `json:"valid_to"`
}

// SmartCASignResult represents the response with simulated CMS/PKCS#7 and TSA token.
type SmartCASignResult struct {
	Environment       string            `json:"environment"` // Always "SANDBOX_STUB"
	Mode              string            `json:"mode"`        // "SANDBOX_STUB"
	SignatureID       string            `json:"signature_id"`
	DocumentID        string            `json:"document_id"`
	DocumentType      string            `json:"document_type"`
	DocumentHash      string            `json:"document_hash"`
	SignatureFormat   string            `json:"signature_format"`    // "CMS/PKCS#7"
	CMSPKCS7Signature string            `json:"cms_pkcs7_signature"` // Base64 CMS/PKCS#7 structure
	TSAToken          TSATimestamp      `json:"tsa_token"`
	Certificate       SignerCertificate `json:"certificate"`
	Status            string            `json:"status"` // "SIGNED"
	SignedAt          string            `json:"signed_at"`
	LegalNotice       string            `json:"legal_notice"`
}

// SmartCAService provides cloud digital signing and TSA timestamping stubs.
type SmartCAService struct{}

// NewSmartCAService initializes a SmartCA sandbox service.
func NewSmartCAService() *SmartCAService {
	return &SmartCAService{}
}

// SignDocument generates simulated CMS/PKCS#7 signature and TSA token,
// updating database transcript snapshots if applicable.
func (s *SmartCAService) SignDocument(ctx context.Context, req SmartCASignRequest) (*SmartCASignResult, error) {
	if strings.TrimSpace(req.DocumentID) == "" {
		return nil, fmt.Errorf("document_id không được để trống")
	}
	if strings.TrimSpace(req.SignerID) == "" || strings.TrimSpace(req.SignerName) == "" {
		return nil, fmt.Errorf("thông tin người ký (signer_id, signer_name) không được để trống")
	}

	docType := strings.ToUpper(strings.TrimSpace(req.DocumentType))
	if docType == "" {
		docType = "TRANSCRIPT_SNAPSHOT"
	}

	// Calculate or verify doc hash
	docHash := req.DocumentHash
	if docHash == "" {
		h := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", req.DocumentID, req.SignerID, time.Now().UTC().Format(time.RFC3339Nano))))
		docHash = hex.EncodeToString(h[:])
	}

	now := time.Now().UTC()
	sigID := uuid.New().String()

	// Simulate CMS/PKCS#7 ASN.1 DER payload base64
	// Header RFC 5652 CMS SignedData structure
	cmsPayload := fmt.Sprintf("CMS_PKCS7_SIGNED_DATA_DOC:%s|SIGNER:%s(%s)|HASH:%s|TSA:%s",
		req.DocumentID, req.SignerName, req.SignerRole, docHash, now.Format(time.RFC3339))
	cmsBase64 := base64.StdEncoding.EncodeToString([]byte(cmsPayload))
	// Prefix with standard CMS/PKCS#7 header signature simulation
	cmsFullSignature := "MIAGCSqGSIb3DQEHAqCAMIACAQExDzANBglghkgBZQMEAgEFADCABgkqhkiG9w0BBwEw" + cmsBase64

	// RFC 3161 TSA Token
	tsaSerial := fmt.Sprintf("TSA-VNPT-2026-%012d", now.UnixNano()%1000000000000)
	tsaPayload := fmt.Sprintf("RFC3161_TSA_TOKEN:SERIAL=%s|HASH=%s|TIME=%s|POLICY=1.3.6.1.4.1.601.10.3.1",
		tsaSerial, docHash, now.Format(time.RFC3339))
	tsaBase64 := base64.StdEncoding.EncodeToString([]byte(tsaPayload))

	tsaToken := TSATimestamp{
		TSAProvider:   "VNPT-TSA National Time-Stamp Authority (Sandbox)",
		SerialNumber:  tsaSerial,
		Timestamp:     now.Format(time.RFC3339),
		PolicyOID:     "1.3.6.1.4.1.601.10.3.1", // Standard Vietnam National TSA Policy OID
		TokenBase64:   tsaBase64,
		HashAlgorithm: "SHA-256",
	}

	certSerial := fmt.Sprintf("VNPT:CA:%08X:%04X", now.Unix(), now.Nanosecond()%65536)
	cert := SignerCertificate{
		SubjectDN:    fmt.Sprintf("CN=%s, TITLE=%s, O=Truong Tieu Hoc EduPrimary, C=VN", req.SignerName, req.SignerRole),
		IssuerDN:     "CN=VNPT SmartCA Cloud Signature Root CA, O=VNPT Group, C=VN",
		SerialNumber: certSerial,
		ValidFrom:    "2026-01-01T00:00:00Z",
		ValidTo:      "2029-01-01T23:59:59Z",
	}

	// Update database record if it's a transcript snapshot
	database, err := db.GetDB()
	if err == nil && database != nil {
		if docUUID, parseErr := uuid.Parse(req.DocumentID); parseErr == nil {
			role := strings.ToUpper(req.SignerRole)
			if role == "BGH" || role == "HIEU_TRUONG" || role == "PHO_HIEU_TRUONG" {
				_, _ = database.ExecContext(ctx, `
					UPDATE academic_transcripts_snapshots
					SET principal_signature = $1, principal_signed_at = $2, timestamp_token = $3
					WHERE id = $4
				`, cmsFullSignature, now, tsaBase64, docUUID)
			} else {
				_, _ = database.ExecContext(ctx, `
					UPDATE academic_transcripts_snapshots
					SET teacher_signature = $1, teacher_signed_at = $2, timestamp_token = $3
					WHERE id = $4
				`, cmsFullSignature, now, tsaBase64, docUUID)
			}
		}
	}

	return &SmartCASignResult{
		Environment:       "SANDBOX_STUB",
		Mode:              "SANDBOX_STUB",
		SignatureID:       sigID,
		DocumentID:        req.DocumentID,
		DocumentType:      docType,
		DocumentHash:      docHash,
		SignatureFormat:   "CMS/PKCS#7",
		CMSPKCS7Signature: cmsFullSignature,
		TSAToken:          tsaToken,
		Certificate:       cert,
		Status:            "SIGNED",
		SignedAt:          now.Format(time.RFC3339),
		LegalNotice:       "Chữ ký số SmartCA kèm dấu thời gian TSA mô phỏng theo Luật Giao dịch điện tử 2023 và Thông tư 03/2018/TT-BGDĐT",
	}, nil
}
