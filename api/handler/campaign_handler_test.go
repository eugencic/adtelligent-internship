package handler_test

import (
	"adtelligent-internship/api"
	"adtelligent-internship/api/handler"
	"adtelligent-internship/api/repository"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
)

func BenchmarkCampaignHandler(b *testing.B) {
	reqCtx := &fasthttp.RequestCtx{}

	reqCtx.QueryArgs().Set("source_id", "1")
	reqCtx.QueryArgs().Set("domain", "example.com")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//start := time.Now()
		handler.CampaignHandler(reqCtx)
		//elapsed := time.Since(start)
		//b.ReportMetric(float64(elapsed.Nanoseconds())/float64(time.Millisecond), "response_time_ms")
	}
}

func TestNewRequestHandler_CampaignsEndpoint(t *testing.T) {
	mockCache := map[int][]repository.Campaign{
		1: {
			{ID: 2, Name: "Campaign 2", FilterType: "white", Domains: map[string]bool{"aol.com": true, "gmail.com": true, "hotmail.com": true, "msn.com": true, "yahoo.com": true}, SourceID: 1},
			{ID: 77, Name: "Campaign 77", FilterType: "white", Domains: map[string]bool{"msn.com": true, "orange.fr": true, "yahoo.com": true}, SourceID: 1},
			{ID: 90, Name: "Campaign 90", FilterType: "black", Domains: map[string]bool{"aol.com": true, "gmail.com": true, "msn.com": true, "yahoo.com": true, "yandex.ru": true}, SourceID: 1},
		},
	}

	repository.PreloadedCache = mockCache

	requestHandler := api.NewRequestHandler()

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/campaigns")
	ctx.Request.Header.SetMethod("GET")
	ctx.QueryArgs().Set("source_id", "1")
	ctx.QueryArgs().Set("domain", "hotmail.com")

	requestHandler(ctx)

	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())

	var parsedCampaigns []repository.Campaign
	err := json.Unmarshal(ctx.Response.Body(), &parsedCampaigns)
	assert.NoError(t, err)

	expectedResult := []repository.Campaign{
		{ID: 2, Name: "Campaign 2", FilterType: "white", Domains: map[string]bool{"aol.com": true, "gmail.com": true, "hotmail.com": true, "msn.com": true, "yahoo.com": true}, SourceID: 1},
		{ID: 90, Name: "Campaign 90", FilterType: "black", Domains: map[string]bool{"aol.com": true, "gmail.com": true, "msn.com": true, "yahoo.com": true, "yandex.ru": true}, SourceID: 1},
	}

	assert.ElementsMatch(t, expectedResult, parsedCampaigns)
}
