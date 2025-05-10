package libs

import (
	"database/sql"
	"fmt"

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
		res := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			res[col] = val
		}
		result = append(result, res)
	}
	return result
}
