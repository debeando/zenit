package aws

import (
	"fmt"
	"regexp"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type InputAWSDiscover struct{}

type MySQL struct {
	Hostname  string `yaml:"hostname"`
	DSN       string `yaml:"dsn"`
	Aurora    bool   `yaml:"aurora"`
	Overflow  bool   `yaml:"overflow"`
	Slave     bool   `yaml:"slave"`
	Status    bool   `yaml:"status"`
	Tables    bool   `yaml:"tables"`
	Variables bool   `yaml:"variables"`
}

var mpm = [3]string{
	"aurora",
	"aurora-mysql",
	"mysql",
}

func (l *InputAWSDiscover) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputAWSDiscover", map[string]interface{}{"error": err})
		}
	}()

	if !config.File.Inputs.AWSDiscover.Enable {
		log.Debug("InputAWSDiscover", map[string]interface{}{"message": "Is not enabled."})
		return
	}

	if len(config.File.General.AWSRegion) == 0 {
		log.Info("InputAWSDiscover", map[string]interface{}{"message": "Require to define aws_region"})
		return
	}

	if len(config.File.General.AWSAccessKeyID) == 0 {
		log.Info("InputAWSDiscover", map[string]interface{}{"message": "Require to define aws_access_key_id"})
		return
	}

	if len(config.File.General.AWSSecretAccessKey) == 0 {
		log.Info("InputAWSDiscover", map[string]interface{}{"message": "Require to define aws_secret_access_key"})
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.File.General.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			config.File.General.AWSAccessKeyID,
			config.File.General.AWSSecretAccessKey,
			"",
		),
	},
	)
	if err != nil {
		log.Error("InputAWSDiscover", map[string]interface{}{"error": err})
		return
	}

	svc := rds.New(sess)

	result, err := svc.DescribeDBInstances(nil)
	if err != nil {
		log.Error("InputAWSDiscover", map[string]interface{}{"error": err})
		return
	}

	for _, d := range result.DBInstances {
		if itemExists(mpm, aws.StringValue(d.Engine)) {
			hostname := aws.StringValue(d.DBInstanceIdentifier)

			if !existInputMySQL(hostname) && matchFilter(hostname) {
				m := MySQL{}
				m.Hostname = hostname
				m.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=3s",
					config.File.Inputs.AWSDiscover.Username,
					config.File.Inputs.AWSDiscover.Password,
					aws.StringValue(d.Endpoint.Address),
					aws.Int64Value(d.Endpoint.Port),
				)
				m.Aurora = config.File.Inputs.AWSDiscover.Plugins.Aurora
				m.Overflow = config.File.Inputs.AWSDiscover.Plugins.Overflow
				m.Slave = config.File.Inputs.AWSDiscover.Plugins.Slave
				m.Status = config.File.Inputs.AWSDiscover.Plugins.Status
				m.Tables = config.File.Inputs.AWSDiscover.Plugins.Tables
				m.Variables = config.File.Inputs.AWSDiscover.Plugins.Variables

				config.File.Inputs.MySQL = append(config.File.Inputs.MySQL, m)

				log.Debug("InputAWSDiscover", map[string]interface{}{
					"name":     aws.StringValue(d.DBInstanceIdentifier),
					"engine":   aws.StringValue(d.Engine),
					"username": config.File.Inputs.AWSDiscover.Username,
					"endpoint": aws.StringValue(d.Endpoint.Address),
					"port":     aws.Int64Value(d.Endpoint.Port),
				})
			}
		}
	}
}

func init() {
	inputs.Add("InputAWSDiscover", func() inputs.Input { return &InputAWSDiscover{} })
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
	for _, mysql := range config.File.Inputs.MySQL {
		if mysql.Hostname == hostname {
			return true
		}
	}
	return false
}

func matchFilter(hostname string) bool {
	if matched, err := regexp.MatchString(config.File.Inputs.AWSDiscover.Filter, hostname); err == nil && matched == true {
		return true
	}

	return false
}
