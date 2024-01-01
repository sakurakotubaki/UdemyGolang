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