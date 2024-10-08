package main

import (
    "fmt"
    "log"
    "os"
    "strings"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Appender interface {
    Define()
    AddCustomer(*Customer)
    Finalize()
    Close()
}

type DB struct {
    conn *gorm.DB
}

type File struct {
    def *os.File
    data *os.File
    written bool
}

const (
    WRITER_ERROR = "connection to writer failed"
) 

func (db DB) Define() {
    conn := db.conn
    if conn == nil {
        log.Fatal(WRITER_ERROR)
    }

    conn.AutoMigrate(&Customer{}, &Transaction{})
}

func (db DB) AddCustomer(customer *Customer) {
    conn := db.conn
    if conn == nil {
        log.Fatal(WRITER_ERROR)
    }
    conn.Save(&customer)
}

func (DB) Finalize() {}

func (db DB) Close() {
    conn, err := db.conn.DB()
    if err != nil {
        return
    }
    conn.Close()
}

func (f File) Define() {
    ptr := f.def
    if ptr == nil {
        log.Fatal(WRITER_ERROR)
    }

    _, err := ptr.Write([]byte("CREATE TABLE IF NOT EXISTS customers (id serial primary key, name varchar(255), age int, country varchar(255), language varchar(50));"))
    if err != nil {
        log.Println(fmt.Sprintf("failed to customers write table definition, error: %s", err))
        return
    }
    log.Println("created table definition successfully")
}

func (f *File) AddCustomer(customer *Customer) {
    ptr := f.data
    if ptr == nil {
        log.Fatal(WRITER_ERROR)
    }

    row := fmt.Sprintf("(%d, '%s', %d)", uint(customer.ID), strings.Replace(customer.Name, "'", "''", -1), customer.Age)
    if f.written {
        row = ",\n" + row
    } else {
        row = "TRUNCATE TABLE customers RESTART IDENTITY CASCADE;\nINSERT INTO customers (id, name, age) VALUES " + row
        f.written = true
    }
    ptr.Write([]byte(row))
}

func (f *File) Finalize() {
    ptr := f.data
    if ptr == nil {
        log.Fatal(WRITER_ERROR)
    }
    ptr.Write([]byte(";\n"))
    (*f).written = false
}

func (f File) Close() {
    f.data.Close()
}

func getAppender() Appender {
    output := os.Getenv("OUTPUT_MEDIUM")
    if output == "" {
        output = "FILE"
    }

    dest := os.Getenv("OUTPUT_DESTINATION")
    if dest == "" {
        dest = "./sql"
    }

    if output == "DATABASE" {
        return getDB(dest)
    }

    return getFile(dest)
}

var file *File
func getFile(dest string) *File {
    if file != nil {
        return file
    }

    data, err := os.Create(fmt.Sprintf("%s/02-populate-tables.sql", dest))
    if err != nil {
        log.Fatal(err)
    }

    def, err := os.Create(fmt.Sprintf("%s/01-create-tables.sql", dest))
    if err != nil {
        log.Fatal(err)
    }

    file = &File{
        def: def,
        data: data,
        written: false,
    }

    return file
}

var db *DB
func getDB(dest string) *DB {
    if db != nil {
        return db
    }

    conn, err := gorm.Open(postgres.Open(dest), &gorm.Config{})
    if err != nil {
        log.Fatal(WRITER_ERROR)
    }
    db = &DB{
        conn: conn,
    }

    return db
}
