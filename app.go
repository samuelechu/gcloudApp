// cloudsql.go - Creates Google Cloud SQL table
package cloudSQLTest

// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample cloudsql demonstrates connection to a Cloud SQL instance from App Engine standard.
import (
        "bytes"
        "database/sql"
        "fmt"
        "net/http"
        "os"
        _ "github.com/go-sql-driver/mysql"
)


type MySQLConfig struct {
        // Optional.
        Username, Password string

        // Host of the MySQL instance.
        //
        // If set, UnixSocket should be unset.
        Host string

        // Port of the MySQL instance.
        //
        // If set, UnixSocket should be unset.
        Port int

        // UnixSocket is the filepath to a unix socket.
        //
        // If set, Host and Port should be unset.
        UnixSocket string
}

// dataStoreName returns a connection string suitable for sql.Open.
func (c MySQLConfig) dataStoreName(databaseName string) string {
        var cred string
        // [username[:password]@]
        if c.Username != "" {
                cred = c.Username
                if c.Password != "" {
                        cred = cred + ":" + c.Password
                }
                cred = cred + "@"
        }

        if c.UnixSocket != "" {
                return fmt.Sprintf("%sunix(%s)/%s", cred, c.UnixSocket, databaseName)
        }
        return fmt.Sprintf("%stcp([%s]:%d)/%s", cred, c.Host, c.Port, databaseName)
}

func getDataStoreName(username, password, instance, databaseName string) string {
        if os.Getenv("GAE_INSTANCE") != "" {
                // Running in production.
                return MySQLConfig{
                        Username:   username,
                        Password:   password,
                        UnixSocket: "/cloudsql/" + instance,
                }.dataStoreName(databaseName)
        }


        // Running locally.
        return MySQLConfig{
                Username: username,
                Password: password,
                Host:     "localhost",
                Port:     3306,
        }.dataStoreName(databaseName)
}



func init() {
        http.HandleFunc("/initDB", handler)
}



func handler(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/initDB" {
                http.NotFound(w, r)
                return
        }

        w.Header().Set("Content-Type", "text/plain")

        dbUserName := "root"
        dbPassword := "dog"
        dbInstance := "gotesting-175718:us-central1:database"
        dbName := "samsDatabase"
        dbOpenString := getDataStoreName(dbUserName, dbPassword, dbInstance, dbName)

        db, err := sql.Open("mysql", dbOpenString)

        if err != nil {
                http.Error(w, fmt.Sprintf("Could not open db: %v", err), 500)
                return
        }
        defer db.Close()

        rows, err := db.Query("SHOW DATABASES")
        if err != nil {
                http.Error(w, fmt.Sprintf("Could not query db: %v", err), 500)
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