/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/spf13/cobra"
)

// stressOutCmd represents the stressOut command
var stressOutCmd = &cobra.Command{
	Use:   "stressOut",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		url, _ := cmd.Flags().GetString("url")
		requests, _ := cmd.Flags().GetInt("requests")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		stressOut(url, requests, concurrency)
	},
}

func init() {
	rootCmd.AddCommand(stressOutCmd)

	stressOutCmd.Flags().StringP("url", "u", "", "Uri will be tested")
	stressOutCmd.Flags().IntP("requests", "r", 0, "Total number os requests")
	stressOutCmd.Flags().IntP("concurrency", "c", 0, "Simultaneous requests")

	stressOutCmd.MarkFlagRequired("url")
	stressOutCmd.MarkFlagRequired("requests")
	stressOutCmd.MarkFlagRequired("concurrency")
}

var requestCounter int32

func stressOut(url string, requests int, concurrency int) {

	report := &Report{}

	http := NewHttpClient(url)

	if concurrency > requests {
		requests = concurrency
	}

	ws := sync.WaitGroup{}

	results := map[int]int{}

	report.Start()

	control := requests

	for control > 0 {

		loop := concurrency

		if concurrency > control {
			loop = control
		}

		ch := make(chan int, loop)

		ws.Add(loop)

		for i := 0; i < loop; i++ {
			go action(http, report, ch, &ws)
		}

		ws.Wait()

		for range loop {
			r := <-ch
			if v, ok := results[r]; !ok {
				results[r] = 1
			} else {
				results[r] = v + 1
			}
		}

		control -= concurrency
	}

	report.Stop()

	report.AddBlueTopic(
		"Stress test initializing",
		"Parameters",
		[]string{
			fmt.Sprintf("url: %s", url),
			fmt.Sprintf("requests: %d", requests),
			fmt.Sprintf("concurrency: %d", concurrency),
		},
	)

	itemsOfResults := []string{}
	for v, x := range results {
		itemsOfResults = append(itemsOfResults, fmt.Sprintf("HttpStatus %d: %d returns", v, x))
	}

	report.AddGreenTopic(
		"Stress test results",
		"",
		itemsOfResults,
	)

	report.AddYellowTopic(
		"Test summary",
		"",
		[]string{
			fmt.Sprintf("Total of requests: %d requests", report.GetRequestCounter()),
			fmt.Sprintf("Operation time: %v", report.GetOperationTime()),
		},
	)

	report.Print()
}

func action(action ActionInterface, report *Report, ch chan<- int, ws *sync.WaitGroup) {

	result := action.Get()

	report.RequestCounterIncrease()

	atomic.AddInt32(&requestCounter, 1)

	ch <- result

	ws.Done()
}
