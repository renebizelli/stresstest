/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/renebizelli/stresstest/utils"
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

		fmt.Printf("\nParameters")
		fmt.Printf("\n- url: %s", url)
		fmt.Printf("\n- requests: %d", requests)
		fmt.Printf("\n- concurrency: %d", concurrency)

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

	fmt.Println(utils.BlueText("\n\n\n***************************************"))
	fmt.Println(utils.BlueText("Stress test initializing"))
	fmt.Println(utils.BlueText("***************************************"))
	fmt.Printf("\nParameters:")
	fmt.Printf("\n- url: %s", url)
	fmt.Printf("\n- requests: %d", requests)
	fmt.Printf("\n- concurrency: %d", concurrency)

	//http := NewHttpClient((url))

	if concurrency > requests {
		requests = concurrency
	}

	ws := sync.WaitGroup{}

	results := map[int]int{}

	t := time.Now()

	for requests > 0 {

		loop := concurrency

		if concurrency > requests {
			loop = requests
		}

		ch := make(chan int, loop)

		ws.Add(loop)

		for i := 0; i < loop; i++ {
			go action(ch, &ws)
		}

		ws.Wait()

		for i := 0; i < loop; i++ {
			r := <-ch
			if v, ok := results[r]; !ok {
				results[r] = 1
			} else {
				results[r] = v + 1
			}
		}

		requests -= concurrency
	}

	fmt.Println(utils.GreenText("\n\n***************************************"))
	fmt.Println(utils.GreenText("Stress test results"))
	fmt.Println(utils.GreenText("***************************************"))

	for v, x := range results {
		fmt.Printf("\n- HttpStatus %d: %d returns", v, x)
	}

	//http.Get()

	fmt.Printf("\n\n %s %d requests", utils.YellowText("Total of requests:"), requestCounter)
	fmt.Printf("\n\n %s: %v", utils.YellowText("Operation time"), time.Since(t))

	fmt.Println(utils.BlueText("\n\nStress test completed\n"))

}

func action(ch chan<- int, ws *sync.WaitGroup) {

	atomic.AddInt32(&requestCounter, 1)

	if requestCounter%2 == 0 {
		ch <- 0
	} else {
		ch <- 1
	}

	ws.Done()

}
