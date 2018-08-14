package main

import (
    "database/sql"
    "fmt"
  	"errors"
  	"strings"

    _ "github.com/lib/pq"
)

const (
    c_host     = "localhost"
    c_port     = 5432
    c_user     = "gouser"
    c_password = "justgo"
    c_dbname   = "gogintest"
)
var db *sql.DB
var err error
func initSql() {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
      "password=%s dbname=%s sslmode=disable",
      c_host, c_port, c_user, c_password, c_dbname)
  db, err = sql.Open("postgres", psqlInfo)
  if err != nil {
      panic(err)
  }

  err = db.Ping()
  if err != nil {
      panic(err)
  }

  fmt.Println("Successfully connected!")
}
func closeSql() {
  defer db.Close()
}

func verifyingUser(i_username, i_password string) bool {
    initSql()
		var (
			name string
			password string
		)
		rows, err := db.Query("select * from users where name = $1", i_username)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
    closeSql()
		for rows.Next() {
			err := rows.Scan(&name, &password)
			if err != nil {
				panic(err)
			}
			fmt.Println(name, password)
      if name != "" || password != "" {
        if password == i_password {
          return true
        }
      }
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
    return false
}

func addNewUser(username, password string) (*user, error) {
  initSql()
  if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}
	u := user{Username: username, Password: password}

  _, err := db.Exec("INSERT INTO users (name, password) VALUES ($1, $2)", u.Username, u.Password)
  if err != nil {
    panic(err)
  }
  closeSql()
	return &u, nil
}
