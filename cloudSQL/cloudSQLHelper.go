package cloudSQL

import (
        "log"
        "net/http"
        _ "github.com/go-sql-driver/mysql"
        "github.com/samuelechu/jsonHelper"
)

var insertUserStmt *mysql.Stmt

func initPrepareStatements() {
    insertUserStmt, err := db.Prepare(`INSERT INTO users (uid, Name, refreshToken) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE
                                refreshToken = ?`)
    checkErr(err)
}

func InsertUser(user_id string, name string, refresh_token string) {
	
    if refresh_token != "" {
        _, err := insertUserStmt.Exec(user_id, name, refresh_token, refresh_token)
        checkErr(err)
        log.Printf("inserted refresh token for %v!", name)
    }
}

func signInHandler(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
                http.NotFound(w, r)
                return
    }

    var u, user jsonHelper.User
    if u, ok := jsonHelper.UnmarshalJSON(w, r, r.Body, u).(jsonHelper.User); ok {
        user = u
        log.Printf("UnmarshalJSON returned %v %v", user.Uid, user.Name)

    }

    stmt, err := db.Prepare("INSERT IGNORE INTO users SET uid=?, Name=?")
    checkErr(err)

    res, err := stmt.Exec(user.Uid, user.Name)
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