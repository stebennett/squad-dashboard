package calculator

import (
	"errors"
	"sort"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dateutil"
	"github.com/stebennett/squad-dashboard/pkg/jiramodels"
)

func CalculateCycleTime(startDate time.Time, completedDate time.Time) (int, error) {
	daysBetween := dateutil.DaysBetween(startDate, completedDate)
	if daysBetween == 0 {
		// add a day for minimum length
		daysBetween = 1
	}

	return daysBetween, nil
}

func CalculateWorkingCycleTime(startDate time.Time, completedDate time.Time) (int, error) {
	weekDaysBetween := dateutil.WeekDaysBetween(startDate, completedDate)
	if weekDaysBetween == 0 {
		// add a day for minimum length
		weekDaysBetween = 1
	}

	return weekDaysBetween, nil
}

func CaculateLeadTime(transitions []jiramodels.JiraTransition, createdAt time.Time, stopState string) (int, error) {
	// filter only to stop states
	stoppingDates := Filter(transitions, func(item jiramodels.JiraTransition) bool {
		return item.ToState == stopState
	})

	if len(stoppingDates) == 0 {
		return 0, errors.New("failed to calculate lead times. No stopping states found in transitions")
	}

	stopDate := stoppingDates[len(stoppingDates)-1].TransitionedAt

	daysBetween := dateutil.DaysBetween(createdAt, stopDate)
	if daysBetween == 0 {
		daysBetween = 1
	}

	return daysBetween, nil
}

func SortByTransitionAtAscending(transitions []jiramodels.JiraTransition) []jiramodels.JiraTransition {
	sort.Slice(transitions, func(i, j int) bool {
		return transitions[i].TransitionedAt.Before(transitions[j].TransitionedAt)
	})
	return transitions
}

func SortByTransitionAtDescending(transitions []jiramodels.JiraTransition) []jiramodels.JiraTransition {
	sort.Slice(transitions, func(i, j int) bool {
		return transitions[i].TransitionedAt.After(transitions[j].TransitionedAt)
	})
	return transitions
}

func Filter(transitions []jiramodels.JiraTransition, fn func(transition jiramodels.JiraTransition) bool) (results []jiramodels.JiraTransition) {
	for _, transition := range transitions {
		if fn(transition) {
			results = append(results, transition)
		}
	}
	return
}
