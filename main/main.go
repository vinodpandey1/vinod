// main.go

//go:generate swagger generate spec

package main

import "github.com/vinodpandey1/main/app"

func main() {
	a := app.App{}
	a.Initialize("root", "", "bookdb")
	a.Run(":8080")
}
