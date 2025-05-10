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
	row, err := db.Query(sql.Value)
	if err != nil {
		panic(err.Error())
	}
	return row
}
