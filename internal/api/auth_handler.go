package api

import (
	"app/myproj/internal/auth"
	"app/myproj/pkg/workflow"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	workflowEngine *workflow.WorkflowEngine
	authService    *auth.AuthService
}

func NewAuthHandler(engine *workflow.WorkflowEngine, authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		workflowEngine: engine,
		authService:    authService,
	}
}

func (h *AuthHandler) StartAuth(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	instance, err := h.workflowEngine.CreateInstance("auth_workflow", map[string]interface{}{
		"phone_number": request.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.workflowEngine.HandleEvent(instance.ID, "submit_phone", map[string]interface{}{
		"phone_number": request.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"instance_id": instance.ID})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
		InstanceID  uint   `json:"instance_id"`
		OTP         string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.workflowEngine.HandleEvent(request.InstanceID, "verify_otp", map[string]interface{}{
		"phone_number":  request.PhoneNumber,
		"submitted_otp": request.OTP,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified"})
}

func (h *AuthHandler) SetPIN(c *gin.Context) {
	var request struct {
		InstanceID uint   `json:"instance_id"`
		PIN        string `json:"pin"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.workflowEngine.HandleEvent(request.InstanceID, "set_pin", map[string]interface{}{
		"pin": request.PIN,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN set successfully"})
}

func (h *AuthHandler) VerifyPIN(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
		PIN         string `json:"pin"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if !h.authService.VerifyPIN(request.PhoneNumber, request.PIN) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN verified"})
}
