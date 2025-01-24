package workflow

import "errors"

var (
	ErrWorkflowNotFound  = errors.New("workflow not found")
	ErrNoValidTransition = errors.New("no valid transition found")
	ErrWorkflowTimeout   = errors.New("workflow timeout exceeded")
	ErrInvalidState      = errors.New("invalid state")
	ErrActionFailed      = errors.New("action execution failed")
)
