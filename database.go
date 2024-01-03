package sqlgodb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	// "os"
	// "strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	// "github.com/go-sql-driver/mysql"
	// "github.com/joho/godotenv"
	// "strings"
)

var dbInstace *sql.DB

func StartDB() {
	if dbInstace == nil {
		dbInstace = initDBConnection(&mysql.Config{User: "root", Passwd: "", Net: "tcp", Addr: "localhost:3306", DBName: "aago_care", AllowNativePasswords: true})
	}
}

func initDBConnection(config *mysql.Config) *sql.DB {
	db, err := sql.Open("mysql", config.FormatDSN())
	// db, err := sql.Open("mysql", "root:@(localhost:3306)/aago_care")

	if err != nil {
		log.Fatalf("Impossible to create the create the connection: %s", err)
		panic(err)
	} else {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	}
	return db
}

func InsertIntoTable(table string, columns []string, values []string) *DBResult {
	queryStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, strings.Join(columns, ","), insertStmtStr(len(values)))

	stmt, prepareError := dbInstace.Prepare(queryStr)
	if prepareError != nil {
		return &DBResult{
			Err: prepareError,
		}
	} else {
		defer stmt.Close()
	}
	newValues := arrayToInterfaceArr[string](values)

	insertResult, insertErr := stmt.Exec(newValues...)
	if insertErr != nil {
		return &DBResult{
			Err: insertErr,
		}
	}
	lastInsertId, lastIdError := insertResult.LastInsertId()
	if lastIdError != nil {
		return &DBResult{
			Err: lastIdError,
		}
	}
	return &DBResult{
		Err:          nil,
		LastInsertId: int(lastInsertId),
	}
}

func QueryData(table string, whereClause string, values []string) *DBResult {
	queryStr := fmt.Sprintf("SELECT * FROM %s", table)
	queryStr = fmt.Sprintf("%s %s", queryStr, whereClause)

	newValues := arrayToInterfaceArr(values)

	stmt, prepareError := dbInstace.Prepare(queryStr)

	if prepareError != nil {
		return &DBResult{
			Err: prepareError,
		}
	} else {
		defer stmt.Close()
	}

	rows, queryError := stmt.Query(newValues...)
	if queryError != nil {
		return &DBResult{
			Err: queryError,
		}
	}

	return &DBResult{
		Rows: rows,
		Err:  nil,
	}
}

func UpdateTable(table string, whereClause string, columns []string, values []string) *DBResult {
	if len(columns) != len(values) {
		return &DBResult{
			Err: errors.New("columns dont match the values"),
		}
	}

	queyStr := fmt.Sprintf("UPDATE %s SET %s %s", table, updateStmtStr(columns), whereClause)
	var newValues = arrayToInterfaceArr(values)
	stmt, prepareError := dbInstace.Prepare(queyStr)

	if prepareError != nil {
		return &DBResult{
			Err: prepareError,
		}
	} else {
		defer stmt.Close()
	}

	update, updateError := stmt.Exec(newValues...)

	if updateError != nil {
		return &DBResult{
			Err: updateError,
		}
	}

	affected, affectedError := update.RowsAffected()

	if affectedError != nil {
		return &DBResult{
			Err: affectedError,
		}
	}
	return &DBResult{
		Affected: int(affected),
		Err:      nil,
	}
}

func DeleteFromTable(table string, whereClause string, values []string) *DBResult {
	queyStr := fmt.Sprintf("DELETE FROM %s %s", table, whereClause)
	newValues := arrayToInterfaceArr[string](values)
	update, updateError := dbInstace.ExecContext(context.Background(), queyStr, newValues...)

	if updateError != nil {
		return &DBResult{
			Err: updateError,
		}
	}

	affected, affectedError := update.RowsAffected()

	if affectedError != nil {
		return &DBResult{
			Err: affectedError,
		}
	}
	return &DBResult{
		Affected: int(affected),
		Err:      nil,
	}
}

type DBResult struct {
	Rows         *sql.Rows
	Err          error
	Affected     int
	LastInsertId int
}
