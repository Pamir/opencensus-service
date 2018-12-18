// Copyright 2018, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample contains a program that exports to the OpenCensus service.
package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

func main() {
	oce, err := ocagent.NewExporter(ocagent.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create ocagent-exporter: %v", err)
	}
	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	// Some stats
        view.SetReportingPeriod(61 * time.Second)
	keyClient, _ := tag.NewKey("client")
	keyMethod, _ := tag.NewKey("method")

        mLatencyMs := stats.Float64("latency", "The latency in milliseconds", "ms")
        latencyView := &view.View{
		Name:        "ocdemo/latency",
		Description: "The various latencies",
		Measure:     mLatencyMs,
		Aggregation: view.Distribution(0, 10, 20, 50, 100, 200, 400, 600, 1000, 1600, 3200, 4800, 5600, 6400),
                TagKeys: []tag.Key{keyClient, keyMethod},
	}
        countView := &view.View{
		Name:        "ocdemo/process_counts",
		Description: "The various counts",
		Measure:     mLatencyMs,
		Aggregation: view.Count(),
                TagKeys: []tag.Key{keyClient, keyMethod},
	}

	// Some metrics
	view.RegisterExporter(oce)
	if err := view.Register(latencyView, countView); err != nil {
		log.Fatalf("Failed to register views for metrics: %v", err)
	}

	ctx, _ := tag.New(context.Background(), tag.Insert(keyMethod, "repl"), tag.Insert(keyClient, "cli"))
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		startTime := time.Now()
		_, span := trace.StartSpan(context.Background(), "Foo")
		time.Sleep(time.Duration(rng.Int63n(3000)) * time.Millisecond)
		span.End()
		latencyMs := float64(time.Since(startTime)) / 1e6
		stats.Record(ctx, mLatencyMs.M(latencyMs))
	}
}
