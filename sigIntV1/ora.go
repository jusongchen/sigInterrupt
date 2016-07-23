package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	// _ "github.com/mattn/go-oci8"
	"github.com/pkg/errors"
	_ "gopkg.in/rana/ora.v3"
)

var (
	dbUser   = flag.String("dbUser", "scott", "The user name used to connect to the Oracle server")
	dbPasswd = flag.String("dbPassword", "", "password of the DB user")
	dbServer = flag.String("dbServer", "//localhost/orcl", "The database server")
)

type record struct {
	owner       string
	objectName  string
	lastDDLtime time.Time
}

func testOracleDB() {

	flag.Parse()

	if *dbPasswd == "" {
		flag.Usage()
		return
	}

	var db *sql.DB
	db, err := OpenDB(*dbUser, *dbPasswd, *dbServer)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = printObjects(db)
	if err != nil {
		log.Print(err)
	}
	return
}

func OpenDB(oraUser, oraPasswd, oraConn string) (*sql.DB, error) {

	// connStr := fmt.Sprintf("%s:%s@%s", oraUser, oraPasswd, oraConn)
	connStr := fmt.Sprintf("%s/%s@%s", oraUser, oraPasswd, oraConn)
	// db, err := sql.Open("oci8", connStr)
	db, err := sql.Open("ora", connStr)
	// defer db.Close()

	if err != nil {
		return nil, errors.Wrapf(err, "connect to oracle %s as user %s failed", oraConn, oraUser)
	}

	//make a SQL query call to make sure the DB server works
	var n int
	err = db.QueryRow("select 1 from dual").Scan(&n)
	if err != nil {
		return nil, errors.Wrapf(err, "connect to oracle %s as user %s failed", oraConn, oraUser)
	}

	if n != 1 {
		panic(fmt.Sprintf("OpenDB:`select 1 from dual` fail. Expecting 1 get %d", n))
	}
	return db, nil

}

func printObjects(db *sql.DB) error {

	rows, err := db.Query(
		//get last modified db objects
		`with rs as (
			select owner, object_name,last_ddl_time 
			from all_objects 
			order by last_ddl_time desc )
			select * from rs
			where rownum<=10`)

	if err != nil {
		return errors.Wrap(err, "DB query fail:")
	}
	defer rows.Close()

	r := record{}
	for rows.Next() {

		if err := rows.Scan(&r.owner, &r.objectName, &r.lastDDLtime); err != nil {
			return errors.Wrap(err, "DB row scan fail")
		}

		fmt.Printf("owner:%s ,obj_name:%s, lastDDLtime:%s\n", r.owner, r.objectName, r.lastDDLtime)

	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "DB query fail")
	}
	return nil

}
