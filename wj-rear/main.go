package main

import (
	"wj_rear/database"
	"wj_rear/router"
)

func main() {
	database.Init()
	router.Run()
}
