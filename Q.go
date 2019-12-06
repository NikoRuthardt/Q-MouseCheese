package main

import (
	"math"
	"math/rand"
	"time"
)

//State interface
type State interface{}

//Agent represents the Q learning Agent
type Agent struct {
	QTable map[State][]float64 // Q - Table
	ε      float64             // exploration
	/* to begin ε=1 */
	α float64 // learning Rate
	/* often used α= 0.1*/
	γ float64 // discount factor
	/*  A factor of 0 for γ will make the agent short sighted  */
	actions int
	/* space of all possible actions */
}

//NewAgent constructor for Q-Agent
func NewAgent(actions int) *Agent {
	return &Agent{
		QTable:  map[State][]float64{},
		ε:       1,
		α:       0.3,
		γ:       0.8,
		actions: actions,
	}
}

func (a *Agent) updateQ(state State, action int, reward float64, newState State) {

	oldQ, _ := maxDir(a.getActions(state))
	optimalFuture, _ := maxDir(a.getActions(newState))
	/*** Q-Function ***/
	/* new Q[s,a] = (1 - α) * Q[s,a] + α * (reward + γ * optimalFuture) */
	a.QTable[state][action] = (1-a.α)*oldQ + a.α*(reward+a.γ*optimalFuture)
}

func (a *Agent) chooseAction(state State) (int, float64) {
	rand.Seed(time.Now().UTC().UnixNano())

	a.ε = math.Max(a.ε*0.995, a.ε*0.005) //let ε --> go slowly to  0

	if rand.Float64() < a.ε {
		return rand.Intn(a.actions), a.ε // choose random Action
	}
	actions := a.getActions(state) // choose action based on Q-Table
	_, action := maxDir(actions)
	return action, a.ε
}

func (a *Agent) getActions(state State) []float64 {
	// if state exist -> return actions
	if actions, ok := a.QTable[state]; ok {
		return actions
	}
	// if state not exist -> create actions for state -> return actions
	a.QTable[state] = make([]float64, a.actions)
	return a.QTable[state]
}

func maxDir(actions []float64) (float64, int) {
	rand.Seed(time.Now().UTC().UnixNano())
	// choose max "direction" from actions
	max := actions[0] // -> 0 UP, 1 Down, 2 Right, 3 Left
	index := 0
	for i, a := range actions {
		if a > max {
			max = a
			index = i
		}
	}
	return max, index

}
