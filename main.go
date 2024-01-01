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