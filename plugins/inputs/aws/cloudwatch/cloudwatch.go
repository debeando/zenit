package cloudwatch

import (
	"time"

	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"
)

type InputAWSCloudWatch struct{}

func (l *InputAWSCloudWatch) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputAWSCloudWatch", map[string]interface{}{"error": err})
		}
	}()

	if !config.File.Inputs.AWSCloudWatch.Enable {
		log.Debug("InputAWSCloudWatch", map[string]interface{}{"message": "Is not enabled."})
		return
	}

	if len(config.File.General.AWSRegion) == 0 {
		log.Info("InputAWSCloudWatch", map[string]interface{}{"message": "Require to define aws_region"})
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.File.General.AWSRegion),
	})

	if err != nil {
		log.Error("InputAWSCloudWatch", map[string]interface{}{"error": err})
		return
	}

	var a = metrics.Load()
	var svcCW = cloudwatch.New(sess)
	var svcRDS = rds.New(sess)

	result, err := svcRDS.DescribeDBInstances(nil)
	if err != nil {
		log.Error("InputAWSCloudWatch", map[string]interface{}{"error": err})
		return
	}

	for _, d := range result.DBInstances {
		search := cloudwatch.GetMetricStatisticsInput{
			StartTime:  aws.Time(time.Now().UTC().Add(time.Second * -300)),
			EndTime:    aws.Time(time.Now().UTC()),
			MetricName: aws.String("CPUUtilization"),
			Period:     aws.Int64(60),
			Statistics: []*string{aws.String("Average")},
			Namespace:  aws.String("AWS/RDS"),
			Dimensions: []*cloudwatch.Dimension{{
				Name:  aws.String("DBInstanceIdentifier"),
				Value: d.DBInstanceIdentifier,
			}},
		}

		resp, err := svcCW.GetMetricStatistics(&search)
		if err != nil {
			log.Error("InputAWSCloudWatch", map[string]interface{}{"error": err})
			return
		}

		if len(resp.Datapoints) > 0 {
			i := len(resp.Datapoints)
			if i > 0 {
				i--
			}

			a.Add(metrics.Metric{
				Key: "aws_cloudwatch_rds",
				Tags: []metrics.Tag{
					{Name: "name", Value: *d.DBInstanceIdentifier},
				},
				Values: []metrics.Value{
					{Key: "percentage", Value: *resp.Datapoints[i].Average},
				},
			})
		}
	}
}

func init() {
	inputs.Add("InputAWSCloudWatch", func() inputs.Input { return &InputAWSCloudWatch{} })
}
