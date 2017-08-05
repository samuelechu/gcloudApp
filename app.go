// cloudsql.go - Creates Google Cloud SQL table
package cloudSQLDatabase

import (
    "bytes"
    "fmt"
    "net/http"

    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func init() {
        http.HandleFunc("/initDatabase", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/initDatabase" {
                http.NotFound(w, r)
                return
        }

        connectionName := "gotesting-175718:database"        
        user := "root"
        password := "dog" // NOTE: password may be empty

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

//}