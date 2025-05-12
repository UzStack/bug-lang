package libs

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/UzStack/bug-lang/internal/runtime/types"
	_ "github.com/lib/pq"
)

func PsqlConnect() any {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, "postgres", "root", "db",
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Query(db *sql.DB, sql *types.StringValue) any {
	rows, err := db.Query(sql.Value)
	if err != nil {
		panic(err)
	}
	return rows
}

func FindAll(rows *sql.Rows) any {
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	var result []any
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))

		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			panic(err)
		}
		res := types.NewMap(make(map[string]any)).(*types.MapValue)
		for i, col := range columns {
			val := values[i]
			switch v := val.(type) {
			case []byte:
				res.Add(col, types.NewString(strings.TrimSpace(string(v))))
			case int:
				res.Add(col, types.NewInt(v))
			case int8:
				res.Add(col, types.NewInt(int(v)))
			case int64:
				res.Add(col, types.NewInt(int(v)))
			default:
				res.Add(col, types.NewString(strings.TrimSpace(v.(string))))
			}
		}
		result = append(result, res)
	}
	return types.NewArray(result)
}

func Find(rows *sql.Rows) any {
	return FindAll(rows).(*types.ArrayValue).Values[0]
}
