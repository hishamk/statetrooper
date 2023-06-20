package main

import (
	"encoding/json"
	"fmt"

	"github.com/hishamk/statetrooper"
)

// Order represents a custom entity with its current state
type Order struct {
	State *statetrooper.FSM[OrderState]
}

// OrderState represents a transition in the order state machine
// It is compliant with the comparable type constraint since it is composed of
// only comparable types
type OrderState struct {
	Name    string
	Group   string
	Version int
}

var (
	Created    = OrderState{Name: "created", Group: "dropship", Version: 1}
	Picked     = OrderState{Name: "picked", Group: "dropship", Version: 1}
	Packed     = OrderState{Name: "packed", Group: "dropship", Version: 2}
	Shipped    = OrderState{Name: "shipped", Group: "dropship", Version: 1}
	Delivered  = OrderState{Name: "delivered", Group: "dropship", Version: 1}
	Canceled   = OrderState{Name: "canceled", Group: "dropship", Version: 4}
	Reinstated = OrderState{Name: "reinstated", Group: "dropship", Version: 1}
)

func (s OrderState) String() string {
	return fmt.Sprintf("%s:%s:v%d", s.Name, s.Group, s.Version)
}

func main() {
	// Create a new order with the initial state
	order := &Order{State: statetrooper.NewFSM[OrderState](Created, 10)}

	// Define the valid state transitions for the order
	order.State.AddRule(Created, Picked, Canceled)    // Created -> Picked or Canceled
	order.State.AddRule(Picked, Packed, Canceled)     // Picked -> Packed or Canceled
	order.State.AddRule(Packed, Shipped)              // Packed -> Shipped
	order.State.AddRule(Shipped, Delivered)           // Shipped -> Delivered
	order.State.AddRule(Canceled, Reinstated)         // Canceled -> Reinstated
	order.State.AddRule(Reinstated, Picked, Canceled) // Reinstated -> Picked or Canceled

	// Check if a transition is valid
	canTransition := order.State.CanTransition(Picked)
	fmt.Printf("Can transition to %s: %t\n", Picked.Name, canTransition)

	// Transition to picked
	_, err := order.State.Transition(Picked, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Check if a transition to canceled is valid
	canTransition = order.State.CanTransition(Canceled)
	fmt.Printf("Can transition to %s: %t\n", Canceled.Name, canTransition)

	// Transition to canceled
	_, err = order.State.Transition(Canceled, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Check if we can resinstate the order
	canTransition = order.State.CanTransition(Reinstated)
	fmt.Printf("Can transition to %s: %t\n", Reinstated.Name, canTransition)

	// Transition to reinstated
	_, err = order.State.Transition(Reinstated, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Transition to picked
	_, err = order.State.Transition(Picked, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Transition to packed
	_, err = order.State.Transition(Packed, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// Transition to shipped
	_, err = order.State.Transition(
		Shipped,
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
	_, err = order.State.Transition(Delivered, nil)
	if err != nil {
		fmt.Println("Transition error:", err)
	} else {
		fmt.Println("Transition successful. Current state:", order.State.CurrentState())
	}

	// print the current FSM data
	fmt.Println("Current FSM data:", order.State)

	// marshal the current FSM data to JSON
	j, err := json.Marshal(order.State)
	if err != nil {
		fmt.Println("JSON error:", err)
	} else {
		fmt.Println("Current FSM data as JSON:", string(j))
	}

	// unmarshal the FSM data from JSON
	newOrder := &Order{State: statetrooper.NewFSM[OrderState](Created, 10)}
	err = json.Unmarshal(j, newOrder.State)
	if err != nil {
		fmt.Println("JSON error:", err)
	} else {
		fmt.Printf("Unmarshalled FSM data:\n%s\n", newOrder.State)
	}

	fmt.Println("FSM rules as Mermaid diagram:")
	m, err := order.State.GenerateMermaidRulesDiagram()
	if err != nil {
		fmt.Println("Mermaid diagram error:", err)
	} else {
		fmt.Println(m)
	}

	fmt.Println("Transition history as Mermaid diagram:")
	m, err = order.State.GenerateMermaidTransitionHistoryDiagram()
	if err != nil {
		fmt.Println("Mermaid diagram error:", err)
	} else {
		fmt.Println(m)
	}

}
