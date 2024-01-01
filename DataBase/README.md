# SQLiteを導入する

## 1. SQLiteをインストールする
```bash
go get -u github.com/mattn/go-sqlite3
```

## 2. データベースを作成する
`go run main.go`を実行した後に、`example.db`が作成されるが、時間がかかる。

```go
package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main()  {
	// SQLiteに接続
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()// deferは、関数が終了する際に実行される処理を指定するためのものです。
  // usersテーブルを作成
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL
	);
  `

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Table Created")
}
```

## 3. データの追加とfetch
`main.go`を以下のように変更する。

```go
package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    _ "github.com/mattn/go-sqlite3"
)

type User struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Age int `json:"age"`
}

func initDB(filepath string) *sql.DB {
    db, err := sql.Open("sqlite3", filepath)
    if err != nil {
        log.Fatal(err)
    }

    // Create table if not exists
    tableCreationQuery := `CREATE TABLE IF NOT EXISTS users (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "name" TEXT,
        "age" INTEGER
    );`

    _, err = db.Exec(tableCreationQuery)
    if err != nil {
        log.Fatal(err)
    }

    return db
}

func main()  {
    db := initDB("./example.db")
    e :=  echo.New()
    e.Use(middleware.Logger())

    e.POST("/users", func(c echo.Context) error {
			u := new(User)
			if err := c.Bind(u); err != nil {
					return err
			}

			result, err := db.Exec("INSERT INTO users(name, age) VALUES(?, ?)", u.Name, u.Age)
			if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			id, err := result.LastInsertId()
			if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get the last inserted id")
			}
			return c.JSON(http.StatusOK, &User{ID: int(id), Name: u.Name, Age: u.Age})
	})

    e.GET("/users", func(c echo.Context) error {
        rows, err := db.Query("SELECT * FROM users")
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
        defer rows.Close()

        var users []User
        for rows.Next() {
            var user User
            if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
                return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
            }
            users = append(users, user)
        }

        return c.JSON(http.StatusOK, users)
    })

    e.Start(":8080")
}
```

POSTのテストをする
```bash
curl -XPOST localhost:8080/users
```

## 3. データを挿入する
nameとageを挿入する
```bash
curl -X POST -H "Content-Type: application/json" -d '{"name":"test", "age": 20}' http://localhost:8080/users
```

データを取得する
```bash
curl -XGET http://localhost:8080/users
```