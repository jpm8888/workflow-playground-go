package workflow

import (
	"gorm.io/gorm"
	"time"
)

type WorkflowEngine struct {
	db        *gorm.DB
	workflows map[string]*WorkflowDefinition
}

func NewWorkflowEngine(db *gorm.DB) *WorkflowEngine {
	return &WorkflowEngine{
		db:        db,
		workflows: make(map[string]*WorkflowDefinition),
	}
}

func (e *WorkflowEngine) RegisterWorkflow(workflow *WorkflowDefinition) {
	e.workflows[workflow.ID] = workflow
}

func (e *WorkflowEngine) CreateInstance(workflowID string, data map[string]interface{}) (*WorkflowInstance, error) {
	workflow, exists := e.workflows[workflowID]
	if !exists {
		return nil, ErrWorkflowNotFound
	}

	instance := &WorkflowInstance{
		WorkflowID:   workflowID,
		CurrentState: workflow.InitialState,
		Data:         data,
		StartTime:    time.Now(),
		LastUpdated:  time.Now(),
		Timeout:      workflow.Timeout,
		Status:       "active",
	}

	if err := e.db.Create(instance).Error; err != nil {
		return nil, err
	}

	return instance, nil
}

func (e *WorkflowEngine) HandleEvent(instanceID uint, event Event, data map[string]interface{}) error {
	var instance WorkflowInstance
	if err := e.db.First(&instance, instanceID).Error; err != nil {
		return err
	}

	if instance.Status != "active" {
		return ErrInvalidState
	}

	if time.Since(instance.StartTime) > instance.Timeout {
		instance.Status = "timeout"
		e.db.Save(&instance)
		return ErrWorkflowTimeout
	}

	workflow := e.workflows[instance.WorkflowID]
	transitions := workflow.States[instance.CurrentState]

	for _, transition := range transitions {
		if transition.Event == event && (transition.Guard == nil || transition.Guard(data)) {
			if err := transition.Action(data); err != nil {
				instance.Status = "failed"
				e.db.Save(&instance)
				return ErrActionFailed
			}

			instance.CurrentState = transition.ToState
			instance.Data = data
			instance.LastUpdated = time.Now()
			return e.db.Save(&instance).Error
		}
	}

	return ErrNoValidTransition
}
