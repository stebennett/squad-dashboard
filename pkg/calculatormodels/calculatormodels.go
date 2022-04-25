package calculatormodels

import (
	"github.com/lib/pq"
)

type IssueCalculations struct {
	IssueKey         string
	CycleTime        int
	LeadTime         int
	SystemDelayTime  int
	IssueCreatedAt   pq.NullTime
	IssueStartedAt   pq.NullTime
	IssueCompletedAt pq.NullTime
}
