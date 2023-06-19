_Tiny, no frills finite state machine for Go_
![](st-logo.png)

[![GoDoc](https://godoc.org/github.com/hishamk/statetrooper?status.png)](https://pkg.go.dev/github.com/hishamk/statetrooper?tab=doc)
[![Go Coverage](https://github.com/hishamk/statetrooper/wiki/coverage.svg)](https://raw.githack.com/wiki/hishamk/statetrooper/coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/hishamk/statetrooper)](https://goreportcard.com/report/github.com/hishamk/statetrooper)
[![MIT](https://img.shields.io/github/license/hishamk/statetrooper)](https://img.shields.io/github/license/hishamk/statetrooper) ![Code size](https://img.shields.io/github/languages/code-size/hishamk/statetrooper)

StateTrooper is a Go package that provides a finite state machine (FSM) for managing states. It allows you to define and enforce state transitions based on predefined rules.

## Features

- Generic support for different comparable types.
- Transition history with metadata.
- Thread safe.
- Super minimal - no triggers/events or actions/callbacks. For my use case I just needed a structured, serializable way to constrain and track state transitions.
- Generate [Mermaid.js](https://mermaid.js.org) diagrams for the transition rules.

![Mermaid.js diagram](order-diagram.png)

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
// Created -> Picked or Canceled
AddRule(StatusCreated, StatusPicked, StatusCanceled)
// Picked -> Packed or Canceled
AddRule(StatusPicked, StatusPacked, StatusCanceled)
// Packed -> Shipped
AddRule(StatusPacked, StatusShipped)
// Shipped -> Delivered
AddRule(StatusShipped, StatusDelivered)
// Canceled -> Reinstated
AddRule(StatusCanceled, StatusReinstated)
// Reinstated -> Picked or Canceled
AddRule(StatusReinstated, StatusPicked, StatusCanceled)
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

Generate Mermaid.js diagram:

```go
diagram, _ :=order.State.GenerateMermaidDiagram()
```

_simply use the generated Mermaid code with your Mermaid visualizer to generate the diagram_

```markdown
graph LR;
shipped;
canceled;
reinstated;
created;
picked;
packed;
created --> picked;
created --> canceled;
picked --> packed;
picked --> canceled;
packed --> shipped;
shipped --> delivered;
canceled --> reinstated;
reinstated --> picked;
reinstated --> canceled;
```

![Mermaid.js diagram](order-diagram.png)

## Benchmarks

| Benchmark                    | Operations | Time per Operation | Memory Allocated per Operation |
| ---------------------------- | ---------- | ------------------ | ------------------------------ |
| Benchmark_singleTransition   | 5,166,985  | 273.8 ns/op        | 314 allocs/op                  |
| Benchmark_twoTransitions     | 2,835,214  | 513.6 ns/op        | 577 allocs/op                  |
| Benchmark_accessCurrentState | 75,695,847 | 14.36 ns/op        | 0 allocs/op                    |
| Benchmark_accessTransitions  | 39,356,628 | 28.74 ns/op        | 48 allocs/op                   |
| Benchmark_marshalJSON        | 1,000,000  | 1,174 ns/op        | 384 allocs/op                  |
| Benchmark_unmarshalJSON      | 318,949    | 3,741 ns/op        | 1,240 allocs/op                |

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

	// Created -> Picked or Canceled
	order.State.AddRule(StatusCreated, StatusPicked, StatusCanceled)
	// Picked -> Packed or Canceled
	order.State.AddRule(StatusPicked, StatusPacked, StatusCanceled)
	// Packed -> Shipped
	order.State.AddRule(StatusPacked, StatusShipped)
	// Shipped -> Delivered
	order.State.AddRule(StatusShipped, StatusDelivered)
	// Canceled -> Reinstated
	order.State.AddRule(StatusCanceled, StatusReinstated)
	// Reinstated -> Picked or Canceled
	order.State.AddRule(StatusReinstated, StatusPicked, StatusCanceled)

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
	_, err = order.State.Transition(StatusPicked, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Transition to packed
	_, err = order.State.Transition(StatusPacked, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Transition to shipped
	_, err = order.State.Transition(
		StatusShipped,
		map[string]string{
			"carrier":         "Aramex",
			"tracking_number": "1234567890",
		})

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

	// print the current FSM data
	fmt.Println("Current FSM data:", order.State)
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

## Serialization

Current state, transition history and any metadata can be marshalled into JSON.

```go
json, err := json.Marshal(order.State)
if err != nil {
	fmt.Println("JSON error:", err)
} else {
	fmt.Println("Current FSM data as JSON:", string(json))
}
```

Output:

```json
{
  "current_state": "delivered",
  "transitions": [
    {
      "from_state": "created",
      "to_state": "picked",
      "timestamp": "2023-06-18T11:44:42.776422+03:00",
      "metadata": null
    },
    {
      "from_state": "picked",
      "to_state": "canceled",
      "timestamp": "2023-06-18T11:44:42.77643+03:00",
      "metadata": null
    },
    {
      "from_state": "canceled",
      "to_state": "reinstated",
      "timestamp": "2023-06-18T11:44:42.776435+03:00",
      "metadata": null
    },
    {
      "from_state": "reinstated",
      "to_state": "picked",
      "timestamp": "2023-06-18T11:44:42.77644+03:00",
      "metadata": null
    },
    {
      "from_state": "picked",
      "to_state": "packed",
      "timestamp": "2023-06-18T11:44:42.776442+03:00",
      "metadata": null
    },
    {
      "from_state": "packed",
      "to_state": "shipped",
      "timestamp": "2023-06-18T11:44:42.776451+03:00",
      "metadata": {
        "carrier": "Aramex",
        "tracking_number": "1234567890"
      }
    },
    {
      "from_state": "shipped",
      "to_state": "delivered",
      "timestamp": "2023-06-18T11:44:42.776454+03:00",
      "metadata": null
    }
  ]
}
```

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for details.

## Contributing

Thank you for your interest in contributing! Feel free to PR bug fixes and documentation improvements. For new features or functional alterations, please open an issue for discussion prior to submitting a PR.

## Logo

Synthwave title text generated courtesy of [Text Effect](https://www.textstudio.com/).
