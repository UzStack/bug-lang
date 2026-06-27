package libs

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"      // postgres driver
	_ "modernc.org/sqlite"     // sqlite driver (pure Go, CGO kerak emas)

	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

// DBConnect ulanish ochadi.
//   postgres: driver "postgres", dsn "host=... dbname=... user=..."
//   sqlite:   driver "sqlite",   dsn "fayl.db" (yoki ":memory:")
func DBConnect(driver *types.StringValue, dsn *types.StringValue) any {
	conn, err := sql.Open(driver.Value, dsn.Value)
	if err != nil {
		panic("db.connect: " + err.Error())
	}
	if err := conn.Ping(); err != nil {
		panic("db.connect: " + err.Error())
	}
	return conn // xom *sql.DB handle — buglang uni o'zgaruvchida saqlaydi
}

// DBQuery parametrli SELECT. params $1, $2 ... bilan bog'lanadi (SQL injection yo'q).
func DBQuery(conn *sql.DB, query *types.StringValue, params *types.ArrayValue) any {
	rows, err := conn.Query(query.Value, unwrapArgs(params)...)
	if err != nil {
		panic("db.query: " + err.Error())
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		panic("db.query: " + err.Error())
	}
	var result []any
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			panic("db.query: " + err.Error())
		}
		row := make(map[string]any, len(cols))
		for i, c := range cols {
			row[c] = normalizeDB(vals[i])
		}
		result = append(result, utils.DecodeBug(row))
	}
	return types.NewArray(result)
}

// DBExec parametrli INSERT/UPDATE/DELETE. Ta'sirlangan qatorlar sonini qaytaradi.
func DBExec(conn *sql.DB, query *types.StringValue, params *types.ArrayValue) any {
	res, err := conn.Exec(query.Value, unwrapArgs(params)...)
	if err != nil {
		panic("db.exec: " + err.Error())
	}
	n, _ := res.RowsAffected()
	return types.NewInt(int(n))
}

func DBClose(conn *sql.DB) any {
	if err := conn.Close(); err != nil {
		panic("db.close: " + err.Error())
	}
	return types.NewNull()
}

// unwrapArgs buglang list elementlarini Go primitivlariga aylantiradi.
func unwrapArgs(params *types.ArrayValue) []any {
	if params == nil {
		return nil
	}
	out := make([]any, len(params.Values))
	for i, a := range params.Values {
		out[i] = utils.EncodeBug(a)
	}
	return out
}

// normalizeDB driver qaytargan tiplarni buglang uchun qulaylashtiradi.
func normalizeDB(v any) any {
	switch t := v.(type) {
	case []byte:
		return string(t)
	case time.Time:
		return t.Format("2006-01-02 15:04:05")
	}
	return v
}
