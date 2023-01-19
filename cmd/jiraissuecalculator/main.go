package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Netflix/go-env"
	"github.com/stebennett/squad-dashboard/cmd/jiraissuecalculator/calculator"
	"github.com/stebennett/squad-dashboard/pkg/configrepository"
	"github.com/stebennett/squad-dashboard/pkg/jira/repo/calculationsrepository"
	"github.com/stebennett/squad-dashboard/pkg/jira/repo/issuerepository"
	"golang.org/x/exp/slices"
)

type Environment struct {
	JiraProject        string `env:"JIRA_PROJECT,required=true"`
	WorkToDoStates     string `env:"WORK_TODO_STATES,required=true"`
	WorkStartStates    string `env:"WORK_START_STATES,required=true"`
	WorkCompleteStates string `env:"WORK_COMPLETE_STATES,required=true"`
	ReportStartDate    string `env:"JIRA_REPORT_START_DATE,required=true"`
	ReportIssueTypes   string `env:"JIRA_REPORT_ISSUE_TYPES,required=true"`
	ReportEndStates    string `env:"JIRA_REPORT_END_STATES,required=true"`
}

func main() {
	// create a new database to store calculations
	db, err := connectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	issueRepo := issuerepository.NewPostgresIssueRepository(db)
	configRepo := configrepository.NewPostgresConfigRepository(db)
	calculationsRepo := calculationsrepository.NewPostgresJiraCalculationsRepository(db)

	// load environment vars
	var environment Environment
	_, err = env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running updates for project %s", environment.JiraProject)

	// drop all calculations
	_, err = calculationsRepo.DropAllCalculations(context.Background(), environment.JiraProject)
	if err != nil {
		log.Fatalf("Failed to drop existing calculations. %s", err)
	}

	// fetch all issues and set create year-week
	_, err = setCreateDates(issueRepo, calculationsRepo, environment.JiraProject)
	if err != nil {
		log.Fatalf("Failed to set created year-week for issues. %s", err)
	}
	log.Println("Completed update of created year-week for issues")

	// fetch all issues started and set started year-week
	_, err = setStartDates(issueRepo, calculationsRepo, environment.JiraProject, strings.Split(environment.WorkStartStates, ","), strings.Split(environment.WorkToDoStates, ","))
	if err != nil {
		log.Fatalf("Failed to set started year-week for issues. %s", err)
	}
	log.Println("Completed update of started year-week for issues")

	// fetch all issues completed and set complete year-week
	_, err = setCompleteDates(issueRepo, calculationsRepo, environment.JiraProject, strings.Split(environment.WorkCompleteStates, ","))
	if err != nil {
		log.Fatalf("Failed to set completed year-week for issues. %s", err)
	}
	log.Println("Completed update of completed year-week for issues")

	// fetch all issues completed and set cycle time (working and complete)
	_, err = setCycleTimeForCompletedIssues(issueRepo, configRepo, calculationsRepo, environment.JiraProject)
	if err != nil {
		log.Fatalf("Failed to set cycle time. %s", err)
	}
	log.Println("Completed updating cycle time for completed issues")

	// set number items completed for a given week
	_, err = setNumberOfItemsCompletedByWeek(issueRepo, calculationsRepo, environment.ReportStartDate, environment.JiraProject, strings.Split(environment.ReportIssueTypes, ","), strings.Split(environment.ReportEndStates, ","))
	if err != nil {
		log.Fatalf("Failed to set number of items completed by week. %s", err)
	}
	log.Println("Completed number of items completed reports")

	// set number items started for a given week
	_, err = setNumberOfItemsStartedByWeek(issueRepo, calculationsRepo, environment.ReportStartDate, environment.JiraProject, strings.Split(environment.ReportIssueTypes, ","))
	if err != nil {
		log.Fatalf("Failed to set number of items completed by week. %s", err)
	}
	log.Println("Completed number of items started reports")

	log.Printf("All calculations complete for project %s", environment.JiraProject)
}

func connectToDatabase() (*sql.DB, error) {
	var err error
	var db *sql.DB
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	fmt.Println("Database initialised")
	return db, nil
}

func setCreateDates(issuesRepo issuerepository.IssueRepository, calaculationsRepo calculationsrepository.JiraCalculationsRepository, project string) (int64, error) {
	issues, err := issuesRepo.GetIssues(context.Background(), project)
	if err != nil {
		return -1, err
	}

	updatedCount := int64(0)

	for _, issue := range issues {
		year, week := issue.CreatedAt.UTC().ISOWeek()

		rowsChanged, err := calaculationsRepo.SaveCreateDates(context.Background(), issue.Key, year, week, issue.CreatedAt)
		if err != nil {
			return updatedCount, err
		}

		updatedCount = updatedCount + rowsChanged
	}

	return updatedCount, nil
}

func setStartDates(issuesRepo issuerepository.IssueRepository, calculationsRepo calculationsrepository.JiraCalculationsRepository, project string, workStartStates []string, workToDoStates []string) (int64, error) {
	transitions, err := issuesRepo.GetTransitionTimeByStateChanges(context.Background(), project, workToDoStates, workStartStates)
	if err != nil {
		return -1, err
	}

	updatedCount := int64(0)

	for issueKey, transitionTime := range transitions {
		year, week := transitionTime.UTC().ISOWeek()

		rowsChanged, err := calculationsRepo.SaveStartDates(context.Background(), issueKey, year, week, transitionTime)
		if err != nil {
			return updatedCount, err
		}

		updatedCount = updatedCount + rowsChanged
	}

	return updatedCount, nil
}

