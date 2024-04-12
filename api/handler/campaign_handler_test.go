package handler_test

import (
	"adtelligent-internship/api/handler"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func BenchmarkCampaignHandler(b *testing.B) {
	reqCtx := &fasthttp.RequestCtx{}

	reqCtx.QueryArgs().Set("source_id", "1")
	reqCtx.QueryArgs().Set("domain", "example.com")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()
		handler.CampaignHandler(reqCtx)
		elapsed := time.Since(start)
		b.ReportMetric(float64(elapsed.Nanoseconds())/float64(time.Millisecond), "response_time_ms")
	}
}
