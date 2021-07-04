package main

import (
	"golang-student/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(":8080")
}
