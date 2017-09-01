package cloudSQL

import (
        "log"
        "net/http"
        _ "github.com/go-sql-driver/mysql"
        "github.com/samuelechu/jsonHelper"
)

func InsertUser(user_id string, name string, refresh_token string) {
	stmt, err := db.Prepare(`INSERT INTO users (uid, Name, refreshToken) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE
								refreshToken = ?`)
	checkErr(err)

	_, err := stmt.Exec(user_id, name, refresh_token, refresh_token)
    checkErr(err)


	//INSERT INTO table (id, name, age) VALUES(1, "A", 19) ON DUPLICATE KEY UPDATE    
	//insert into users (uid, Name, refreshToken) values("0", "testUser", "333ff") on duplicate key update  Name = "ffd", refreshToken = "Woww Ia m refresh";
    log.Printf("inserted refresh token for %v!", name)
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