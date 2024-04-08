package main

import (
	"adtelligent-internship/db"
	"github.com/valyala/fasthttp"
	"testing"
)

func BenchmarkCampaignHandler(b *testing.B) {
	database, err := db.ConnectToDB()
	if err != nil {
		b.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := &fasthttp.RequestCtx{}

	ctx.QueryArgs().Set("source_id", "1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CampaignHandler(ctx, database)
	}
}
