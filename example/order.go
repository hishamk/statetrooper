package main

import (
	"encoding/json"
	"fmt"

	"github.com/hishamk/statetrooper"
)

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

	// print the current FSM data as JSON
	json, err := json.Marshal(order.State)
	if err != nil {
		fmt.Println("JSON error:", err)
	} else {
		fmt.Println("Current FSM data as JSON:", string(json))
	}

	fmt.Println("FSM rules as Mermaid diagram:")
	fmt.Println(order.State.GenerateMermaidDiagram())
}
