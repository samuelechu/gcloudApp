// cloudsql.go - Creates Google Cloud SQL table
package cloudSQL

// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample cloudsql demonstrates connection to a Cloud SQL instance from App Engine standard.
import (
        "google.golang.org/appengine"
        "bytes"
        "database/sql"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
        initDB()
        http.HandleFunc("/signIn", signInHandler)
        http.HandleFunc("/showDatabases", showDatabases)
}

// func GetDB() *sql.DB {
//     return db
// }

func InsertUser(user_id string, name string, refresh_token string) {
    
}

func initDB(){
    var err error

    user := "root"
    password := "dog"
    instance := "gotesting-175718:us-central1:database"
    dbName := "mailMigrationDatabase"
    
    // dbOpenString := "root:dog@cloudsql(gotesting-175718:us-central1:database)/samsDatabase"
    dbOpenString := fmt.Sprintf("%s:%s@cloudsql(%s)/%s", user, password, instance, dbName)

    if appengine.IsDevAppServer() {
            dbOpenString = fmt.Sprintf("%s:%s@tcp([localhost]:3306)/%s", user, password, dbName)
    }

    db, err = sql.Open("mysql", dbOpenString)

    if err != nil {
        log.Print("Could not open db: %v", err)
        return    
    }

    _, err = db.Exec(
                `CREATE TABLE IF NOT EXISTS users
                (uid VARCHAR(64) UNIQUE,
                firstname VARCHAR(64) NULL DEFAULT NULL,
                PRIMARY KEY (uid))`)


    if err != nil {
        log.Printf("CREATE TABLE failed: %v", err)
    }
}

func showDatabases(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")

    rows, err := db.Query("SHOW DATABASES")
    if err != nil {
            http.Error(w, fmt.Sprintf("Could not query db: %v.", err), 500)
            return
    }
    defer rows.Close()

    buf := bytes.NewBufferString("Databases:\n")

    for rows.Next() {
            var dbName string
            if err := rows.Scan(&dbName); err != nil {
                    http.Error(w, fmt.Sprintf("Could not scan result: %v", err), 500)
                    return
            }
            fmt.Fprintf(buf, "- %s\n", dbName)
    }

    w.Write(buf.Bytes())
}

type User struct{
    Uid     string
    Name    string
}

func signInHandler(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
                http.NotFound(w, r)
                return
    }

    var u User
    if r.Body == nil {
        http.Error(w, "Please send a request body", 400)
        return
    }
    err := json.NewDecoder(r.Body).Decode(&u)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    stmt, err := db.Prepare("INSERT IGNORE INTO users SET uid=?, Name=?")
    checkErr(err)

    res, err := stmt.Exec(u.Uid, u.Name)
    checkErr(err)

    id, err := res.RowsAffected()
    // checkErr(err)

    log.Println(id)
}

func checkErr(err error) {
        if err != nil {
            panic(err)
        }
}