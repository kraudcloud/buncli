package main

import (
	"context"
	"os"
	"strings"
	"time"

	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"

	//"github.com/uptrace/bun/extra/bundebug"

	"github.com/oiime/logrusbun"
	logrus "github.com/sirupsen/logrus"
)

const DEBUG = true

var DB *bun.DB

func DBInit() {

	DB = bun.NewDB(
		sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(strings.ReplaceAll(os.Getenv("DATABASE_URL"), "cockroachdb", "postgres")))),
		pgdialect.New(),
		bun.WithDiscardUnknownColumns(),
	)

	DB.AddQueryHook(bunotel.NewQueryHook())

	if DEBUG {
		DB.AddQueryHook(logrusbun.NewQueryHook(logrusbun.QueryHookOptions{
			LogSlow:    200 * time.Millisecond,
			QueryLevel: logrus.DebugLevel,
			ErrorLevel: logrus.ErrorLevel,
			SlowLevel:  logrus.WarnLevel,
			Logger:     log,
		}))
		//DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	err := DB.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS zones (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL UNIQUE
		)`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS racks (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL UNIQUE,
			zone_id UUID REFERENCES zones(id)
		)`)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
