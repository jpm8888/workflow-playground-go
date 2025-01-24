package workflow

import (
	"fmt"
	"sync"
)

// WorkflowState represents the different states a workflow can be in
type State string

const (
	StateCreated    State = "CREATED"
	StatePending    State = "PENDING"
	StateInProgress State = "IN_PROGRESS"
	StateReviewing  State = "REVIEWING"
	StateApproved   State = "APPROVED"
	StateRejected   State = "REJECTED"
	StateCancelled  State = "CANCELLED"
	StateCompleted  State = "COMPLETED"
)

// Workflow represents a workflow instance
type Workflow struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Description  string
	CurrentState State
	Metadata     map[string]interface{} `gorm:"type:jsonb"`
}

// WorkflowStateMachine manages workflow state transitions
type StateMachine struct {
	mu sync.RWMutex
}

// NewWorkflowStateMachine creates a new state machine
func NewWorkflowStateMachine() *StateMachine {
	return &StateMachine{}
}

// ValidateTransition checks if a state transition is allowed
func (sm *StateMachine) ValidateTransition(from, to State) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Define valid state transitions
	validTransitions := map[State][]State{
		StateCreated:    {StatePending, StateCancelled},
		StatePending:    {StateInProgress, StateReviewing, StateCancelled},
		StateInProgress: {StateReviewing, StateCancelled},
		StateReviewing:  {StateApproved, StateRejected},
		StateApproved:   {StateCompleted},
		StateRejected:   {StatePending, StateCancelled},
		StateCancelled:  {},
		StateCompleted:  {},
	}

	allowedStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, state := range allowedStates {
		if state == to {
			return true
		}
	}

	return false
}

// TransitionState attempts to transition a workflow to a new state
func (sm *StateMachine) TransitionState(workflow *Workflow, newState State) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Validate the state transition
	if !sm.ValidateTransition(workflow.CurrentState, newState) {
		return fmt.Errorf("invalid state transition from %s to %s", workflow.CurrentState, newState)
	}

	// Perform any pre-transition validations or actions
	switch newState {
	case StateInProgress:
		// Example: Check if workflow can be started
		if workflow.CurrentState != StatePending {
			return fmt.Errorf("workflow must be in PENDING state to start")
		}
	case StateApproved:
		// Example: Check if workflow meets approval criteria
		if err := sm.validateApproval(workflow); err != nil {
			return err
		}
	}

	// Update the state
	workflow.CurrentState = newState
	return nil
}

// validateApproval contains custom logic for workflow approval
func (sm *StateMachine) validateApproval(workflow *Workflow) error {
	// Example custom validation logic
	// This could check metadata, perform additional checks, etc.
	if workflow.CurrentState != StateReviewing {
		return fmt.Errorf("workflow must be in REVIEWING state for approval")
	}
	return nil
}
