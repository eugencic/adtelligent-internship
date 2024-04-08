package api

import (
	"adtelligent-internship/db"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
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
		//CampaignHandler(ctx, database)
		start := time.Now()
		CampaignHandler(ctx, database)
		elapsed := time.Since(start)
		b.ReportMetric(float64(elapsed.Nanoseconds())/float64(time.Millisecond), "response_time_ms")
	}
}

func BenchmarkCachedCampaignHandler(b *testing.B) {
	database, err := db.ConnectToDB()
	if err != nil {
		b.Fatalf("Failed to connect to database: %v", err)
	}

	err = PreloadData(database)
	if err != nil {
		b.Fatalf("Failed to preload data: %v", err)
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.QueryArgs().Set("source_id", "1")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()
		CachedCampaignHandler(ctx)
		elapsed := time.Since(start)
		b.ReportMetric(float64(elapsed.Nanoseconds())/float64(time.Millisecond), "response_time_ms")
	}
}
