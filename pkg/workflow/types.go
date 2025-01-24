package workflow

import (
	"gorm.io/gorm"
	"time"
)

type State string
type Event string

type WorkflowInstance struct {
	gorm.Model
	WorkflowID   string
	CurrentState State
	Data         map[string]interface{} `gorm:"serializer:json"`
	StartTime    time.Time
	Timeout      time.Duration
	LastUpdated  time.Time
	Status       string // "active", "completed", "failed", "timeout"
}

type StateTransition struct {
	FromState State
	Event     Event
	ToState   State
	Action    func(data map[string]interface{}) error
	Guard     func(data map[string]interface{}) bool
}

type WorkflowDefinition struct {
	ID           string
	InitialState State
	States       map[State][]StateTransition
	Timeout      time.Duration
}
