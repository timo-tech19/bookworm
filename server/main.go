package main

import (
	"github.com/timotech-19/bookworm/routes"
)

func main() {
	h := routes.NewHandler()

	h.Serve()
}
