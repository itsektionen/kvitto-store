package main

import (
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Kvitto struct {
	Timestamp int64         `json:"timestamp"`
	Sold      []SoldProduct `json:"sold"`
	SoldBy    string        `json:"sold_by"`
	Raw       []byte        `json:"-"`
}

type SoldProduct struct {
	Category string `json:"category"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

func (k Kvitto) Time() time.Time {
	return time.Unix(k.Timestamp, 0)
}

func (k Kvitto) toPoints() []*write.Point {
	rawP := influxdb2.NewPoint("kvitto_raw",
		map[string]string{
			"sold_by": strings.ToLower(k.SoldBy),
		},
		map[string]interface{}{
			"json": string(k.Raw),
		},
		k.Time(),
	)

	points := make([]*write.Point, 0, len(k.Sold)+1)
	points = append(points, rawP)

	for _, s := range k.Sold {
		tags := map[string]string{
			"category": s.Category,
			"sold_by":  strings.ToLower(k.SoldBy),
			"product":  s.Name,
		}
		fields := map[string]interface{}{
			"count": s.Count,
		}
		points = append(points, write.NewPoint("kvitto", tags, fields, k.Time()))
	}

	return points
}
