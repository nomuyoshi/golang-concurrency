package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest(123, "NomuYoshi")
}

func ProcessRequest(userID int, authToken string) {
	ctx := context.WithValue(context.Background(), "userID", userID)
	ctx = context.WithValue(ctx, "authToken", authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	// ctx.Valueでコンテキストから値を取得
	fmt.Printf("handling response for %v (%v)\n", ctx.Value("userID"), ctx.Value("authToken"))
}
