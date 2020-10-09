package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.56.38:3306)/go_product?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		ID       int
		userName string
	)

	rows, err := db.Query("SELECT `ID`, `userName` FROM `user` WHERE `ID`= ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 必须把 rows 里面的内容计完，或者显示的调用Close（） 方法，
	// 否则 defer 的rows.Close(0执行之前连接永远不会释放

	for rows.Next() {
		err := rows.Scan(&ID, &userName)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(ID, userName)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
