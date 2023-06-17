/*
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
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

// CustomStateEnum represents a custom state enum for testing
type CustomStateEnum string

// Enum values for custom state
const (
	CustomStateEnumA CustomStateEnum = "A"
	CustomStateEnumB CustomStateEnum = "B"
	CustomStateEnumC CustomStateEnum = "C"
	CustomStateEnumD CustomStateEnum = "D"
)

func Test_canTransition(t *testing.T) {
	fsm := NewFSM[CustomStateEnum](CustomStateEnumA)
	fsm.AddRule(CustomStateEnumA, CustomStateEnumB)
	fsm.AddRule(CustomStateEnumB, CustomStateEnumC)
	fsm.AddRule(CustomStateEnumC, CustomStateEnumD)

	tests := []struct {
		currentState CustomStateEnum
		targetState  CustomStateEnum
		expected     bool
	}{
		{CustomStateEnumA, CustomStateEnumB, true},
		{CustomStateEnumA, CustomStateEnumC, false},
		{CustomStateEnumB, CustomStateEnumA, false},
		{CustomStateEnumB, CustomStateEnumC, true},
		{CustomStateEnumC, CustomStateEnumA, false},
		{CustomStateEnumC, CustomStateEnumB, false},
		{CustomStateEnumC, CustomStateEnumC, false},
		{CustomStateEnumC, CustomStateEnumD, true},
	}

	for _, test := range tests {
		result := fsm.canTransition(test.currentState, test.targetState)
		if result != test.expected {
			t.Errorf("canTransition(%v, %v) = %v, expected %v", test.currentState, test.targetState, result, test.expected)
		}
	}
}

func Test_transition(t *testing.T) {
	fsm := NewFSM[CustomStateEnum](CustomStateEnumA)
	fsm.AddRule(CustomStateEnumA, CustomStateEnumB)
	fsm.AddRule(CustomStateEnumB, CustomStateEnumC)

	tests := []struct {
		targetState CustomStateEnum
		expected    CustomStateEnum
		wantErr     bool
	}{
		{CustomStateEnumB, CustomStateEnumB, false}, // Valid state transition
		{CustomStateEnumB, CustomStateEnumB, true},  // Invalid state transition (already in target state)
		{CustomStateEnumA, CustomStateEnumB, true},  // Invalid state transition (no transition from current state to target state)
		{CustomStateEnumC, CustomStateEnumC, false}, // Valid state transition
		{CustomStateEnumD, CustomStateEnumC, true},  // Invalid state transition (no transition from current state to target state)
	}

	for _, test := range tests {
		newState, err := fsm.Transition(test.targetState, nil)
		if (err != nil) != test.wantErr {
			t.Errorf("Transition(%v, %v) returned error: %v, wantErr: %v", *fsm.CurrentState, test.targetState, err, test.wantErr)
		}

		if *fsm.CurrentState != test.expected {
			t.Errorf("Transition(%v, %v) did not update the current state to %v", *fsm.CurrentState, test.targetState, test.expected)
		}

		if newState != nil && *newState != test.expected {
			t.Errorf("Transition(%v, %v) did not return the expected new state of %v", *fsm.CurrentState, test.targetState, test.expected)
		}
	}
}

func Test_transitionTracking(t *testing.T) {
	fsm := NewFSM[CustomStateEnum](CustomStateEnumA)
	fsm.AddRule(CustomStateEnumA, CustomStateEnumB)
	fsm.AddRule(CustomStateEnumB, CustomStateEnumC)

	metadata1 := map[string]string{
		"requested_by":  "Mahmoud",
		"logic_version": "1.0",
	}

	// Perform the first transition
	_, err := fsm.Transition(CustomStateEnumB, metadata1)
	if err != nil {
		t.Errorf("Transition(%v, %v) returned an error: %v", fsm.CurrentState, CustomStateEnumB, err)
	}

	time.Sleep(1 * time.Millisecond) // Add slight delay between transitions

	metadata2 := map[string]string{
		"requested_by":  "John",
		"logic_version": "1.1",
	}

	// Perform the second transition
	_, err = fsm.Transition(CustomStateEnumC, metadata2)
	if err != nil {
		t.Errorf("Transition(%v, %v) returned an error: %v", fsm.CurrentState, CustomStateEnumC, err)
	}

	// Retrieve the transition tracker
	transitionTrack := fsm.Transitions

	// Verify the number of entries in the transition tracker
	if len(transitionTrack) != 2 {
		t.Errorf("Transition tracker does not contain the expected number of entries. Got %d, expected 2", len(transitionTrack))
	}

	// Get the transition timestamps in order
	var timestamps []time.Time
	for timestamp := range transitionTrack {
		timestamps = append(timestamps, timestamp)
	}
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i].Before(timestamps[j])
	})

	// Check each transition in the tracker
	expectedTransitions := []struct {
		FromState CustomStateEnum
		ToState   CustomStateEnum
		Timestamp time.Time
		Metadata  map[string]string
	}{
		{
			FromState: CustomStateEnumA,
			ToState:   CustomStateEnumB,
			Timestamp: timestamps[0],
			Metadata:  metadata1,
		},
		{
			FromState: CustomStateEnumB,
			ToState:   CustomStateEnumC,
			Timestamp: timestamps[1],
			Metadata:  metadata2,
		},
	}

	for i, timestamp := range timestamps {
		tracker := transitionTrack[timestamp]
		expected := expectedTransitions[i]

		if tracker.FromState != expected.FromState {
			t.Errorf("Transition tracker has incorrect FromState. Got %v, expected %v", tracker.FromState, expected.FromState)
		}

		if tracker.ToState != expected.ToState {
			t.Errorf("Transition tracker has incorrect ToState. Got %v, expected %v", tracker.ToState, expected.ToState)
		}

		// Allow a small delta in the timestamp comparison due to slight time difference
		if tracker.Timestamp.Sub(expected.Timestamp) > time.Second {
			t.Errorf("Transition tracker has incorrect Timestamp. Got %v, expected within 1 second", tracker.Timestamp)
		}

		// Deep compare metadata
		if !reflect.DeepEqual(tracker.Metadata, expected.Metadata) {
			t.Errorf("Transition tracker has incorrect Metadata. Got %v, expected %v", tracker.Metadata, expected.Metadata)
		}
	}
}

func Test_concurrency(t *testing.T) {
	fsm := NewFSM[CustomStateEnum](CustomStateEnumA)
	fsm.AddRule(CustomStateEnumA, CustomStateEnumB)
	fsm.AddRule(CustomStateEnumB, CustomStateEnumC)

	var wg sync.WaitGroup

	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 1000; j++ {
				fsm.Transition(CustomStateEnumB, nil)
				fsm.Transition(CustomStateEnumC, nil)
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func Test_jsonMarshal(t *testing.T) {
	fsm := NewFSM[CustomStateEnum](CustomStateEnumA)
	fsm.AddRule(CustomStateEnumA, CustomStateEnumB)
	fsm.AddRule(CustomStateEnumB, CustomStateEnumC)

	fsm.Transition(
		CustomStateEnumB,
		map[string]string{
			"requested_by":  "Mahmoud",
			"logic_version": "1.0",
		})

	fsm.Transition(
		CustomStateEnumC,
		map[string]string{
			"requested_by":  "John",
			"logic_version": "1.1",
		})

	_, err := json.Marshal(fsm)

	if err != nil {
		t.Errorf("JSON() returned an error: %v", err)
	}
}

func Benchmark_transition(b *testing.B) {
	// CustomEntity represents a custom entity with its current state
	type CustomEntity struct {
		State *CustomStateEnum
	}

	entity := &CustomEntity{State: new(CustomStateEnum)}

	fsm := NewFSM[CustomStateEnum](CustomStateEnumA)
	fsm.AddRule(CustomStateEnumA, CustomStateEnumB)
	fsm.AddRule(CustomStateEnumB, CustomStateEnumC)
	fsm.AddRule(CustomStateEnumC, CustomStateEnumB)
	fsm.AddRule(CustomStateEnumB, CustomStateEnumA)

	tests := []struct {
		targetState CustomStateEnum
	}{
		{CustomStateEnumB},
		{CustomStateEnumC},
		{CustomStateEnumB},
		{CustomStateEnumC},
		{CustomStateEnumB},
		{CustomStateEnumA},
	}

	var err error

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			entity.State, err = fsm.Transition(test.targetState, nil)
			if err != nil {
				b.Errorf("Transition returned an error: %v", err)
			}
		}
	}
}
