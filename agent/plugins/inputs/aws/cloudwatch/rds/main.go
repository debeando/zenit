package rds

import (
	"time"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/debeando/go-common/log"
)

type Plugin struct {
	Counter int64
}

var plugin = new(Plugin)

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	log.DebugWithFields(name, log.Fields{
		"enable":   cnf.Inputs.AWS.CloudWatch.Enable,
		"interval": cnf.Inputs.AWS.CloudWatch.Interval,
		"counter":  p.Counter,
	})

	if !cnf.Inputs.AWS.CloudWatch.Enable {
		return
	}

	if !p.isTimeToCollect(cnf.Inputs.AWS.CloudWatch.Interval) {
		return
	}

	if len(cnf.General.AWSRegion) == 0 {
		log.InfoWithFields(name, log.Fields{"message": "Require to define aws_region"})
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cnf.General.AWSRegion),
	})

	if err != nil {
		log.ErrorWithFields(name, log.Fields{"message": err})
		return
	}

	var a = metrics.Load()
	var svcCW = cloudwatch.New(sess)
	var svcRDS = rds.New(sess)

	result, err := svcRDS.DescribeDBInstances(nil)
	if err != nil {
		log.ErrorWithFields(name, log.Fields{"message": err})
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
			log.ErrorWithFields(name, log.Fields{"message": err})
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
					{Key: "cpu", Value: *resp.Datapoints[i].Average},
				},
			})

			log.DebugWithFields(name, log.Fields{
				"name": *d.DBInstanceIdentifier,
				"cpu":  *resp.Datapoints[i].Average,
			})
		}
	}
}

func (p *Plugin) isTimeToCollect(i int) bool {
	if p.Counter == 0 || int(time.Since(time.Unix(p.Counter, 0)).Seconds()) >= i {
		(*p).Counter = int64(time.Now().Unix())

		return true
	}

	return false
}

func init() {
	inputs.Add("InputAWSCloudWatchRDS", func() inputs.Input { return plugin })
}
