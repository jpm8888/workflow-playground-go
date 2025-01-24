package repository

import (
	"app/auth/internal/workflow"
	"errors"

	"gorm.io/gorm"
)

// WorkflowRepository handles database operations for workflows
type WorkflowRepository struct {
	db *gorm.DB
	sm *workflow.StateMachine
}

// NewWorkflowRepository creates a new workflow repository
func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{
		db: db,
		sm: workflow.NewWorkflowStateMachine(),
	}
}

// Create a new workflow
func (r *WorkflowRepository) Create(wf *workflow.Workflow) error {
	// Set initial state
	wf.CurrentState = workflow.StateCreated
	return r.db.Create(wf).Error
}

// Update updates an existing workflow
func (r *WorkflowRepository) Update(wf *workflow.Workflow) error {
	return r.db.Save(wf).Error
}

// TransitionState handles state transition with database update
func (r *WorkflowRepository) TransitionState(wf *workflow.Workflow, newState workflow.State) error {
	// Begin a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Attempt to transition the state
	if err := r.sm.TransitionState(wf, newState); err != nil {
		tx.Rollback()
		return err
	}

	// Save the updated workflow
	if err := tx.Save(wf).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// GetByID retrieves a workflow by its ID
func (r *WorkflowRepository) GetByID(id uint) (*workflow.Workflow, error) {
	var wf workflow.Workflow
	result := r.db.First(&wf, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, result.Error
	}
	return &wf, nil
}

// ListWorkflows retrieves workflows with optional filtering
func (r *WorkflowRepository) ListWorkflows(state *workflow.State, limit, offset int) ([]workflow.Workflow, int64, error) {
	var workflows []workflow.Workflow
	var count int64

	query := r.db.Model(&workflow.Workflow{})

	// Optional state filtering
	if state != nil {
		query = query.Where("current_state = ?", *state)
	}

	// Count total workflows
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve workflows with pagination
	err := query.Limit(limit).Offset(offset).Find(&workflows).Error
	if err != nil {
		return nil, 0, err
	}

	return workflows, count, nil
}
