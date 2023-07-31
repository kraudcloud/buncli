package main

import (
	"github.com/uptrace/bun"
)

type Zone struct {
	bun.BaseModel `bun:"table:zones"`

	ID   string `bun:"type:uuid,pk,nullzero"`
	Name string `bun:"name,nullzero,notnull,unique"`
}

type Rack struct {
	bun.BaseModel `bun:"table:racks"`

	ID   string `bun:"type:uuid,pk,nullzero"`
	Name string `bun:",nullzero,notnull,unique"`

	ZoneID string `bun:"type:uuid,notnull"`
	Zone   *Zone  `bun:"rel:belongs-to"`
}
