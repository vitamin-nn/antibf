package psql

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	_ "github.com/jackc/pgx/v4/stdlib"
	outErr "github.com/vitamin-nn/otus_anti_bruteforce/internal/error"
)

const (
	ConstraintViolationCode = "23"
)

var settingTableList = map[string]string{
	"white": "ip_white_list",
	"black": "ip_black_list",
}

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
	return sr.addToList(ctx, inet, getWhiteTableName())
}

func (sr *Psql) AddToBlackList(ctx context.Context, inet *net.IPNet) error {
	return sr.addToList(ctx, inet, getBlackTableName())
}

func (sr *Psql) DeleteFromWhiteList(ctx context.Context, inet *net.IPNet) error {
	return sr.deleteFromList(ctx, inet, getWhiteTableName())
}

func (sr *Psql) DeleteFromBlackList(ctx context.Context, inet *net.IPNet) error {
	return sr.deleteFromList(ctx, inet, getBlackTableName())
}

func (sr *Psql) GetWhiteList(ctx context.Context) ([]*net.IPNet, error) {
	return sr.getNetList(ctx, getWhiteTableName())
}

func (sr *Psql) GetBlackList(ctx context.Context) ([]*net.IPNet, error) {
	return sr.getNetList(ctx, getBlackTableName())
}

func (sr *Psql) getNetList(ctx context.Context, tableName string) ([]*net.IPNet, error) {
	// теортически, здесь в базе может быть очень много настроек для белых/черных списков
	// и тогда, наверное, лучше возвращать ссылку на курсор,
	// но, во-первых, такое все же маловероятно,
	// а во вторых - в таком случае реализация с хранением в памяти не подходит;
	// поэтому с курсором здесь не вижу смысла заморачиваться
	rows, err := sr.db.QueryContext(ctx, fmt.Sprintf("SELECT ip_network FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var netList []*net.IPNet

	for rows.Next() {
		var net pgtype.Inet

		if err := rows.Scan(
			&net,
		); err != nil {
			return nil, err
		}

		netList = append(netList, net.IPNet)
	}
	return netList, rows.Err()
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

func (sr *Psql) deleteFromList(ctx context.Context, inet *net.IPNet, tableName string) error {
	res, err := sr.db.ExecContext(
		ctx,
		fmt.Sprintf("DELETE FROM %s WHERE ip_network = $1", tableName),
		inet,
	)
	if err != nil {
		return fmt.Errorf("delete from events error: %v", err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt < 1 {
		return outErr.ErrINetNotFound
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

func getWhiteTableName() string {
	return settingTableList["white"]
}

func getBlackTableName() string {
	return settingTableList["black"]
}
