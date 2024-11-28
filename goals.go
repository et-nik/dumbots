package main

type Goal int

const (
	// GoalUnknown is an unknown goal.
	GoalUnknown Goal = iota
	GoalEnemy        // GoalEnemy is the goal to find and kill the enemy.
	GoalMove         // GoalMove is the goal to move to a specific location.
	GoalSecret       // GoalSecret is the goal to find and secure the secret.
)
