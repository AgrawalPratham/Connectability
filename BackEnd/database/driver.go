package database

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifeTime = 5 * time.Minute

//NewDatabase creates a new database connection for the application
func NewDatabase(dsn string) (*sql.DB, error){
	db, err := sql.Open("pgx", dsn)
	if err != nil{
		fmt.Print("Error creating connection with the database")
		return nil, err
	}
	if err := db.Ping(); err!=nil{
		fmt.Println("Error during ping to the database")
		return nil, err
	}
	return db, nil
}

//ConnectSQL creates database pool for postgres
func ConnectSQL( dsn string) ( *DB, error){
	d, err := NewDatabase(dsn)
	if(err != nil){
		fmt.Println("Error while creating the pool of database : ", err)
	}

	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetMaxIdleConns(maxIdleDBConn)
	d.SetConnMaxLifetime(maxDBLifeTime)

	dbConn.SQL = d
	err = d.Ping()
	if( err != nil){
		fmt.Println("Error while ping to database after creating pool of connection : ", err)
		return nil, err
	}
	return dbConn, nil
}