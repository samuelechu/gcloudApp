package cloudSQL

import (
        "google.golang.org/appengine"
        "bytes"
        "database/sql"
        "fmt"
        "log"
        "net/http"
        _ "github.com/go-sql-driver/mysql"
        "github.com/samuelechu/jsonHelper"
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
    log.Print("In InsertUser!")
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


