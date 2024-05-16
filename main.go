package main

import (
	"context"
)

func Handler(ctx context.Context) (string, error) {
	return "Hello from Lambda!", nil
}

func main() {}
