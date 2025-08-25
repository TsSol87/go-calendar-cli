package priority

import (
	"errors"
)

type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

var ErrIsValidPriority = errors.New("priority: invalid priority value")

func (p Priority) Validate() error {
	switch p {
	case PriorityHigh, PriorityMedium, PriorityLow:
		return nil
	default:
		return ErrIsValidPriority
	}
}
