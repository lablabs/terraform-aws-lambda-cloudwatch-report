package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gocarina/gocsv"
)

var (
	region          string
	metricName      string
	metricDimension string
	metricNamespace string
	emailSource     string
	emailTarget     string
	startTime       time.Time
	endTime         time.Time
)

type request struct{}

type timeData struct {
	time.Time
}

func (time *timeData) MarshalCSV() (string, error) {
	return time.Time.Format("2006-01-02T15:04:05"), nil
}

type statusData struct {
	float64
}

func (status *statusData) MarshalCSV() (string, error) {
	str := fmt.Sprintf("%.0f", status.float64)
	return str, nil
}

type sample struct {
	Time  timeData   `csv:"timestamp"`
	Value statusData `csv:"status_code"`
	Url   string     `csv:"url"`
}

func handleRequest(ctx context.Context, req request) (string, error) {

	region = os.Getenv("REGION")
	metricName = os.Getenv("METRIC_NAME")
	metricNamespace = os.Getenv("METRIC_NAMESPACE")
	metricDimension = os.Getenv("METRIC_DIMENSION")
	emailSource = os.Getenv("EMAIL_SOURCE")
	emailTarget = os.Getenv("EMAIL_TARGET")

	endTime = time.Now()
	startTime = endTime.AddDate(0, 0, -1)

	sendEmail(getMetricData())
	return "Finished", nil
}

func buildEmailInput(source, destination, subject, message string,
	csvFile []byte) (*ses.SendRawEmailInput, error) {

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// email main header:
	h := make(textproto.MIMEHeader)
	h.Set("From", source)
	h.Set("To", destination)
	h.Set("Return-Path", source)
	h.Set("Subject", subject)
	h.Set("Content-Language", "en-US")
	h.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
	h.Set("MIME-Version", "1.0")
	_, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// body:
	h = make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "7bit")
	h.Set("Content-Type", "text/plain; charset=us-ascii")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte(message))
	if err != nil {
		return nil, err
	}

	fn := "report_" + startTime.Format("2006-01-02T15-04-05") + "_" + endTime.Format("2006-01-02T15-04-05") + ".csv"
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition", "attachment; filename="+fn)
	h.Set("Content-Type", "text/csv; x-unix-mode=0644; name=\""+fn+"\"")
	h.Set("Content-Transfer-Encoding", "7bit")
	part, err = writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(csvFile)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Strip boundary line before header (doesn't work with it present)
	s := buf.String()
	if strings.Count(s, "\n") < 2 {
		return nil, fmt.Errorf("invalid e-mail content")
	}
	s = strings.SplitN(s, "\n", 2)[1]

	raw := ses.RawMessage{
		Data: []byte(s),
	}
	input := &ses.SendRawEmailInput{
		Destinations: []*string{aws.String(destination)},
		Source:       aws.String(source),
		RawMessage:   &raw,
	}

	return input, nil
}

func sendEmail(attachement []byte) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		log.Fatal("Error creating session")
	}

	svc := ses.New(sess)

	input, err := buildEmailInput(
		emailSource,
		emailTarget,
		"Healthcheck report for "+metricDimension,
		"Report for "+metricDimension+" is attached to this message\n"+"from "+startTime.Format("2006-01-02T15:04:05")+"\n"+"to  "+endTime.Format("2006-01-02T15:04:05")+"\n",
		attachement,
	)

	if err != nil {
		panic(err)
	}

	_, err = svc.SendRawEmail(input)
	if err != nil {
		panic(err)
	}
}

func getMetricData() []byte {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		log.Fatal("Error creating session")
	}

	svc := cloudwatch.New(sess)

	metricDataQuery := cloudwatch.MetricDataQuery{
		Id: aws.String("result"),
		MetricStat: &cloudwatch.MetricStat{
			Period: aws.Int64(300),
			Stat:   aws.String("Maximum"),
			Metric: &cloudwatch.Metric{
				MetricName: aws.String(metricName),
				Namespace:  aws.String(metricNamespace),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("HealthcheckTarget"),
						Value: aws.String(metricDimension),
					},
				},
			},
		},
		ReturnData: aws.Bool(true),
	}

	result, err := svc.GetMetricData(&cloudwatch.GetMetricDataInput{
		StartTime: &startTime,
		EndTime:   &endTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			&metricDataQuery,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	var samples []*sample

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		samples = append([]*sample{&sample{Time: timeData{*timestamp}, Value: statusData{*result.MetricDataResults[0].Values[i]}, Url: metricDimension}}, samples...)
	}

	csvContent, err := gocsv.MarshalBytes(&samples)
	if err != nil {
		panic(err)
	}

	return csvContent
}

func main() {
	// handleRequest(nil, request{})
	lambda.Start(handleRequest)
}
