package main

import (
	"context"
	"log"
	"strings"

	"github.com/Netflix/go-env"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type Environment struct {
	AWSTagName  string `env:"AWS_COST_TAG_NAME,required=true"`
	AWSTagValue string `env:"AWS_COST_TAG_VALUE,required=true"`
}

func main() {
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := costexplorer.NewFromConfig(cfg)

	startDate := "2022-06-01"
	endDate := "2022-12-31"
	tagKey := environment.AWSTagName
	tagValues := strings.Split(environment.AWSTagValue, ",")
	var nextPageToken string

	params := costexplorer.GetCostAndUsageInput{
		Granularity: types.GranularityMonthly,
		Metrics:     []string{"NetUnblendedCost"},
		TimePeriod: &types.DateInterval{
			Start: &startDate,
			End:   &endDate,
		},
		Filter: &types.Expression{
			Tags: &types.TagValues{
				Key:    &tagKey,
				Values: tagValues,
			},
		},
		NextPageToken: &nextPageToken,
	}

	output, err := client.GetCostAndUsage(context.TODO(), &params)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range output.ResultsByTime {
		for k, v := range result.Total {
			log.Printf("%s => %s [%s]; %s - %s", k, *v.Amount, *v.Unit, *result.TimePeriod.Start, *result.TimePeriod.End)
		}
	}

}
