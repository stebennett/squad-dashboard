package calculator

import (
	"errors"
	"sort"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dateutil"
	"github.com/stebennett/squad-dashboard/pkg/jiramodels"
)

func calculateCycleTime(transitions []jiramodels.JiraTransition, startState string, stopState string) (int, error) {
	// filter only to start states
	startingDates := filter(transitions, func(item jiramodels.JiraTransition) bool {
		return item.ToState == startState
	})
	// filter only to stop states
	stoppingDates := filter(transitions, func(item jiramodels.JiraTransition) bool {
		return item.ToState == stopState
	})

	// sort both ascending
	sort.Slice(startingDates, func(i, j int) bool {
		return startingDates[i].TransitionedAt.Before(startingDates[j].TransitionedAt)
	})
	sort.Slice(stoppingDates, func(i, j int) bool {
		return stoppingDates[i].TransitionedAt.Before(stoppingDates[j].TransitionedAt)
	})

	if len(startingDates) == 0 {
		return 0, errors.New("Failed to calculate cycle times. No starting states found in transitions.")
	}

	if len(stoppingDates) == 0 {
		return 0, errors.New("Failed to calculate cycle times. No stopping states found in transitions.")
	}

	startDate := startingDates[0].TransitionedAt
	stopDate := stoppingDates[len(stoppingDates)-1].TransitionedAt

	daysBetween := dateutil.DaysBetween(startDate, stopDate)
	if daysBetween == 0 {
		// add a day for minimum length
		daysBetween = 1
	}

	return daysBetween, nil
}

func caculateLeadTime(transitions []jiramodels.JiraTransition, createdAt time.Time, stopState string) (int, error) {
	// filter only to stop states
	stoppingDates := filter(transitions, func(item jiramodels.JiraTransition) bool {
		return item.ToState == stopState
	})

	if len(stoppingDates) == 0 {
		return 0, errors.New("Failed to calculate lead times. No stopping states found in transitions.")
	}

	stopDate := stoppingDates[len(stoppingDates)-1].TransitionedAt

	daysBetween := dateutil.DaysBetween(createdAt, stopDate)
	if daysBetween == 0 {
		daysBetween = 1
	}

	return daysBetween, nil
}

func sortByTransitionAtAscending(transitions []jiramodels.JiraTransition) []jiramodels.JiraTransition {
	sort.Slice(transitions, func(i, j int) bool {
		return transitions[i].TransitionedAt.Before(transitions[j].TransitionedAt)
	})
	return transitions
}

func sortByTransitionAtDescending(transitions []jiramodels.JiraTransition) []jiramodels.JiraTransition {
	sort.Slice(transitions, func(i, j int) bool {
		return transitions[i].TransitionedAt.After(transitions[j].TransitionedAt)
	})
	return transitions
}

func filter(transitions []jiramodels.JiraTransition, fn func(transition jiramodels.JiraTransition) bool) (results []jiramodels.JiraTransition) {
	for _, transition := range transitions {
		if fn(transition) {
			results = append(results, transition)
		}
	}
	return
}
