package command

import (
  "fmt"
  "os"
  "strings"
  "time"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/slack"
)

const (
  MESSAGE_TEXT = "*From*: %s (%s)\nFinish executing the command on the server"
  MESSAGE_ATTACHMENT_TEXT = "*Command:* `%s`\n*Start at:* %s\n*End at:* %s\n*Duration:* %d seconds\n*Exit code:* %d"
)

func Run(cmd string) {
  if len(cmd) == 0 {
    os.Exit(1)
  }

  if len(config.SLACK_TOKEN) == 0 {
    fmt.Printf("Environment variable not defined: SLACK_TOKEN\n")
    os.Exit(1)
  }

  fmt.Printf("==> Run command...\n")
  fmt.Printf("--> Press Ctrl+C to end.\n")
  start := current_timestamp()
  fmt.Printf("--> Start at: %s\n", start)
  fmt.Printf("--> Wait to finish command: %s\n", cmd)
  stdout, exitcode := common.ExecCommand(cmd)
  stdout = clear_stdout(stdout)
  fmt.Printf("--> Stdout: %s", stdout)
  fmt.Printf("--> Exit code: %d\n", exitcode)
  end  := current_timestamp()
  diff := duration(start, end)
  fmt.Printf("--> End at: %s\n", end)
  fmt.Printf("--> Duration: %d\n", diff)

  msg := &slack.Message{
    Text: fmt.Sprintf(MESSAGE_TEXT, common.Hostname(), common.IpAddress()),
    Channel: config.SLACK_CHANNEL,
  }
  msg.AddAttachment(&slack.Attachment{
    Color: get_color(exitcode),
    Text: fmt.Sprintf(MESSAGE_ATTACHMENT_TEXT, cmd, start, end, diff, exitcode),
  })

  fmt.Printf("--> Slack response code: %d\n", slack.Send(msg))
}

func current_timestamp() string {
  t := time.Now().UTC()
  return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func duration(start string, end string) int {
  parsed_start , _ := time.Parse("2006-01-02 15:04:05", start);
  parsed_end   , _ := time.Parse("2006-01-02 15:04:05", end);

  return int(parsed_end.Sub(parsed_start).Seconds())
}

func clear_stdout(stdout string) string {
  if strings.HasPrefix(stdout, "\n") == false {
    stdout = "\n" + stdout
  }
  if strings.HasSuffix(stdout, "\n") == false {
    stdout = stdout + "\n"
  }
  return stdout
}

func get_color(exitcode int) string {
    if exitcode != 0 {
      return "danger"
    }
    return "good"
}
