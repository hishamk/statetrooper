```markdown
   ______       __        ______                          
  / __/ /____ _/ /____   /_  __/______  ___  ___  ___ ____
 _\ \/ __/ _ `/ __/ -_)   / / / __/ _ \/ _ \/ _ \/ -_) __/
/___/\__/\_,_/\__/\__/   /_/ /_/  \___/\___/ .__/\__/_/   
                                          /_/              
```
*Finite State Machine for Go*

[![GoDoc](https://godoc.org/github.com/hishamk/statetrooper?status.png)](https://pkg.go.dev/github.com/hishamk/statetrooper/v2?tab=doc)
[![Go report card](https://goreportcard.com/badge/github.com/hishamk/statetrooper)](https://goreportcard.com/report/github.com/hishamk/statetrooper)
[![Test coverage](http://gocover.io/_badge/github.com/hishamk/statetrooper)](https://gocover.io/github.com/hishamk/statetrooper)



StateTrooper is a Go package that provides a finite state machine (FSM) for managing states. It allows you to define and enforce state transitions based on predefined rules.

## Features
- Generic support for different state types
- Ability to define valid state transitions
- Transition error handling
- Transition history
- Thread safe

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

   Add valid transitions between states:

   ```go
   fsm.AddRule(CustomStateEnumA, CustomStateEnumB)
   fsm.AddRule(CustomStateEnumB, CustomStateEnumC)
   ```

   Check if a transition from the current state to the target state is valid:

   ```go
   canTransition := fsm.CanTransition(targetState)
   ```

   Transition the entity from the current state to the target state:

   ```go
   newState, err := fsm.Transition(targetState)
   if err != nil {
       // Handle the error
   }
   ```
   :

   ```go
   fsm.Transitions
   ```


## Benchmark
| Benchmark                | Iterations | Time per Iteration | Memory Allocation per Iteration | Allocations per Iteration |
|--------------------------|------------|--------------------|---------------------------------|---------------------------|
| Benchmark_transition-8   | 442,840    | 3,162 ns/op        | 2,117 B/op                      | 18 allocs/op              |


## Example
Here's an example usage with a custom entity struct and state enum:

```go
package main

import (
	"fmt"
	"github.com/hishamk/statetrooper"
)

// CustomStateEnum represents the state enum for the custom entity
type CustomStateEnum int

// Enum values for the custom entity
const (
	CustomStateEnumA CustomStateEnum = iota
	CustomStateEnumB
	CustomStateEnumC
)

/* Note that states can be defined using any comparable type, such as strings, e.g.:

type CustomStateEnum string

const (
	CustomStateEnumA CustomStateEnum = "Created"
	CustomStateEnumB CustomStateEnum = "Packed"
	CustomStateEnumC CustomStateEnum = "Shipped"
)

*/

// CustomEntity represents a custom entity with its current state
type CustomEntity struct {
	State *CustomStateEnum
}

func main() {
	entity := &CustomEntity{State: new(CustomStateEnum)}
	fsm := statetrooper.NewFSM[CustomStateEnum](CustomStateEnumA)
	fsm.AddRule(CustomStateEnumA, CustomStateEnumB)
	fsm.AddRule(CustomStateEnumB, CustomStateEnumC)

	// Check if a transition is valid
	canTransition := fsm.CanTransition(CustomStateEnumB)
	fmt.Println("Can transition to B:", canTransition)

	// Transition to a new state
	entity.State, err := fsm.Transition(CustomStateEnumB)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", entity.State)
	}
}
```

## License
This package is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for details.
