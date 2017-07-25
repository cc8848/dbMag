package main

import (
	"fmt"
	"replication"
	"os"
	"mysql"
	"github.com/juju/errors"
	"golang.org/x/net/context"
)
func main() {

	cfg := replication.BinlogSyncerConfig{
		ServerID: 101,
		Flavor:   "mariadb",

		Host:            "180.97.81.42",
		Port:            33068,
		User:            "repl",
		Password:        "123",
		RawModeEnabled:  false,
		SemiSyncEnabled: false,
	}

	b := replication.NewBinlogSyncer(&cfg)

	//show master info 构造一个位置
	pos := mysql.Position{"mybinlog.000002", 4}

	s, err := b.StartSync(pos)
	if err != nil {
		fmt.Printf("Start sync error: %v\n", errors.ErrorStack(err))
		return
	}

	for {
		e, err := s.GetEvent(context.Background())
		if err != nil {
			fmt.Printf("Get event error: %v\n", errors.ErrorStack(err))
			return
		}

		e.Dump(os.Stdout)
	}

	fmt.Println("hello world!")
}