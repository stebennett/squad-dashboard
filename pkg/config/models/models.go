package models

type Config struct {
	JiraProject          string
	JiraToDoStates       []string
	JiraInProgressStates []string
	JiraDoneStates       []string
	NonWorkingDays       []string
}
