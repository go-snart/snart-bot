package token

import (
	"context"
	"fmt"

	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/logs"
)

func table(ctx context.Context, d *db.DB) {
	const (
		e = `CREATE TABLE IF NOT EXISTS token(
			value TEXT
		)`
	)

	_, err := d.Conn(&ctx).Exec(ctx, e)
	if err != nil {
		err = fmt.Errorf("exec %#q: %w", e, err)

		logs.Warn.Println(err)

		return
	}
}

// SelectTokens retrieves bot tokens from a DB.
func SelectTokens(ctx context.Context, d *db.DB) ([]string, error) {
	logs.Debug.Println("enter->table")

	table(ctx, d)

	logs.Debug.Println("table->query")

	const q = `SELECT value FROM token`

	rows, err := d.Conn(&ctx).Query(ctx, q)
	if err != nil {
		err = fmt.Errorf("query %#q: %w", q, err)

		logs.Warn.Println(err)

		return nil, err
	}
	defer rows.Close()

	logs.Debug.Println("query->scan")

	toks := []string(nil)

	for rows.Next() {
		tok := ""

		err = rows.Scan(&tok)
		if err != nil {
			err = fmt.Errorf("scan tok: %w", err)

			logs.Warn.Println(err)

			return nil, err
		}

		toks = append(toks, tok)
	}

	logs.Debug.Println("scan->err")

	if err := rows.Err(); err != nil {
		err = fmt.Errorf("rows: %w", err)

		logs.Warn.Println(err)

		return nil, err
	}

	logs.Debug.Println("err->done")

	return toks, nil
}

// InsertTokens adds tokens to the database so that they're persistent.
func InsertTokens(ctx context.Context, d *db.DB, toks []string) {
	table(ctx, d)

	e := `INSERT INTO token(value) VALUES`
	vals := []interface{}(nil)

	for n, tok := range toks {
		e += fmt.Sprintf(`($%d)`, n)

		if n < len(toks)-1 {
			e += `,`
		}

		vals = append(vals, tok)
	}

	_, err := d.Conn(&ctx).Exec(ctx, e, vals...)
	if err != nil {
		err = fmt.Errorf("exec %#q (%#v): %w", e, vals, err)

		logs.Warn.Println(err)

		return
	}
}
