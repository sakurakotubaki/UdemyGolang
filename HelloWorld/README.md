# Go言語でHello World

Hello Worldするだけなら、`go mod init`は不要

```go
package main

import "fmt"// goにもともとあるパッケージ

func main()  {
	fmt.Println("Hello World")
}
```