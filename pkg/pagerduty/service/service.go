package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/pagerduty/models"
)

type PagerDutyParams struct {
	AuthToken string
	BaseUrl   string
}

type PagerDutyService struct {
	params     PagerDutyParams
	httpClient *http.Client
}

type PagerDutyTimestamp struct {
	time.Time
}

type PagerDutyResponseEntity struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HtmlUrl string `json:"html_url"`
}

type PagerDutyOncallResponse struct {
	OnCalls []struct {
		User             PagerDutyResponseEntity `json:"user"`
		Schedlule        PagerDutyResponseEntity `json:"schedule"`
		EscalationPolicy PagerDutyResponseEntity `json:"escalation_policy"`
		EscalationLevel  int                     `json:"escalation_level"`
		Start            PagerDutyTimestamp      `json:"start"`
		End              PagerDutyTimestamp      `json:"end"`
	} `json:"oncalls"`
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
	More   bool `json:"more"`
	Total  int  `json:"total"`
}

func NewPagerDutyService(httpClient *http.Client, params PagerDutyParams) *PagerDutyService {
	return &PagerDutyService{
		params:     params,
		httpClient: httpClient,
	}
}

func (pds *PagerDutyService) GetOnCalls(since time.Time, until time.Time) ([]models.OnCall, error) {
	return pds.getOnCalls(0, 25, since, until, []models.OnCall{})
}

func (pds *PagerDutyService) getOnCalls(offset int, limit int, since time.Time, until time.Time, result []models.OnCall) ([]models.OnCall, error) {
	log.Printf("Fetching PagerDuty Oncalls with offset %d", offset)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/oncalls", pds.params.BaseUrl), nil)
	if err != nil {
		return []models.OnCall{}, err
	}

	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", pds.params.AuthToken))

	query := req.URL.Query()
	query.Add("limit", strconv.Itoa(limit))
	query.Add("offset", strconv.Itoa(offset))
	query.Add("since", since.Format(time.RFC3339))
	query.Add("until", until.Format(time.RFC3339))
	req.URL.RawQuery = query.Encode()

	log.Printf("query - %s", req.URL.RawQuery)

	resp, err := pds.httpClient.Do(req)
	if err != nil {
		return []models.OnCall{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []models.OnCall{}, err
	}

	var oncallPageResponse PagerDutyOncallResponse
	err = json.Unmarshal(body, &oncallPageResponse)
	if err != nil {
		return []models.OnCall{}, err
	}

	for _, onCallResp := range oncallPageResponse.OnCalls {
		onCall := models.OnCall{
			User: models.PagerDutyEntitySummary{
				Id:   onCallResp.User.Id,
				Name: onCallResp.User.Summary,
			},
			Schedule: models.PagerDutyEntitySummary{
				Id:   onCallResp.Schedlule.Id,
				Name: onCallResp.Schedlule.Summary,
			},
			EscalationPolicy: models.PagerDutyEntitySummary{
				Id:   onCallResp.EscalationPolicy.Id,
				Name: onCallResp.EscalationPolicy.Summary,
			},
			EscalationLevel: onCallResp.EscalationLevel,
			Start:           onCallResp.Start.Time,
			End:             onCallResp.End.Time,
		}
		result = append(result, onCall)
	}

	if oncallPageResponse.More {
		newOffset := offset + limit
		newResult, err := pds.getOnCalls(newOffset, limit, since, until, result)
		if err != nil {
			return result, err
		}
		result = append(result, newResult...)
	}

	return result, nil
}

func (p *PagerDutyTimestamp) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		fmt.Printf("Failed to unmarshal timestamp - %s", err)
		return err
	}

	if len(raw) == 0 {
		return err
	}

	p.Time, err = time.Parse(time.RFC3339, raw)
	return err
}
