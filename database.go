package main

import (
	"fmt"
	"strings"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type item_t struct {
	id int
	name string
}

func InitDB(){
	DB, _ = sql.Open("sqlite3", "./eqparse.sqlite")
	itemTableSQL := `CREATE TABLE IF NOT EXISTS item (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT UNIQUE NOT NULL
		);`

	DB.Exec(itemTableSQL)

	buyItemTableSQL := `CREATE TABLE IF NOT EXISTS item_buy (
		"id" integer NOT NULL UNIQUE,
		"platinum" integer default null,
		"gold" integer default null,
		"silver" integer default null,
		"copper" integer default null
		);`

	DB.Exec(buyItemTableSQL)

	sellItemTableSQL := `CREATE TABLE IF NOT EXISTS item_sell (
		"id" integer NOT NULL UNIQUE,
		"platinum" integer default null,
		"gold" integer default null,
		"silver" integer default null,
		"copper" integer default null
		);`

	DB.Exec(sellItemTableSQL)

	NPCTableSQL := `CREATE TABLE IF NOT EXISTS npc (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT UNIQUE NOT NULL,
		"loc" TEXT NOT NULL
		);`

	DB.Exec(NPCTableSQL)

	NPCItemsTableSQL := `CREATE TABLE IF NOT EXISTS npc_item (
		"npc_id" integer NOT NULL,
		"item_id" integer NOT NULL
		);`

	DB.Exec(NPCItemsTableSQL)

	NPCQuestTableSQL := `CREATE TABLE IF NOT EXISTS npc_quest (
		"npc_id" integer NOT NULL,
		"tag" TEXT NOT NULL
		);`

	DB.Exec(NPCQuestTableSQL)
}

func WriteNPCToDB(npc string) int64{
	query := `INSERT INTO npc (name, loc) VALUES ($1, $2)`
	res, err := DB.Exec(query, npc, ZONE)

	var id int64
	if err != nil {
		query := `SELECT id FROM npc where name=$1`
		rows, _ := DB.Query(query, npc)

		rows.Next()
		rows.Scan(&id)
		rows.Close()
	}else{
		id, _ = res.LastInsertId()
	}

	return id
}

func WriteBuyToDB(parts [][]string){
	query := `INSERT INTO item (name) VALUES ($1)`
	res, err := DB.Exec(query, parts[0][ITEM])

	var id int64
	if err != nil {
		query := `SELECT id FROM item where name=$1`
		rows, _ := DB.Query(query, parts[0][ITEM])

		rows.Next()
		rows.Scan(&id)
		rows.Close()
	}else{
		id, _ = res.LastInsertId()
	}

	query = `SELECT id FROM item_buy WHERE id=$1`
	rows, _ := DB.Query(query, id)

	if !rows.Next() {
		query = `INSERT INTO item_buy (id, platinum, gold, silver, copper) VALUES ($1, $2, $3, $4, $5)`
		DB.Exec(query, id, parts[0][PLAT], parts[0][GOLD], parts[0][SILVER], parts[0][COPPER])


	}

	rows.Close()

	npcID := WriteNPCToDB(parts[0][NAME])
	query = `SELECT item_id FROM npc_item WHERE npc_id=$1 and item_id=$2`

	rows, _ = DB.Query(query, npcID, id)

	if !rows.Next() {
		query = `INSERT INTO npc_item (npc_id, item_id) VALUES ($1, $2)`
		DB.Exec(query, npcID, id)
	}

	rows.Close()
}

func WriteSellToDB(parts [][]string){
	query := `INSERT INTO item (name) VALUES ($1)`
	res, err := DB.Exec(query, parts[0][ITEM])

	var id int64
	if err != nil {
		query := `SELECT id FROM item where name=$1`
		rows, _ := DB.Query(query, parts[0][ITEM])

		rows.Next()
		rows.Scan(&id)
		rows.Close()
	}else{
		id, _ = res.LastInsertId()
	}

	query = `INSERT INTO item_sell (id, platinum, gold, silver, copper) VALUES ($1, $2, $3, $4, $5)`
	DB.Exec(query, id, parts[0][PLAT], parts[0][GOLD], parts[0][SILVER], parts[0][COPPER])
}

func GetItemDB(item string){
	query := `SELECT id, name FROM item WHERE LOWER (name) LIKE $1`

	item = strings.ReplaceAll(item, " ", "%")
	item = "%"+item+"%"
	item = strings.ToLower(item)

	var items []item_t
	rows, _ := DB.Query(query, item)
	defer rows.Close()
	for rows.Next(){
		var item item_t
		rows.Scan(&item.id, &item.name)
		items = append(items, item)
	}
	rows.Close()

	for _, e := range items {
		var bPlat, bGold, bSilver, bCopper int
		query = `SELECT platinum, gold, silver, copper FROM item_buy WHERE id=$1`
		rows, _ = DB.Query(query, e.id)

		if rows.Next(){
			rows.Scan(&bPlat, &bGold, &bSilver, &bCopper)
		}

		rows.Close()

		var sPlat, sGold, sSilver, sCopper int
		query = `SELECT platinum, gold, silver, copper FROM item_sell WHERE id=$1`
		rows, _ = DB.Query(query, e.id)

		if rows.Next(){
			rows.Scan(&sPlat, &sGold, &sSilver, &sCopper)
		}

		rows.Close()

		fmt.Println("-----")
		fmt.Printf("%s\nBuy: %d P %d G %d S %d C\nSell %d P %d G %d S %d C\n", e.name, bPlat, bGold, bSilver, bCopper, sPlat, sGold, sSilver, sCopper)
	}


}
