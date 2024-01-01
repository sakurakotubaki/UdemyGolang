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
curl -X POST -H "Content-Type: application/json" -d '{"name":"kurokawa", "age": 19}' http://localhost:8080/users
```

データを取得する
```bash
curl -XGET http://localhost:8080/users
```

指定したidのデータを取得する
```bash
curl -XGET http://localhost:8080/users/3
```

id1のデータを更新する
```bash
curl -XPUT -H "Content-Type: application/json" -d '{"name":"JboyHashimoto", "age": 25}' http://localhost:8080/users/1
```

追加、表紙、更新、削除、特定のidのデータを検索するコード

```go
package main

import (
    "database/sql"
    "log"
    "net/http"
    "strconv"

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

func validateUser(name string, age int) error {
    if name == "" {
        return echo.NewHTTPError(http.StatusBadRequest, "name is empty")
    }

    if len(name) > 100 { // let(name) を len(name) に修正
        return echo.NewHTTPError(http.StatusBadRequest, "name is too long")
    }

    if age < 0 || age >= 200 {
        return echo.NewHTTPError(http.StatusBadRequest, "age must be between 0 and 200")
    }
    return nil
}

func main()  {
    db := initDB("./example.db")
    e :=  echo.New()
    e.Use(middleware.Logger())

    // データの削除
    e.DELETE("/users/:id", func(c echo.Context) error {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }

        result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }

        rowsAffected, err := result.RowsAffected()
        if rowsAffected == 0 {
            return echo.NewHTTPError(http.StatusNotFound, "not Found")
        }

        return c.NoContent(http.StatusNoContent)
    })
    // データを追加
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

    e.PUT("/users/:id", func(c echo.Context) error {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }

        u := new(User)
        if err := c.Bind(u); err != nil {
            return err
        }

        if err := validateUser(u.Name, u.Age); err != nil { // name, age を u.Name, u.Age に修正
            return err
        }

        result, err := db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", u.Name, u.Age, id) // name, age を u.Name, u.Age に修正
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }

        rows, err := result.RowsAffected()
        if rows == 0 {
            return echo.NewHTTPError(http.StatusNotFound, "not Found")
        }

        return c.JSON(http.StatusOK, u)
    })

    // データを取得
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

    // 特定のユーザーのデータを取得
    e.GET("/users/:id", func(c echo.Context) error {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
        // SQL文で特定のユーザーのデータを取得
        row := db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", id)
        // rowでScanをUser構造体に対して実行して、データを取得
        var user User
        if err := row.Scan(&user.ID, &user.Name, &user.Age); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }

        return c.JSON(http.StatusOK, user)

    })

    e.Start(":8080")
}
```