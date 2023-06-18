/*
Package statetrooper provides a finite state machine (FSM) implementation for managing states.

MIT License

Copyright (c) 2023 Hisham Khalifa

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package statetrooper

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Transition represents information about a state transition
type Transition[T comparable] struct {
	FromState T                 `json:"from_state"`
	ToState   T                 `json:"to_state"`
	Timestamp *time.Time        `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
}

// FSM represents the finite state machine for managing states
type FSM[T comparable] struct {
	currentState T
	transitions  []Transition[T]
	ruleset      map[T][]T
	mu           sync.Mutex
}

// NewFSM creates a new instance of FSM with predefined transitions
func NewFSM[T comparable](initialState T) *FSM[T] {
	return &FSM[T]{
		currentState: initialState,
		ruleset:      make(map[T][]T),
	}
}

// CanTransition checks if a transition from the current state to the target state is valid
func (fsm *FSM[T]) CanTransition(targetState T) bool {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	return fsm.canTransition(&fsm.currentState, &targetState)
}

// canTransition checks if a transition from one state to another state is valid
func (fsm *FSM[T]) canTransition(fromState *T, toState *T) bool {
	validTransitions, ok := fsm.ruleset[*fromState]
	if !ok {
		return false
	}

	for _, validState := range validTransitions {
		if validState == *toState {
			return true
		}
	}

	return false
}

// AddRule adds a valid transition between two states
func (fsm *FSM[T]) AddRule(fromState T, toState ...T) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	fsm.ruleset[fromState] = append(fsm.ruleset[fromState], toState...)
}

// Transition transitions the entity from the current state to the target state
// if the transition is invalid, an error is returned and the current state is not changed
func (fsm *FSM[T]) Transition(targetState T, metadata map[string]string) (T, error) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	if !fsm.canTransition(&fsm.currentState, &targetState) {
		return fsm.currentState, TransitionError[T]{
			FromState: fsm.currentState,
			ToState:   targetState,
		}
	}

	// Track the transition
	tn := time.Now()
	fsm.transitions = append(
		fsm.transitions,
		Transition[T]{
			FromState: fsm.currentState,
			ToState:   targetState,
			Timestamp: &tn,
			Metadata:  metadata,
		})

	fsm.currentState = targetState

	return fsm.currentState, nil
}

// CurrentState returns the current state of the FSM
func (fsm *FSM[T]) CurrentState() T {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	return fsm.currentState
}

// Transitions returns a slice of all transitions
func (fsm *FSM[T]) Transitions() []Transition[T] {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	// return a copy of the transitions
	transitions := make([]Transition[T], len(fsm.transitions))

	copy(transitions, fsm.transitions)

	return transitions
}

// MarshalJSON serializes the FSM to JSON
func (fsm *FSM[T]) MarshalJSON() ([]byte, error) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	type FSMExport struct {
		CurrentState T               `json:"current_state"`
		Transitions  []Transition[T] `json:"transitions"`
	}

	export := FSMExport{
		CurrentState: fsm.currentState,
		Transitions:  fsm.transitions,
	}

	return json.Marshal(export)
}

// String returns a string representation of the FSM
func (fsm *FSM[T]) String() string {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	// print current state
	currentState := fmt.Sprintf("Current State: %v\n", fsm.currentState)

	// print rules
	rules := "Rules:\n"
	for fromState, toStates := range fsm.ruleset {
		rules += fmt.Sprintf("\t%v -> %v\n", fromState, toStates)
	}

	// print transitions
	transitions := "Transitions:\n"
	for _, transition := range fsm.transitions {
		transitions += fmt.Sprintf("\t%v\n", transition)
	}

	return currentState + rules + transitions

}

// String returns a string representation of the Transition
func (t *Transition[T]) String() string {
	return fmt.Sprintf("Transition from %v to %v at %v with metadata %v", t.FromState, t.ToState, t.Timestamp, t.Metadata)
}
