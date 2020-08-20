package psql

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"

	// PG driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	outErr "github.com/vitamin-nn/otus_anti_bruteforce/internal/error"
)

const (
	ConstraintViolationCode = "23"
	WhiteListTable          = "ip_white_list"
	BlackListTable          = "ip_black_list"
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
	return sr.addToList(ctx, inet, WhiteListTable)
}

func (sr *Psql) AddToBlackList(ctx context.Context, inet *net.IPNet) error {
	return sr.addToList(ctx, inet, BlackListTable)
}

func (sr *Psql) DeleteFromWhiteList(ctx context.Context, inet *net.IPNet) error {
	return sr.deleteFromList(ctx, inet, WhiteListTable)
}

func (sr *Psql) DeleteFromBlackList(ctx context.Context, inet *net.IPNet) error {
	return sr.deleteFromList(ctx, inet, BlackListTable)
}

func (sr *Psql) GetWhiteList(ctx context.Context) ([]*net.IPNet, error) {
	return sr.getNetList(ctx, WhiteListTable)
}

func (sr *Psql) GetBlackList(ctx context.Context) ([]*net.IPNet, error) {
	return sr.getNetList(ctx, BlackListTable)
}

func (sr *Psql) getNetList(ctx context.Context, tableName string) ([]*net.IPNet, error) {
	// здесь не нашел другого решения для того чтобы линтер не ругался на sprintf для вставки имени таблицы
	// пробовал испльзовать pgx.Identifier.Sanitize(), но линтер на это не обращает внимания
	sql := fmt.Sprintf("SELECT ip_network FROM %s", tableName) // nolint: gosec
	// теортически, здесь в базе может быть очень много настроек для белых/черных списков
	// и тогда, наверное, лучше возвращать ссылку на курсор,
	// но, во-первых, такое все же маловероятно,
	// а во вторых - в таком случае реализация с хранением в памяти может не подойти;
	// поэтому с курсором здесь не вижу смысла заморачиваться

	rows, err := sr.db.QueryContext(ctx, sql)
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
	sql := fmt.Sprintf("INSERT INTO %s(ip_network) VALUES($1) RETURNING id", tableName)
	_, err := sr.db.ExecContext(ctx, sql, inet)
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
	sql := fmt.Sprintf("DELETE FROM %s WHERE ip_network = $1", tableName)
	res, err := sr.db.ExecContext(ctx, sql, inet)
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
