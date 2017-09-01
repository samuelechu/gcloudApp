package cloudSQL

import (
        "log"
        "net/http"
        _ "github.com/go-sql-driver/mysql"
        "github.com/samuelechu/jsonHelper"
)

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