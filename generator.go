package main

import (
	"bytes"
	"database/sql"
	"fmt"
)

type idGenerator struct {
	db        *sql.DB
	tableName string
	increment int
}

func newIdGenerator(host string, tableName string, increment int) *idGenerator {
	db, err := sql.Open("mysql", host)
	if err != nil {
		panic(err)
	}

	return &idGenerator{db: db, tableName: tableName, increment: increment}
}

func (idG idGenerator) get(amount int) []int64 {

	getSqlValues := func(amount int) string {
		var buffer bytes.Buffer
		for i := 1; i <= amount; i++ {
			buffer.WriteString("('a'),")
		}
		values := buffer.String()

		return values[:len(values)-1]
	}

	calculateIds := func(firstId int64) []int64 {
		id := firstId
		ids := make([]int64, amount)
		for i := 0; i < amount; i++ {
			ids[i] = id
			id = id + int64(idG.increment)
		}

		return ids
	}

	result, err := idG.db.Exec(
		fmt.Sprintf("REPLACE INTO %s (stub) VALUES %s", idG.tableName, getSqlValues(amount)),
	)
	if err != nil {
		panic(err)
	}
	firstId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	return calculateIds(firstId)
}
