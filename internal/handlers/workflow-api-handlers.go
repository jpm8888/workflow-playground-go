package handlers

import (
	"app/auth/internal/repository"
	"app/auth/internal/workflow"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// WorkflowHandler manages HTTP endpoints for workflow operations
type WorkflowHandler struct {
	repo *repository.WorkflowRepository
}

// NewWorkflowHandler creates a new workflow handler
func NewWorkflowHandler(repo *repository.WorkflowRepository) *WorkflowHandler {
	return &WorkflowHandler{repo: repo}
}

// CreateWorkflow handles workflow creation
func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	var wf workflow.Workflow
	if err := c.ShouldBindJSON(&wf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&wf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wf)
}

// TransitionWorkflowState handles state transitions
func (h *WorkflowHandler) TransitionWorkflowState(c *gin.Context) {
	// Parse workflow ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Get the workflow
	wf, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Parse new state from request
	var transitionRequest struct {
		NewState workflow.State `json:"new_state"`
	}
	if err := c.ShouldBindJSON(&transitionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Attempt state transition
	if err := h.repo.TransitionState(wf, transitionRequest.NewState); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wf)
}

// ListWorkflows retrieves workflows with optional filtering
func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	// Parse query parameters
	stateStr := c.Query("state")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	// Convert state if provided
	var state *workflow.State
	if stateStr != "" {
		workflowState := workflow.State(stateStr)
		state = &workflowState
	}

	// Retrieve workflows
	workflows, total, err := h.repo.ListWorkflows(state, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflows": workflows,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// SetupRoutes configures workflow-related routes
func SetupRoutes(r *gin.Engine, repo *repository.WorkflowRepository) {
	handler := NewWorkflowHandler(repo)

	// Workflow routes
	workflows := r.Group("/workflows")
	{
		workflows.POST("", handler.CreateWorkflow)
		workflows.GET("", handler.ListWorkflows)
		workflows.PUT("/:id/transition", handler.TransitionWorkflowState)
	}
}
