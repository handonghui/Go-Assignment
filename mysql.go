package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	xerrors "github.com/pkg/errors"
	"os"
)

var errSQL = errors.New("SQL")

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/book_manager?charset=utf8")
	if err != nil {
		fmt.Print(`err`)
		fmt.Println(err)
		return
	}
	fmt.Println(db)
	err = showTable(db)
	if err != nil {
		fmt.Printf("original error: %T %v\n", xerrors.Cause(err), xerrors.Cause(err))
		fmt.Printf("stack trace:\n %+v\n", err)
		os.Exit(1)
	}

}

func showTable(db *sql.DB) error {
	//表结构
	type info struct {
		Id         int    `db:"id"`
		Name       string `db:"name"`
		Author     string `db:"author"`
		Status     string `db:"status"`
		CreateTime string `db:"create_time"`
	}
	sqlString := "SELECT * FROM book"
	rows, err := db.Query(sqlString)
	if err != nil {
		fmt.Println("sql:" + sqlString)
		return xerrors.Wrap(errSQL, "sql:"+sqlString)
	}

	for rows.Next() {
		var s info
		err = rows.Scan(&s.Id, &s.Name, &s.Author, &s.Status, &s.CreateTime)
		if err != nil {
			break
		}
		fmt.Println(s)
	}
	//用完关闭
	err = rows.Close()
	if err != nil {
		return err
	}
	return nil
}
