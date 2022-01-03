package main

import (
	"os"
	"database/sql"
)

const (
	NAME = iota + 1
	PLAT
	GOLD
	SILVER
	COPPER
	ITEM = 8
)

var DB *sql.DB
var ZONE string

func main() {

	InitDB()

	GetZone()

	file := "/home/namll/.wine/drive_c/Program Files/Sony/EverQuest/Logs/eqlog_Oxnull_P1999Green.txt"

	f, _ := os.Create(file)
	f.Close()

	log := GetLog(file)

	if log == nil { return }

	defer log.Close()

	go GetInput()

	for{ ParseLog(log) }
}
