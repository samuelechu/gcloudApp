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
        "fmt"
        "net/http"
        _ "github.com/go-sql-driver/mysql"
)

func init() {
        http.HandleFunc("/initDB", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {

        w.Header().Set("Content-Type", "text/plain")

        user := "root"
        password := "dog"
        instance := "gotesting-175718:us-central1:database"
        dbName := "samsDatabase"
        
        // dbOpenString := "root:dog@cloudsql(gotesting-175718:us-central1:database)/samsDatabase"
        dbOpenString := fmt.Sprintf("%s:%s@cloudsql(%s)/%s", user, password, instance, dbName)

        if appengine.IsDevAppServer() {
                dbOpenString = fmt.Sprintf("%s:%s@tcp([localhost]:3306)/%s", user, password, dbName)
        }
    
        db, err := sql.Open("mysql", dbOpenString)

        if err != nil {
                http.Error(w, fmt.Sprintf("Could not open db: %v", err), 500)
                return    
        }
        defer db.Close()

        rows, err := db.Query("SHOW DATABASES")
        if err != nil {
                http.Error(w, fmt.Sprintf("Could not query db: %v. DBString: %s", err, dbOpenString), 500)
                return
        }
        defer rows.Close()

        buf := bytes.NewBufferString("dbOpenString: " + dbOpenString + "\nDatabases:\n")
        for rows.Next() {
                var dbName string
                if err := rows.Scan(&dbName); err != nil {
                        http.Error(w, fmt.Sprintf("Could not scan result: %v", err), 500)
                        return
                }
                fmt.Fprintf(buf, "- %s\n", dbName)
        }







        // insert
        stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
        checkErr(err)

        res, err := stmt.Exec("Sam", "comp sci", "2012-12-09")
        checkErr(err)

        id, err := res.LastInsertId()
        checkErr(err)

        fmt.Println(id)
        // update
        stmt, err = db.Prepare("update userinfo set username=? where uid=?")
        checkErr(err)

        res, err = stmt.Exec("samupdate", id)
        checkErr(err)

        affect, err := res.RowsAffected()
        checkErr(err)

        fmt.Println(affect)

        // query
        rows, err := db.Query("SELECT * FROM userinfo")
        checkErr(err)

        for rows.Next() {
            var uid int
            var username string
            var department string
            var created string
            err = rows.Scan(&uid, &username, &department, &created)
            checkErr(err)
            fmt.Println(uid)
            fmt.Println(username)
            fmt.Println(department)
            fmt.Println(created)
        }

        // delete
        stmt, err = db.Prepare("delete from userinfo where uid=?")
        checkErr(err)

        res, err = stmt.Exec(id)
        checkErr(err)

        affect, err = res.RowsAffected()
        checkErr(err)

        fmt.Println(affect)

        db.Close()

        w.Write(buf.Bytes())
}

func checkErr(err error) {
        if err != nil {
            panic(err)
        }
}