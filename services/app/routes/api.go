package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stebennett/squad-dashboard/pkg/dateutil"
	"github.com/stebennett/squad-dashboard/pkg/statsservice"
)

type Api struct {
	StatsService         statsservice.StatsService
	CycleTimeIssueTypes  []string
	ThroughputIssueTypes []string
	DefectIssueType      string
}

type WorkingCycleTimeQuery struct {
	WeekCount    int       `form:"weekCount"`
	Percentile   int       `form:"percentile"`
	StartTime    time.Time `form:"startTime" time_format:"2006-01-02" time_utc="1"`
	EndOfWeekDay string    `form:"endOfWeekDay"`
}

type ThroughputQuery struct {
	WeekCount    int       `form:"weekCount"`
	StartTime    time.Time `form:"startTime" time_format:"2006-01-02" time_utc="1"`
	EndOfWeekDay string    `form:"endOfWeekDay"`
}

type EscapedDefectQuery struct {
	WeekCount    int       `form:"weekCount"`
	StartTime    time.Time `form:"startTime" time_format:"2006-01-02" time_utc="1"`
	EndOfWeekDay string    `form:"endOfWeekDay"`
}

func (a Api) WorkingCycleTimeGET(c *gin.Context) {
	params := WorkingCycleTimeQuery{
		WeekCount:    12,
		Percentile:   75,
		StartTime:    time.Now(),
		EndOfWeekDay: "Friday",
	}

	c.BindQuery(&params)

	endOfWeekDay, err := dateutil.ParseDayOfWeek(params.EndOfWeekDay)
	if err != nil {
		c.Errors = append(c.Errors, c.Copy().Error(err))
	}

	if len(c.Errors) > 0 {
		c.AbortWithStatus(400)
	}

	project := c.Param("project")

	report, err := a.StatsService.GenerateCycleTime(params.WeekCount, float64(params.Percentile)/100.0, project, a.CycleTimeIssueTypes, params.StartTime, endOfWeekDay)
	if err != nil {
		c.AbortWithStatus(500)
	}

	c.JSON(http.StatusOK, report)
}

func (a Api) ThroughputGET(c *gin.Context) {
	params := ThroughputQuery{
		WeekCount:    12,
		StartTime:    time.Now(),
		EndOfWeekDay: "Friday",
	}

	c.BindQuery(&params)

	endOfWeekDay, err := dateutil.ParseDayOfWeek(params.EndOfWeekDay)
	if err != nil {
		c.Errors = append(c.Errors, c.Copy().Error(err))
	}

	if len(c.Errors) > 0 {
		c.AbortWithStatus(400)
	}

	project := c.Param("project")

	report, err := a.StatsService.GenerateThroughput(params.WeekCount, project, a.ThroughputIssueTypes, params.StartTime, endOfWeekDay)
	if err != nil {
		c.AbortWithStatus(500)
	}

	c.JSON(http.StatusOK, report)
}

func (a Api) EscapedDefectsGET(c *gin.Context) {
	params := EscapedDefectQuery{
		WeekCount:    12,
		StartTime:    time.Now(),
		EndOfWeekDay: "Friday",
	}

	c.BindQuery(&params)

	endOfWeekDay, err := dateutil.ParseDayOfWeek(params.EndOfWeekDay)
	if err != nil {
		c.Errors = append(c.Errors, c.Copy().Error(err))
	}

	if len(c.Errors) > 0 {
		c.AbortWithStatus(400)
	}

	project := c.Param("project")

	report, err := a.StatsService.GenerateEscapedDefects(params.WeekCount, project, a.DefectIssueType, params.StartTime, endOfWeekDay)
	if err != nil {
		c.AbortWithStatus(500)
	}

	c.JSON(http.StatusOK, report)
}

func (a Api) UnplannedWorkGET(c *gin.Context) {
	params := ThroughputQuery{
		WeekCount:    12,
		StartTime:    time.Now(),
		EndOfWeekDay: "Friday",
	}

	c.BindQuery(&params)

	endOfWeekDay, err := dateutil.ParseDayOfWeek(params.EndOfWeekDay)
	if err != nil {
		c.Errors = append(c.Errors, c.Copy().Error(err))
	}

	if len(c.Errors) > 0 {
		c.AbortWithStatus(400)
	}

	project := c.Param("project")

	report, err := a.StatsService.GenerateUnplannedWorkReport(params.WeekCount, project, a.ThroughputIssueTypes, params.StartTime, endOfWeekDay)
	if err != nil {
		c.AbortWithStatus(500)
	}

	c.JSON(http.StatusOK, report)
}
