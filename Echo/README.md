# Echoを使ってみる

## インストール
```bash
go get -u github.com/labstack/echo/v4
```

go.sumというファイルが生成される。この中には依存ライブラリーが記述されている。

`go mod tidy`を実行すると、go.modに記述されている依存ライブラリーがgo.sumに記述される。

Echoのミドルウェアを使うとルーティングのログを出力や、リクエストのヘッダーをチェックすることができる。

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main()  {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.GET("/users/", func(c echo.Context) error {
		return c.String(200, "Hello, Users!")
	})

	e.Use(middleware.Logger())

	e.Start(":8080")
}
```