package cloudSQL

import (
        "log"
        "errors"
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
)

var insertUserStmt *sql.Stmt
var insertThreadStmt *sql.Stmt
var getRefTokenStmt *sql.Stmt

func initPrepareStatements() {
    var err error
    
    insertUserStmt, err = db.Prepare(`INSERT INTO users (uid, Name, refreshToken) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE
                                refreshToken = ?`)
    checkErr(err)

    insertThreadStmt, err = db.Prepare(`INSERT IGNORE INTO threads (uid, thread_id) VALUES(?, ?)` )
    checkErr(err)

    getRefTokenStmt, err = db.Prepare(`SELECT refreshToken FROM users WHERE uid = ?`)
    checkErr(err)

}

func InsertUser(user_id, name, refresh_token string) {
	
    if refresh_token != "" {
        _, err := insertUserStmt.Exec(user_id, name, refresh_token, refresh_token)
        checkErr(err)
        log.Printf("inserted refresh token for %v!", name)
    } else {
        stmt, err := db.Prepare("INSERT IGNORE INTO users SET uid=?, Name=?")
        checkErr(err)

        _, err := stmt.Exec(user_id, name)
        checkErr(err)
    }
}

func InsertThread(uid, thread_id string) {
    _, err := insertThreadStmt.Exec(uid, thread_id)
        checkErr(err)

    log.Printf("inserted thread %v!", thread_id)

} 

func GetRefreshToken(uid string) (string, error){
    
    result, err := getRefTokenStmt.Query(uid)
    checkErr(err)
    result.Next()

    var refToken string
    err = result.Scan(&refToken)
    checkErr(err)

    if refToken == "" {
        return refToken, errors.New("Error: refreshToken not found")
    }

    return refToken, nil
}


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}