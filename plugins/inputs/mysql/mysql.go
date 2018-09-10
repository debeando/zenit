package mysql

import (
	"log"
	"strings"

	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
)

func Check() bool {
	log.Printf("I! - MySQL - DSN: %s\n", config.File.MySQL.DSN)
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	if err != nil {
		log.Printf("E! - MySQL - Impossible to connect: %s\n", err)
		return false
	}

	log.Println("I! - MySQL - Connected successfully.")
	conn.Close()
	return true
}

func ClearUser(u string) string {
	index := strings.Index(u, "[")
	if index > 0 {
		return u[0:index]
	}
	return u
}
