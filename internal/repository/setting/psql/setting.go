package psql

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	outErr "github.com/vitamin-nn/otus_anti_bruteforce/internal/error"
)

const (
	ConstraintViolationCode = "23"
)

type Psql struct {
	dsn string
	db  *sql.DB
}

func NewSettingRepo(dsn string) *Psql {
	return &Psql{
		dsn: dsn,
	}
}

func (sr *Psql) Connect(ctx context.Context) error {
	db, err := sql.Open("pgx", sr.dsn)
	if err != nil {
		return err
	}
	sr.db = db
	sr.db.Stats()
	return sr.db.PingContext(ctx)
}

func (sr *Psql) Close() error {
	return sr.db.Close()
}

func (sr *Psql) AddToWhiteList(ctx context.Context, inet *net.IPNet) error {
	return sr.addToList(ctx, inet, "ip_white_list")
}

func (sr *Psql) AddToBlackList(ctx context.Context, inet *net.IPNet) error {
	return sr.addToList(ctx, inet, "ip_black_list")
}

func (sr *Psql) addToList(ctx context.Context, inet *net.IPNet, tableName string) error {
	_, err := sr.db.ExecContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s(ip_network) VALUES($1) RETURNING id", tableName),
		inet.String(),
	)
	if err != nil {
		specErr := getSpecificError(err)
		if specErr == nil {
			specErr = fmt.Errorf("insert error: %v", err)
		}
		return specErr
	}

	return nil
}

func getSpecificError(err error) error {
	if errPg, ok := err.(*pgconn.PgError); ok {
		if sqlState := errPg.SQLState(); len(sqlState) > 2 && sqlState[0:2] == ConstraintViolationCode {
			return outErr.ErrINetAlreadySet
		}
	}
	return nil
}