func setCompleteDates(issuesRepo issuerepository.IssueRepository, calaculationsRepo calculationsrepository.JiraCalculationsRepository, project string, workCompleteStates []string) (int64, error) {
	transitions, err := issuesRepo.GetTransitionTimeByToState(context.Background(), project, workCompleteStates)
	if err != nil {
		return -1, err
	}

	updatedCount := int64(0)

	for issueKey, transitionTime := range transitions {
		// get the current state of the issue
		allTransitions, err := issuesRepo.GetTransitionsForIssue(context.Background(), issueKey)
		if err != nil {
			return updatedCount, err
		}

		sort.Slice(allTransitions, func(i, j int) bool {
			return allTransitions[i].TransitionedAt.After(allTransitions[j].TransitionedAt)
		})
		// is the issue actually complete?
		if slices.Contains(workCompleteStates, allTransitions[0].ToState) {
			year, week := transitionTime.UTC().ISOWeek()

			endState, err := issuesRepo.GetEndStateForIssue(context.Background(), issueKey, transitionTime)
			if err != nil {
				return updatedCount, err
			}

			rowsChanged, err := calaculationsRepo.SaveCompleteDates(context.Background(), issueKey, year, week, transitionTime, endState)
			if err != nil {
				return updatedCount, err
			}

			updatedCount = updatedCount + rowsChanged
		}
	}

	return updatedCount, nil
}

func setCycleTimeForCompletedIssues(issuesRepo issuerepository.IssueRepository, configRepo configrepository.ConfigRepository, calaculationsRepo calculationsrepository.JiraCalculationsRepository, project string) (int64, error) {
	calculations, err := calaculationsRepo.GetCompletedIssues(context.Background(), project)
	if err != nil {
		return -1, err
	}

	datesToExclude, err := configRepo.GetNonWorkingDays(context.Background(), project)
	if err != nil {
		return -1, err
	}

	updatedCount := int64(0)

	for issueKey, calculations := range calculations {
		if !calculations.IssueCompletedAt.Valid {
			return updatedCount, fmt.Errorf("null completed time for issue %s", issueKey)
		}

		var cycleTime, workingCycleTime int
		if !calculations.IssueStartedAt.Valid {
			cycleTime, workingCycleTime = 1, 1 // default value when something moves to done from todo state (ie no start time)
		} else {
			cycleTime, err = calculator.CalculateCycleTime(calculations.IssueStartedAt.Time, calculations.IssueCompletedAt.Time)
			if err != nil {
				return updatedCount, err
			}

			workingCycleTime, err = calculator.CalculateWorkingCycleTime(calculations.IssueStartedAt.Time, calculations.IssueCompletedAt.Time, datesToExclude)
			if err != nil {
				return updatedCount, err
			}
		}

		rowsChanged, err := calaculationsRepo.SaveCycleTime(context.Background(), issueKey, cycleTime, workingCycleTime)
		if err != nil {
			return updatedCount, err
		}

		updatedCount = updatedCount + rowsChanged
	}

	return updatedCount, nil
}

func setNumberOfItemsCompletedByWeek(repo issuerepository.IssueRepository, calaculationsRepo calculationsrepository.JiraCalculationsRepository, startDateStr string, project string, issueTypes []string, endStates []string) (int64, error) {
	startDate, err := time.Parse("2006-01-02T15:04:05Z", startDateStr)
	if err != nil {
		return -1, err
	}

	totalUpdates := int64(0)

	for {
		if startDate.After(time.Now()) {
			break
		}

		endDate := startDate.AddDate(0, 0, 7)
		issues, err := calaculationsRepo.GetIssuesCompletedBetweenDates(context.Background(), project, startDate, endDate, issueTypes, endStates)
		if err != nil {
			return totalUpdates, err
		}

		rowsUpdated, err := repo.SetIssuesCompletedInWeekStarting(context.Background(), project, startDate, len(issues))
		if err != nil {
			return totalUpdates, err
		}

		totalUpdates = totalUpdates + rowsUpdated
		startDate = endDate
	}

	return totalUpdates, nil
}

func setNumberOfItemsStartedByWeek(repo issuerepository.IssueRepository, calaculationsRepo calculationsrepository.JiraCalculationsRepository, startDateStr string, project string, issueTypes []string) (int64, error) {
	startDate, err := time.Parse("2006-01-02T15:04:05Z", startDateStr)
	if err != nil {
		return -1, err
	}

	totalUpdates := int64(0)

	for {
		if startDate.After(time.Now()) {
			break
		}

		endDate := startDate.AddDate(0, 0, 7)

		issues, err := calaculationsRepo.GetIssuesStartedBetweenDates(context.Background(), project, startDate, endDate, issueTypes)
		if err != nil {
			return totalUpdates, err
		}

		rowsUpdated, err := repo.SetIssuesStartedInWeekStarting(context.Background(), project, startDate, len(issues))
		if err != nil {
			return totalUpdates, err
		}

		totalUpdates = totalUpdates + rowsUpdated
		startDate = endDate
	}

	return totalUpdates, nil
}
