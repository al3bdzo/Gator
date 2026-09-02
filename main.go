package main

import (
	"github.com/al3bdzo/Gator/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.ReadJson()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cfg.SetUser("al3bdzo")
	if err != nil {
		fmt.Println(err)
		return
	}
	rcfg, err := config.ReadJson()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rcfg.Db_url)
	fmt.Println(rcfg.Current_user_name)
}