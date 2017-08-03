// cloudsql1.go - Creates Google Cloud SQL table
package main

import (
    "google.golang.org/appengine"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
)

func main() {
    fmt.Println("In Main!")
    const dbUserName = "root"
    const dbPassword = "dog"
    const dbName = "cloudsqltest-175704:us-central1:database2"
    //const dbIP = "2001:4860:4864:1:de34:1928:6ae4:7058"
    const dbIP = "tcp(130.211.122.232:3306)"
    const dbOpenString = dbUserName + ":" + dbPassword + "@" + dbIP + "/" + dbName
    fmt.Println("attempt open database!")
    db, err := sql.Open("mysql", dbOpenString);
    if err != nil {
        fmt.Println("sql.Open(" +
            dbOpenString +
            "\"mysql, \"")
    }
    defer db.Close()
    fmt.Println("Pinging database. This may take a moment.")
    err = db.Ping()
    if err != nil {
        fmt.Println("db.Ping() call failed:");
        fmt.Println(err)
    }
    _, err = db.Exec(
        `CREATE TABLE IF NOT EXISTS exercisecloudsql101
        (id INT NOT NULL AUTO_INCREMENT,
        name VARCHAR(100) NOT NULL,
        description TEXT, PRIMARY KEY (id))`)

    if err != nil {
        fmt.Println("CREATE TABLE failed:")
        fmt.Println(err) 
    }

    appengine.Main()

}