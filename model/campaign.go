package model

import (
	"fmt"
	"math/rand"
	"time"
)

type Campaign struct {
	ID         int
	Name       string
	FilterType string
	Domains    map[string]bool
	SourceID   int
}

type CampaignWithPrice struct {
	Campaign Campaign
	Price    int
}

func (c Campaign) Call() CampaignWithPrice {
	price := rand.Intn(100) + 1
	fmt.Printf("Campaign %d generated price: %d\n", c.ID, price)

	delay := rand.Intn(5)
	fmt.Printf("Campaign %d simulating response delay: %d seconds\n", c.ID, delay)
	time.Sleep(time.Duration(delay) * time.Second)

	return CampaignWithPrice{Campaign: c, Price: price}
}

type CampaignSlice struct {
	ID         int
	Name       string
	FilterType string
	Domains    []string
	SourceID   int
}

type CampaignSliceWithPrice struct {
	CampaignSlice CampaignSlice
	Price         int
}

func (c CampaignSlice) Call() CampaignSliceWithPrice {
	price := rand.Intn(100) + 1
	fmt.Printf("Campaign %d generated price: %d\n", c.ID, price)

	delay := rand.Intn(5)
	fmt.Printf("Campaign %d simulating response delay: %d seconds\n", c.ID, delay)
	time.Sleep(time.Duration(delay) * time.Second)

	return CampaignSliceWithPrice{CampaignSlice: c, Price: price}
}
