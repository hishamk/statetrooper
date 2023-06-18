```
   ______       __        ______
  / __/ /____ _/ /____   /_  __/______  ___  ___  ___ ____
 _\ \/ __/ _ `/ __/ -_)   / / / __/ _ \/ _ \/ _ \/ -_) __/
/___/\__/\_,_/\__/\__/   /_/ /_/  \___/\___/ .__/\__/_/
                                          /_/
```

_Tiny, no frills finite state machine for Go_

[![GoDoc](https://godoc.org/github.com/hishamk/statetrooper?status.png)](https://pkg.go.dev/github.com/hishamk/statetrooper?tab=doc)
[![Go Coverage](https://github.com/hishamk/statetrooper/wiki/coverage.svg)](https://raw.githack.com/wiki/hishamk/statetrooper/coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/hishamk/statetrooper)](https://goreportcard.com/report/github.com/hishamk/statetrooper)
[![MIT](https://img.shields.io/github/license/hishamk/statetrooper)](https://img.shields.io/github/license/hishamk/statetrooper) ![Code size](https://img.shields.io/github/languages/code-size/hishamk/statetrooper)

StateTrooper is a Go package that provides a finite state machine (FSM) for managing states. It allows you to define and enforce state transitions based on predefined rules.

## Features

- Generic support for different comparable types
- Transition history with metadata
- Thread safe
- Super minimal - no triggers/events or actions/callbacks. For my use case I just needed a structured, serializable way to constrain and track state transitions.

## Installation

To install StateTrooper, use the following command:

```shell
go get github.com/hishamk/statetrooper
```

## Usage

Import the `statetrooper` package into your Go code:

```go
import "github.com/hishamk/statetrooper"
```

Create an instance of the FSM with the desired state enum type and initial state:

```go
fsm := statetrooper.NewFSM[CustomStateEnum](CustomStateEnumA)
```

Add valid transitions between states. AddRule takes variadic parameters for the allowed states:

```go
AddRule(StatusCreated, StatusPicked, StatusCanceled)    // Created -> Picked or Canceled
AddRule(StatusPicked, StatusPacked, StatusCanceled)     // Picked -> Packed or Canceled
AddRule(StatusPacked, StatusShipped)                    // Packed -> Shipped
AddRule(StatusShipped, StatusDelivered)                 // Shipped -> Delivered
AddRule(StatusCanceled, StatusReinstated)               // Canceled -> Reinstated
AddRule(StatusReinstated, StatusPicked, StatusCanceled) // Reinstated -> Picked or Canceled
```

Check if a transition from the current state to the target state is valid:

```go
canTransition := fsm.CanTransition(targetState)
```

Transition the entity from the current state to the target state with no metadata:

```go
newState, err := fsm.Transition(targetState, nil)
if err != nil {
    // Handle the error
}
```

Transition the entity from the current state to the target state with metadata:

```go
	newState, err := fsm.Transition(
		CustomStateEnumB,
		map[string]string{
			"requested_by":  "Mahmoud",
			"logic_version": "1.0",
		})
```

## Benchmark

| Benchmark              | Iterations | Time per Iteration | Memory Allocation per Iteration | Allocations per Iteration |
| ---------------------- | ---------- | ------------------ | ------------------------------- | ------------------------- |
| Benchmark_transition-8 | 363,970    | 2,975 ns/op        | 2,187 B/op                      | 12 allocs/op              |

## Example

Here's an example usage with a custom entity struct and state enum:

```go
type OrderStatusEnum string

// Enum values for the custom entity
const (
	StatusCreated    OrderStatusEnum = "created"
	StatusPicked     OrderStatusEnum = "picked"
	StatusPacked     OrderStatusEnum = "packed"
	StatusShipped    OrderStatusEnum = "shipped"
	StatusDelivered  OrderStatusEnum = "delivered"
	StatusCanceled   OrderStatusEnum = "canceled"
	StatusReinstated OrderStatusEnum = "reinstated"
)

// Order represents a custom entity with its current state
type Order struct {
	State *statetrooper.FSM[OrderStatusEnum]
}

func main() {
	// Create a new order with the initial state
	order := &Order{State: statetrooper.NewFSM[OrderStatusEnum](StatusCreated)}

	// Define the valid state transitions for the order
	order.State.AddRule(StatusCreated, StatusPicked, StatusCanceled)    // Created -> Picked or Canceled
	order.State.AddRule(StatusPicked, StatusPacked, StatusCanceled)     // Picked -> Packed or Canceled
	order.State.AddRule(StatusPacked, StatusShipped)                    // Packed -> Shipped
	order.State.AddRule(StatusShipped, StatusDelivered)                 // Shipped -> Delivered
	order.State.AddRule(StatusCanceled, StatusReinstated)               // Canceled -> Reinstated
	order.State.AddRule(StatusReinstated, StatusPicked, StatusCanceled) // Reinstated -> Picked or Canceled

	// Check if a transition is valid
	canTransition := order.State.CanTransition(StatusPicked)
	fmt.Printf("Can transition to %s: %t\n", StatusPicked, canTransition)

	// Transition to picked
	_, err := order.State.Transition(StatusPicked, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Check if a transition to canceled is valid
	canTransition = order.State.CanTransition(StatusCanceled)
	fmt.Printf("Can transition to %s: %t\n", StatusCanceled, canTransition)

	// Transition to canceled
	_, err = order.State.Transition(StatusCanceled, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Check if we can resinstate the order
	canTransition = order.State.CanTransition(StatusReinstated)
	fmt.Printf("Can transition to %s: %t\n", StatusReinstated, canTransition)

	// Transition to reinstated
	_, err = order.State.Transition(StatusReinstated, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Transition to picked
	_, err = order.State.Transition(StatusPacked, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Transition to shipped
	_, err = order.State.Transition(StatusShipped, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Transition to delivered
	_, err = order.State.Transition(StatusDelivered, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}
}

```

Note that states can be defined using any comparable type, such as strings, int, etc e.g.:

```go
// CustomStateEnum represents the state enum for the custom entity
type CustomStateEnum int

// Enum values for the custom entity
const (
	CustomStateEnumA CustomStateEnum = iota
	CustomStateEnumB
	CustomStateEnumC
)

```

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for details.

## Contributing

Thank you for your interest in contributing! Feel free to PR bug fixes and documentation improvements. For new features or functional alterations, please open an issue for discussion prior to submitting a PR.
