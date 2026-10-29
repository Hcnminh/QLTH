package handlers

import (
	"net/http"
	"strings"

	"qlth/internal/integrations/sandbox"

	"github.com/gin-gonic/gin"
)

// SignSmartCA handles POST /api/v1/integrations/smartca/sign
// Simulates digital signing with SmartCA cloud provider and RFC 3161 TSA timestamp.
func (h *APIHandler) SignSmartCA(c *gin.Context) {
	var req sandbox.SmartCASignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu ký số SmartCA không hợp lệ: " + err.Error(),
		})
		return
	}

	result, err := h.smartCASandbox.SignDocument(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Lỗi ký số SmartCA: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Ký số SmartCA và cấp dấu thời gian TSA thành công (Sandbox Stub)",
		"data":    result,
	})
}

// GenerateVietQR handles POST /api/v1/integrations/vietqr/generate
// Generates dynamic EMVCo-compliant QR string and image URL according to NAPAS 247.
func (h *APIHandler) GenerateVietQR(c *gin.Context) {
	var req sandbox.VietQRGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu sinh VietQR không hợp lệ: " + err.Error(),
		})
		return
	}

	result, err := h.vietQRSandbox.Generate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Lỗi tạo mã VietQR: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Sinh mã VietQR động chuẩn EMVCo thành công (Sandbox Stub)",
		"data":    result,
	})
}

// SyncMOET handles POST /api/v1/integrations/moet/sync
// Processes the outbound queue of student profiles and transcripts to MOET CSDL.
func (h *APIHandler) SyncMOET(c *gin.Context) {
	var req sandbox.MOETSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu đồng bộ MOET CSDL không hợp lệ: " + err.Error(),
		})
		return
	}

	if req.SchoolID == "" {
		req.SchoolID = c.GetHeader("X-School-ID")
	}

	result, err := h.moetSandbox.ProcessOutboxSync(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Lỗi xử lý hàng đợi đồng bộ MOET: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Xử lý hàng đợi Outbox đồng bộ CSDL Ngành Giáo dục MOET thành công (Sandbox Stub)",
		"data":    result,
	})
}

// GetMOETQueue handles GET /api/v1/integrations/moet/queue
// Inspects current status and pending counts in the outbound MOET queue.
func (h *APIHandler) GetMOETQueue(c *gin.Context) {
	schoolID := c.Query("school_id")
	if schoolID == "" {
		schoolID = c.GetHeader("X-School-ID")
	}
	if schoolID == "" {
		schoolID = "00000000-0000-0000-0000-000000000001"
	}

	status, err := h.moetSandbox.GetQueueStatus(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    500,
			"message": "Lỗi kiểm tra hàng đợi MOET: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Lấy trạng thái hàng đợi Outbox MOET thành công",
		"data":    status,
	})
}

// VerifyVNeID handles POST /api/v1/integrations/vneid/verify
// Simulates citizen Level-2 authentication verification according to Decree 69/2024 & Decree 13/2023.
func (h *APIHandler) VerifyVNeID(c *gin.Context) {
	var req sandbox.VNeIDVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu định danh VNeID không hợp lệ: " + err.Error(),
		})
		return
	}

	result, err := h.vneidSandbox.VerifyLevel2(c.Request.Context(), req)
	if err != nil {
		code := http.StatusBadRequest
		if strings.Contains(err.Error(), "Nghị định 13/2023") {
			code = http.StatusUnprocessableEntity // 422 for consent missing
		}
		c.JSON(code, gin.H{
			"success": false,
			"code":    code,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Xác thực định danh điện tử VNeID Mức 2 thành công (Sandbox Stub)",
		"data":    result,
	})
}

// IssueEInvoice handles POST /api/v1/integrations/einvoice/issue
// Simulates electronic invoice issuance compliant with Decree 123/2020/ND-CP.
func (h *APIHandler) IssueEInvoice(c *gin.Context) {
	headerSchoolID := c.GetHeader("X-School-ID")
	var req sandbox.EInvoiceIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Dữ liệu phát hành hóa đơn không hợp lệ: " + err.Error(),
		})
		return
	}

	if req.SchoolID == "" {
		req.SchoolID = headerSchoolID
	}
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Danh sách mặt hàng hóa đơn không được rỗng",
		})
		return
	}

	invResult, invErr := h.einvoiceSandbox.IssueInvoice(c.Request.Context(), req)
	if invErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"message": "Lỗi phát hành hóa đơn điện tử: " + invErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Phát hành hóa đơn điện tử theo Nghị định 123/2020/NĐ-CP thành công (Sandbox Stub)",
		"data":    invResult,
	})
}

// GetSandboxStatus handles GET /api/v1/integrations/sandbox/status
// Returns the status and specifications of all 5 external integration sandbox stubs.
func (h *APIHandler) GetSandboxStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "Trạng thái các cổng kết nối ngoại vi môi trường Sandbox Stub",
		"data": gin.H{
			"environment": "SANDBOX_STUB",
			"integrations": []gin.H{
				{
					"subsystem":   "Phân hệ 9 & 10: Định danh điện tử công dân VNeID",
					"service":     "VNeID Level 2 Authentication",
					"status":      "ACTIVE_SANDBOX",
					"legal_basis": "Nghị định 69/2024/NĐ-CP & Nghị định 13/2023/NĐ-CP",
					"endpoint":    "POST /api/v1/integrations/vneid/verify",
				},
				{
					"subsystem":   "Phân hệ 9: Ký số điện tử & Dấu thời gian TSA",
					"service":     "SmartCA Digital Signing with TSA",
					"status":      "ACTIVE_SANDBOX",
					"legal_basis": "Luật Giao dịch điện tử 2023 & Thông tư 03/2018/TT-BGDĐT",
					"endpoint":    "POST /api/v1/integrations/smartca/sign",
				},
				{
					"subsystem":   "Phân hệ 10: Đồng bộ CSDL Ngành Giáo dục",
					"service":     "MOET Sector Outbox Queue Sync",
					"status":      "ACTIVE_SANDBOX",
					"legal_basis": "Thông tư 42/2021/TT-BGDĐT",
					"endpoint":    "POST /api/v1/integrations/moet/sync",
				},
				{
					"subsystem":   "Phân hệ 7: Tài chính & Cổng thanh toán bán trú",
					"service":     "VietQR Dynamic EMVCo Payment Generation",
					"status":      "ACTIVE_SANDBOX",
					"legal_basis": "Quyết định 1927/QĐ-NHNN & Chuẩn EMVCo NAPAS 247",
					"endpoint":    "POST /api/v1/integrations/vietqr/generate",
				},
				{
					"subsystem":   "Phân hệ 7: Tài chính & Thu phí học đường",
					"service":     "E-Invoice Simulation",
					"status":      "ACTIVE_SANDBOX",
					"legal_basis": "Nghị định 123/2020/NĐ-CP & Thông tư 78/2021/TT-BTC",
					"endpoint":    "POST /api/v1/integrations/einvoice/issue",
				},
			},
		},
	})
}
