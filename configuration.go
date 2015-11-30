package main

import "math/rand"

type conf struct {
	// the auto-increment-increment on the mysql server
	// note all servers must be equal
	increment int
	// the max amount of ids that can be requested
	maxAmount int
	// table name
	tableName string
	// source unix: user:pass@unix(/full-path-to-sock-file)/database
	// source tcp: user:pass@tcp(host:port)/database
	hosts map[string]string
}

func (c conf) getRandomHost() string {
	i := int(float32(len(c.hosts)) * rand.Float32())
	for _, v := range c.hosts {
		if i == 0 {
			return v
		} else {
			i--
		}
	}
	panic("Not able to return a random host")
}
