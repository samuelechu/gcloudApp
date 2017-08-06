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
        "log"
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

        const dbUserName = "root"
        const dbPassword = "dog"
        const dbInstance = "gotesting-175718:us-central1:database"
        const dbName = "samsDatabase"
        const dbOpenString = getDataStoreName(dbUserName, dbPassword, dbInstance, dbName)
        db, err := sql.Open("mysql", dbOpenString);
        if err != nil {
        log.Println("sql.Open(" +
            dbOpenString +
            "\"mysql, \"")
        }

        connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
        user := mustGetenv("CLOUDSQL_USER")
        password := os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty

        w.Header().Set("Content-Type", "text/plain")

        db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/", user, password, connectionName))
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



// func mains() {
//     const dbUserName = "root"
//     const dbPassword = "dog"
//     const dbName = "database2"
//     //const dbIP = "2001:4860:4864:1:de34:1928:6ae4:7058"
//     const dbIP = "tcp(130.211.122.232:3306)"
//     const dbOpenString = dbUserName + ":" + dbPassword + "@" + dbIP + "/" + dbName
//     db, err := sql.Open("mysql", dbOpenString);
//     if err != nil {
//         log.Println("sql.Open(" +
//             dbOpenString +
//             "\"mysql, \"")
//     }
//     defer db.Close()
//     log.Println("Pinging database. This may take a moment.")
//     err = db.Ping()
//     if err != nil {
//         log.Println("db.Ping() call failed:");
//         log.Println(err)
//     }
//     _, err = db.Exec(
//         `CREATE TABLE IF NOT EXISTS exercisecloudsql101
//         (id INT NOT NULL AUTO_INCREMENT,
//         name VARCHAR(100) NOT NULL,
//         description TEXT, PRIMARY KEY (id))`)

//     if err != nil {
//         log.Println("CREATE TABLE failed:")
//         log.Println(err) 
//     }

//     appengine.Main()

// }