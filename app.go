// cloudsql.go - Creates Google Cloud SQL table
package cloudSQLDatabase

import (
    "google.golang.org/appengine"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

func main() {
    const dbUserName = "root"
    const dbPassword = "dog"
    const dbName = "database2"
    //const dbIP = "2001:4860:4864:1:de34:1928:6ae4:7058"
    const dbIP = "tcp(130.211.122.232:3306)"
    const dbOpenString = dbUserName + ":" + dbPassword + "@" + dbIP + "/" + dbName
    db, err := sql.Open("mysql", dbOpenString);
    if err != nil {
        log.Println("sql.Open(" +
            dbOpenString +
            "\"mysql, \"")
    }
    defer db.Close()
    log.Println("Pinging database. This may take a moment.")
    err = db.Ping()
    if err != nil {
        log.Println("db.Ping() call failed:");
        log.Println(err)
    }
    _, err = db.Exec(
        `CREATE TABLE IF NOT EXISTS exercisecloudsql101
        (id INT NOT NULL AUTO_INCREMENT,
        name VARCHAR(100) NOT NULL,
        description TEXT, PRIMARY KEY (id))`)

    if err != nil {
        log.Println("CREATE TABLE failed:")
        log.Println(err) 
    }

    appengine.Main()

}