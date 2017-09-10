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

        _, err = stmt.Exec(user_id, name)
        checkErr(err)
    }
}

func InsertThread(uid, thread_id string) {
    _, err := insertThreadStmt.Exec(uid, thread_id)
        checkErr(err)

    log.Printf("inserted thread %v!", thread_id)

}

func GetThreadsForUser(uid string) []string {

    getThreadsStmt, err := db.Prepare(`SELECT thread_id FROM threads WHERE uid=? AND done='F'`)
    checkErr(err)

    rows, err := getThreadsStmt.Query(uid)
    checkErr(err)

    var threads []string
    defer rows.Close()
    for rows.Next() {
        var thread_id string
        err = rows.Scan(&thread_id)
        threads = append(threads, thread_id)
        checkErr(err)
    }
    // get any error encountered during iteration
    err = rows.Err()
    checkErr(err)

    return threads
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