package sandbox

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// InvoiceItem represents a single line item in an electronic invoice.
type InvoiceItem struct {
	Name      string  `json:"name" binding:"required"`
	Unit      string  `json:"unit"`
	Quantity  float64 `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Amount    float64 `json:"amount"`
	VATRate   float64 `json:"vat_rate"` // e.g. 0.0 (non-taxable for primary education)
	VATAmount float64 `json:"vat_amount"`
}

// EInvoiceIssueRequest represents the request to create and issue an e-invoice.
type EInvoiceIssueRequest struct {
	SchoolID      string        `json:"school_id" binding:"required"`
	BuyerName     string        `json:"buyer_name" binding:"required"` // Guardian / Parent name
	BuyerTaxCode  string        `json:"buyer_tax_code,omitempty"`
	BuyerPhone    string        `json:"buyer_phone,omitempty"`
	BuyerEmail    string        `json:"buyer_email,omitempty"`
	StudentID     string        `json:"student_id,omitempty"`
	StudentName   string        `json:"student_name,omitempty"`
	StudentClass  string        `json:"student_class,omitempty"`
	InvoiceType   string        `json:"invoice_type,omitempty"`   // "BOARDING_FEE", "TUITION", "FACILITY"
	PaymentMethod string        `json:"payment_method,omitempty"` // "CHUYEN_KHOAN_VIETQR", "TIEN_MAT"
	Items         []InvoiceItem `json:"items" binding:"required"`
}

// EInvoiceIssueResult represents the simulated Decree 123 e-invoice.
type EInvoiceIssueResult struct {
	Environment      string        `json:"environment"` // Always "SANDBOX_STUB"
	Mode             string        `json:"mode"`        // "SANDBOX_STUB"
	InvoiceID        string        `json:"invoice_id"`
	InvoiceSeries    string        `json:"invoice_series"`     // "1C26TAA"
	InvoiceNumber    string        `json:"invoice_number"`     // 8 digits
	TaxAuthorityCode string        `json:"tax_authority_code"` // 23 characters
	LookupCode       string        `json:"lookup_code"`
	IssueDate        string        `json:"issue_date"`
	SellerName       string        `json:"seller_name"`
	SellerTaxCode    string        `json:"seller_tax_code"`
	SellerAddress    string        `json:"seller_address"`
	BuyerName        string        `json:"buyer_name"`
	BuyerTaxCode     string        `json:"buyer_tax_code,omitempty"`
	StudentName      string        `json:"student_name,omitempty"`
	StudentClass     string        `json:"student_class,omitempty"`
	PaymentMethod    string        `json:"payment_method"`
	Items            []InvoiceItem `json:"items"`
	Subtotal         float64       `json:"subtotal"`
	VATAmount        float64       `json:"vat_amount"`
	TotalAmount      float64       `json:"total_amount"`
	AmountInWords    string        `json:"amount_in_words"`
	SignedBy         string        `json:"signed_by"`
	SignedAt         string        `json:"signed_at"`
	Status           string        `json:"status"` // "ISSUED"
	PDFURL           string        `json:"pdf_url"`
	XMLURL           string        `json:"xml_url"`
	LegalNotice      string        `json:"legal_notice"`
}

// EInvoiceService handles simulation of electronic invoices compliant with Decree 123/2020/ND-CP.
type EInvoiceService struct{}

// NewEInvoiceService initializes a new e-invoice simulation service.
func NewEInvoiceService() *EInvoiceService {
	return &EInvoiceService{}
}

// IssueInvoice validates items and creates a simulated electronic invoice.
func (s *EInvoiceService) IssueInvoice(ctx context.Context, req EInvoiceIssueRequest) (*EInvoiceIssueResult, error) {
	if strings.TrimSpace(req.SchoolID) == "" {
		return nil, fmt.Errorf("school_id không được để trống")
	}
	if strings.TrimSpace(req.BuyerName) == "" {
		return nil, fmt.Errorf("tên người mua (buyer_name) không được để trống")
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("danh sách dịch vụ/mặt hàng (items) không được rỗng")
	}

	var subtotal float64
	var totalVAT float64
	var processedItems []InvoiceItem

	for _, item := range req.Items {
		itemAmt := item.Quantity * item.UnitPrice
		vatAmt := itemAmt * item.VATRate
		processedItems = append(processedItems, InvoiceItem{
			Name:      item.Name,
			Unit:      item.Unit,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Amount:    itemAmt,
			VATRate:   item.VATRate,
			VATAmount: vatAmt,
		})
		subtotal += itemAmt
		totalVAT += vatAmt
	}

	totalAmount := subtotal + totalVAT
	now := time.Now().UTC()

	// Decree 123/2020/ND-CP compliance format:
	// Ký hiệu mẫu số: 1 (Hóa đơn GTGT)
	// Ký hiệu hóa đơn: C (Có mã CQT), 26 (Năm 2026), T (Điện tử), AA (Dải ký hiệu) -> "1C26TAA"
	series := "1C26TAA"
	invNum := fmt.Sprintf("%08d", (now.UnixNano()%100000000)+1)
	cqtCode := fmt.Sprintf("T26%010X%010X", now.Unix(), now.UnixNano()%10000000000)
	lookupCode := fmt.Sprintf("INV-%04X-%04X-%04X", now.Unix()%65536, (now.Unix()/65536)%65536, now.Nanosecond()%65536)

	paymentMethod := req.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "CHUYEN_KHOAN_VIETQR"
	}

	invoiceID := uuid.New().String()

	return &EInvoiceIssueResult{
		Environment:      "SANDBOX_STUB",
		Mode:             "SANDBOX_STUB",
		InvoiceID:        invoiceID,
		InvoiceSeries:    series,
		InvoiceNumber:    invNum,
		TaxAuthorityCode: cqtCode,
		LookupCode:       lookupCode,
		IssueDate:        now.Format("2006-01-02"),
		SellerName:       "Trường Tiểu học EduPrimary",
		SellerTaxCode:    "0102030405",
		SellerAddress:    "Phố Tràng Thi, Phường Hàng Trống, Quận Hoàn Kiếm, TP. Hà Nội",
		BuyerName:        req.BuyerName,
		BuyerTaxCode:     req.BuyerTaxCode,
		StudentName:      req.StudentName,
		StudentClass:     req.StudentClass,
		PaymentMethod:    paymentMethod,
		Items:            processedItems,
		Subtotal:         subtotal,
		VATAmount:        totalVAT,
		TotalAmount:      totalAmount,
		AmountInWords:    formatVietnameseCurrencyInWords(totalAmount),
		SignedBy:         "Chữ ký số tập trung Trường Tiểu học EduPrimary (SmartCA VNPT)",
		SignedAt:         now.Format(time.RFC3339),
		Status:           "ISSUED",
		PDFURL:           fmt.Sprintf("https://einvoice.sandbox.eduprimary.vn/view/%s.pdf", lookupCode),
		XMLURL:           fmt.Sprintf("https://einvoice.sandbox.eduprimary.vn/xml/%s.xml", lookupCode),
		LegalNotice:      "Hóa đơn điện tử có mã của cơ quan thuế mô phỏng theo Nghị định 123/2020/NĐ-CP và Thông tư 78/2021/TT-BTC",
	}, nil
}

// formatVietnameseCurrencyInWords converts common integer currency amounts to Vietnamese text.
func formatVietnameseCurrencyInWords(amount float64) string {
	n := int64(amount)
	if n <= 0 {
		return "Không đồng chẵn"
	}

	units := []string{"", "nghìn", "triệu", "tỷ"}
	var parts []string
	idx := 0

	for n > 0 && idx < len(units) {
		rem := n % 1000
		if rem > 0 {
			partStr := fmt.Sprintf("%d %s", rem, units[idx])
			parts = append([]string{strings.TrimSpace(partStr)}, parts...)
		}
		n /= 1000
		idx++
	}

	result := strings.Join(parts, " ") + " đồng chẵn"
	return strings.ToUpper(result[:1]) + result[1:]
}
