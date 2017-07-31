package main

import (
	"config"
	"fmt"
)

func main() {

	config.WriteKey("/Users/ghc/IdeaProjects/dbMag/src/etc/my.cnf", "client", "max_packet_size", "32M")

	size := config.ReadSectionKey("/users/ghc/IdeaProjects/dbMag/src/etc/my.cnf", "client", "max_packet_size")
	fmt.Println(size)
	//config.WriteSection("/tmp/my.cnf","client")
}
