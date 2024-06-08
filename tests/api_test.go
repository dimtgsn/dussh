package tests

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const baseURL = "http://localhost:8082/"

func TestGetCourseMethod(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 10 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    path.Join(baseURL, "courses/31"),
	})

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewHDRHistogramPlotReporter(&metrics)

	file, err := os.Create("get_section")
	if err != nil {
		t.Fatal(err)
	}

	if err := reporter.Report(file); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func TestPostCourseMethod(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 10 * time.Second
	body := fmt.Sprintf("{\"name\": \"swimming%d\",    \"monthly_subscription_cost\": 500,    \"events\": [        {            \"description\": \"2222222222\",            \"start_date\": \"2024-10-15 14:00:00\",            \"recurrent_count\": 60,            \"period_freq\": 1,            \"period_type\": \"day\"        }    ],    \"employees\": [9]}",
		rand.Int())
	b, _ := json.Marshal(&body)
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		Body:   b,
		URL:    path.Join(baseURL, "courses"),
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewHDRHistogramPlotReporter(&metrics)

	file, err := os.Create("post_section")
	if err != nil {
		t.Fatal(err)
	}

	if err := reporter.Report(file); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func TestPostUserMethod(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 10 * time.Second
	body := fmt.Sprintf("{\"name\": \"swimming%d\",    \"monthly_subscription_cost\": 500,    \"events\": [        {            \"description\": \"2222222222\",            \"start_date\": \"2024-10-15 14:00:00\",            \"recurrent_count\": 60,            \"period_freq\": 1,            \"period_type\": \"day\"        }    ],    \"employees\": [9]}",
		rand.Int())
	b, _ := json.Marshal(&body)
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		Body:   b,
		URL:    path.Join(baseURL, "users"),
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewHDRHistogramPlotReporter(&metrics)

	file, err := os.Create("post_user")
	if err != nil {
		t.Fatal(err)
	}

	if err := reporter.Report(file); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func TestGetAllCoursesMethod(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 10 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "PUT",
		URL:    path.Join(baseURL, "courses"),
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewHDRHistogramPlotReporter(&metrics)

	file, err := os.Create("get_all_sections")
	if err != nil {
		t.Fatal(err)
	}

	if err := reporter.Report(file); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func TestGetUserMethod(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 10 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    path.Join(baseURL, "users/1"),
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewHDRHistogramPlotReporter(&metrics)

	file, err := os.Create("get_user")
	if err != nil {
		t.Fatal(err)
	}

	if err := reporter.Report(file); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func TestGetAllUsersMethod(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 10 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "PUT",
		URL:    path.Join(baseURL, "users"),
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewHDRHistogramPlotReporter(&metrics)

	file, err := os.Create("get_all_users")
	if err != nil {
		t.Fatal(err)
	}

	if err := reporter.Report(file); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}
