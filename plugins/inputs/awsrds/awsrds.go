package awsrds

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type InputAWSRDS struct{}

type MySQL struct {
	Hostname  string `yaml:"hostname"`
	DSN       string `yaml:"dsn"`
	Aurora    bool `yaml:"aurora"`
	Overflow  bool `yaml:"overflow"`
	Slave     bool `yaml:"slave"`
	Status    bool `yaml:"status"`
	Tables    bool `yaml:"tables"`
	Variables bool `yaml:"variables"`
}

func (l *InputAWSRDS) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputAWSRDS - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Inputs.AWSRDS.Enable {
		return
	}

	if len(config.File.General.AWSRegion) == 0 {
		log.Info("Plugin - InputAWSRDS - Require to define aws_region")
		return
	}

	if len(config.File.General.AWSAccessKeyID) == 0 {
		log.Info("Plugin - InputAWSRDS - Require to define aws_access_key_id")
		return
	}

	if len(config.File.General.AWSSecretAccessKey) == 0 {
		log.Info("Plugin - InputAWSRDS - Require to define aws_secret_access_key")
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

	svc := rds.New(sess)

	result, err := svc.DescribeDBInstances(nil)
	if err != nil {
		log.Debug(fmt.Sprintf("Plugin - InputAWSRDS - Unable to list instances."))
		return
	}

	for _, d := range result.DBInstances {
		if ! existInputMySQL(aws.StringValue(d.DBInstanceIdentifier)) {
			log.Debug(fmt.Sprintf("Plugin - InputAWSRDS - Register new server: %s", aws.StringValue(d.DBInstanceIdentifier)))

			m := MySQL{}
			m.Hostname = aws.StringValue(d.DBInstanceIdentifier)
			m.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/",
				config.File.Inputs.AWSRDS.Username,
				config.File.Inputs.AWSRDS.Password,
				aws.StringValue(d.Endpoint.Address),
				aws.Int64Value(d.Endpoint.Port),
			)
			m.Aurora    = config.File.Inputs.AWSRDS.Plugins.Aurora
			m.Overflow  = config.File.Inputs.AWSRDS.Plugins.Overflow
			m.Slave     = config.File.Inputs.AWSRDS.Plugins.Slave
			m.Status    = config.File.Inputs.AWSRDS.Plugins.Status
			m.Tables    = config.File.Inputs.AWSRDS.Plugins.Tables
			m.Variables = config.File.Inputs.AWSRDS.Plugins.Variables

			config.File.Inputs.MySQL = append(config.File.Inputs.MySQL, m)
		}
	}
}

func existInputMySQL(hostname string) bool {
	for _, mysql := range config.File.Inputs.MySQL {
		if mysql.Hostname == hostname {
			return true
		}
	}
	return false
}

func init() {
	inputs.Add("InputAWSRDS", func() inputs.Input { return &InputAWSRDS{} })
}
