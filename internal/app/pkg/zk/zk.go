package zk

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/samuel/go-zookeeper/zk"
)

func NewZkConn() (*zk.Conn, error) {
	services := []string{"127.0.0.1:2181"}
	username := ""
	passwd := ""

	if v, ok := os.LookupEnv(constant.EnvNameZkServices); ok {
		services = strings.Split(v, ",")
	}
	if v, ok := os.LookupEnv(constant.EnvNameZkUsername); ok {
		username = v
	}
	if v, ok := os.LookupEnv(constant.EnvNameZkPasswd); ok {
		passwd = v
	}

	conn, _, err := zk.Connect(services, 10*time.Second)
	if err != nil {
		return nil, err
	}
	if username != "" {
		authData := fmt.Sprintf("%s:%s", username, passwd)
		if err := conn.AddAuth("digest", []byte(authData)); err != nil {
			return nil, err
		}
	}
	return conn, nil
}
