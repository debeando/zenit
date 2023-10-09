package rds

import (
	"fmt"
	"regexp"

	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/debeando/go-common/log"
)

type Plugin struct{}

var mpm = [3]string{
	"aurora",
	"aurora-mysql",
	"mysql",
}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	if !cnf.Inputs.AWS.Discover.Enable {
		log.DebugWithFields(name, log.Fields{"message": "Is not enabled."})
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

	svc := rds.New(sess)

	result, err := svc.DescribeDBInstances(nil)
	if err != nil {
		log.ErrorWithFields(name, log.Fields{"message": err})
		return
	}

	for _, d := range result.DBInstances {
		if itemExists(mpm, aws.StringValue(d.Engine)) {
			hostname := aws.StringValue(d.DBInstanceIdentifier)

			if !existInputMySQL(hostname) && matchFilter(hostname) {
				m := config.MySQL{}
				m.Engine = aws.StringValue(d.Engine)
				m.Hostname = hostname
				m.Enable = cnf.Inputs.AWS.Discover.Plugins.MySQL.Enable
				m.Overflow = cnf.Inputs.AWS.Discover.Plugins.MySQL.Overflow
				m.Replica = cnf.Inputs.AWS.Discover.Plugins.MySQL.Replica
				m.Status = cnf.Inputs.AWS.Discover.Plugins.MySQL.Status
				m.Tables = cnf.Inputs.AWS.Discover.Plugins.MySQL.Tables
				m.Variables = cnf.Inputs.AWS.Discover.Plugins.MySQL.Variables
				m.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=3s",
					cnf.Inputs.AWS.Discover.Username,
					cnf.Inputs.AWS.Discover.Password,
					aws.StringValue(d.Endpoint.Address),
					aws.Int64Value(d.Endpoint.Port),
				)

				if m.Engine == "aurora-mysql" {
					m.Aurora = true
				}

				cnf.Inputs.MySQL = append(cnf.Inputs.MySQL, m)

				log.DebugWithFields(name, log.Fields{
					"dsn":       m.DSN,
					"enable":    m.Enable,
					"engine":    m.Engine,
					"hostname":  m.Hostname,
					"overflow":  m.Overflow,
					"replica":   m.Replica,
					"status":    m.Status,
					"tables":    m.Tables,
					"variables": m.Variables,
				})
			}
		}
	}
}

func init() {
	inputs.Add("InputAWSDiscoverRDS", func() inputs.Input { return &Plugin{} })
}

func itemExists(items [3]string, str string) bool {
	for _, i := range items {
		if i == str {
			return true
		}
	}
	return false
}

func existInputMySQL(hostname string) bool {
	for _, mysql := range config.GetInstance().Inputs.MySQL {
		if mysql.Hostname == hostname {
			return true
		}
	}
	return false
}

func matchFilter(hostname string) bool {
	if matched, err := regexp.MatchString(config.GetInstance().Inputs.AWS.Discover.Filter, hostname); err == nil && matched == true {
		return true
	}

	return false
}
