package statetrooper

import "fmt"

// TransitionError represents an error that occurs during a state transition
type TransitionError[T comparable] struct {
	FromState T
	ToState   T
}

func (err TransitionError[T]) Error() string {
	return fmt.Sprintf("invalid state transition from %v to %v", err.FromState, err.ToState)
}
