package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gaRep "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("starting connection...")
	config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		// RedirectURL:  "https://developers.google.com/oauthplayground",
	}

	token := &oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		Expiry:       time.Now(),
	}

	ctx := context.Background()
	gaService, err := gaRep.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if err != nil {
		fmt.Println("error creating service:", err)
	}

	report, err := getReport(gaService)
	if err != nil {
		fmt.Println("error creating report:", err)
	}

	printResponse(report)

}

func getReport(svc *gaRep.Service) (*gaRep.GetReportsResponse, error) {
	req := &gaRep.GetReportsRequest{
		// Our request contains only one request
		// So initialise the slice with one ga.ReportRequest object
		ReportRequests: []*gaRep.ReportRequest{
			// Create the ReportRequest object.
			{
				ViewId: "261593436",
				DateRanges: []*gaRep.DateRange{
					// Create the DateRange object.
					{StartDate: "2022-09-01", EndDate: "today"},
				},
				Metrics: []*gaRep.Metric{
					// Create the Metrics object.
					{Expression: "ga:users", Alias: "Users"},
					{Expression: "ga:sessions", Alias: "Sessions"},
					{Expression: "ga:transactions", Alias: "Transactions"},
					{Expression: "ga:transactionRevenue", Alias: "Revenue"},
				},
				MetricFilterClauses: []*gaRep.MetricFilterClause{
					{
						Filters: []*gaRep.MetricFilter{{
							MetricName:      "ga:transactions",
							Operator:        "GREATER_THAN",
							ComparisonValue: "20",
						},
						},
					},
				},
				Dimensions: []*gaRep.Dimension{
					{Name: "ga:country"},
					{Name: "ga:city"},
				},
			},
		},
	}
	return svc.Reports.BatchGet(req).Do()
}

// printResponse parses and prints the Analytics Reporting API V4 response.
func printResponse(res *gaRep.GetReportsResponse) {
	for _, report := range res.Reports {
		header := report.ColumnHeader
		dimHdrs := header.Dimensions
		metricHdrs := header.MetricHeader.MetricHeaderEntries
		rows := report.Data.Rows

		if rows == nil {
			fmt.Println("No data found for given view.")
		}

		for i := 0; i < len(dimHdrs); i++ {
			fmt.Printf(" %s ", dimHdrs[i])
		}

		for j := 0; j < len(metricHdrs); j++ {
			fmt.Printf(" %s ", metricHdrs[j].Name)
		}
		fmt.Println()

		for _, row := range rows {
			dims := row.Dimensions
			metrics := row.Metrics

			for i := 0; i < len(dims); i++ {
				fmt.Printf(" %s ", dims[i])

			}

			for _, metric := range metrics {
				// We have only 1 date range in the example
				// So it'll always print "Date Range (0)"
				// log.Infof("Date Range (%d)", idx)
				for j := 0; j < len(metric.Values); j++ {
					fmt.Printf(" %s ", metric.Values[j])
				}
			}

			fmt.Println("")
		}
	}

}
