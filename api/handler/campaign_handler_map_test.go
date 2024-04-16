package handler_test

import (
	"adtelligent-internship/api"
	"adtelligent-internship/api/handler"
	"adtelligent-internship/api/repository"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	_ "net/http/pprof"
	"testing"
	"time"
)

func BenchmarkCampaignHandler(b *testing.B) {
	repository.PreloadedCache = make(map[int][]repository.Campaign)
	repository.PreloadedCache[123] = []repository.Campaign{
		{ID: 1, Name: "Campaign 1", FilterType: "black", Domains: map[string]bool{"example.com": true}, SourceID: 123},
		{ID: 2, Name: "Campaign 2", FilterType: "white", Domains: map[string]bool{"example.com": false}, SourceID: 123},
	}

	reqCtx := &fasthttp.RequestCtx{}

	reqCtx.QueryArgs().Set("source_id", "1")
	reqCtx.QueryArgs().Set("domain", "example.com")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()
		handler.CampaignHandlerMap(reqCtx)
		elapsed := time.Since(start)
		b.ReportMetric(float64(elapsed.Nanoseconds())/float64(time.Millisecond), "response_time_ms")
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
	ctx.Request.SetRequestURI("/campaigns_map")
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

type Campaign struct {
	ID         int
	Name       string
	FilterType string
	Domains    map[string]bool
	SourceID   int
}

type CampaignSlice struct {
	ID         int
	Name       string
	FilterType string
	Domains    []string
	SourceID   int
}

var (
	mockLargeCampaignData = generateLargeMockCampaignData()
)

var PreloadedCache = make(map[int][]Campaign)
var PreloadedSliceCache = make(map[int][]CampaignSlice)

func generateLargeMockCampaignData() map[int][]Campaign {
	mockData := make(map[int][]Campaign)

	for sourceID := 1; sourceID <= 100; sourceID++ {
		var campaigns []Campaign
		for campaignID := 1; campaignID <= 10; campaignID++ {
			campaign := Campaign{
				ID:         campaignID,
				Name:       fmt.Sprintf("Campaign %d", campaignID),
				FilterType: "white",
				Domains:    map[string]bool{"example.com": true, "test.com": true},
				SourceID:   sourceID,
			}
			campaigns = append(campaigns, campaign)
		}
		mockData[sourceID] = campaigns
	}
	return mockData
}

func BenchmarkMapLookupLargeData(b *testing.B) {
	PreloadedCache = mockLargeCampaignData

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sourceID := 99
		_ = LookupCampaignsBySourceIDMap(sourceID)
	}
}

func BenchmarkSliceLookupLargeData(b *testing.B) {
	for sourceID, campaigns := range mockLargeCampaignData {
		var campaignSlices []CampaignSlice
		for _, campaign := range campaigns {
			campaignSlices = append(campaignSlices, CampaignSlice{
				ID:         campaign.ID,
				Name:       campaign.Name,
				FilterType: campaign.FilterType,
				Domains:    campaign.getDomainSlice(),
				SourceID:   campaign.SourceID,
			})
		}
		PreloadedSliceCache[sourceID] = campaignSlices
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sourceID := 99
		_ = LookupCampaignsBySourceIDSlice(sourceID)
	}
}

func LookupCampaignsBySourceIDMap(sourceID int) []Campaign {
	if campaigns, ok := PreloadedCache[sourceID]; ok {
		return campaigns
	}
	return nil
}

func LookupCampaignsBySourceIDSlice(sourceID int) []CampaignSlice {
	if campaigns, ok := PreloadedSliceCache[sourceID]; ok {
		return campaigns
	}
	return nil
}

func (c *Campaign) getDomainSlice() []string {
	var domains []string
	for domain := range c.Domains {
		domains = append(domains, domain)
	}
	return domains
}
