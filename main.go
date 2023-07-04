package main

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"os"
)

func main() {
	client := godo.NewFromToken(os.Getenv("DIGITAL_OCEAN_PAT"))
	ctx := context.TODO()
	fmt.Println(client.BillingHistory.List(ctx, nil))
}
