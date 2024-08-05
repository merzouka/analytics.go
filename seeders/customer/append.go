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
    AddTransaction(*Transaction)
    Finalize()
    Close()
}

type DB struct {
    conn *gorm.DB
}

type File struct {
    ptr *os.File
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

func (db DB) AddTransaction(transaction *Transaction) {
    conn := db.conn
    if conn == nil {
        log.Fatal(WRITER_ERROR)
    }
    conn.Save(&transaction)
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
    ptr := f.ptr
    if ptr == nil {
        log.Fatal(WRITER_ERROR)
    }

    _, err := ptr.Write([]byte("CREATE TABLE IF NOT EXISTS customers (id serial primary key, name varchar(255), age int, country varchar(255), language varchar(50));\nCREATE TABLE IF NOT EXISTS transactions (customer_id bigint references customers(id), transaction_id bigint);"))
    if err != nil {
        log.Println(fmt.Sprintf("failed to customers write table definition, error: %s", err))
        return
    }
}

func (f *File) AddCustomer(customer *Customer) {
    ptr := f.ptr
    if ptr == nil {
        log.Fatal(WRITER_ERROR)
    }

    row := fmt.Sprintf("(%d, '%s', %d)", uint(customer.ID), strings.Replace(customer.Name, "'", "''", -1), customer.Age)
    if f.written {
        row = ",\n" + row
    } else {
        row = "INSERT INTO customers (id, name, age) VALUES " + row
        f.written = true
    }
    ptr.Write([]byte(row))
}

func (f *File) AddTransaction(transaction *Transaction) {
    ptr := f.ptr
    if ptr == nil {
        log.Fatal(WRITER_ERROR)
    }

    row := fmt.Sprintf("(%d, %d)", uint(transaction.CustomerID), uint(transaction.TransactionID))
    if f.written {
        row = ",\n" + row
    } else {
        row = "INSERT INTO transactions (customer_id, transaction_id) VALUES " + row
        f.written = true
    }
    ptr.Write([]byte(row))
}

func (f *File) Finalize() {
    ptr := f.ptr
    if ptr == nil {
        log.Fatal(WRITER_ERROR)
    }
    ptr.Write([]byte(";\n"))
    (*f).written = false
}

func (f File) Close() {
    f.ptr.Close()
}

func getAppender() Appender {
    output := os.Getenv("OUTPUT_MEDIUM")
    if output == "" {
        output = "FILE"
    }

    dest := os.Getenv("OUTPUT_DESTINATION")
    if dest == "" {
        dest = "./sql/init.sql"
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

    ptr, err := os.Create(dest)
    if err != nil {
        log.Fatal(WRITER_ERROR)
    }
    file = &File{
        ptr: ptr,
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
