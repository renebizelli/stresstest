package cmd

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/renebizelli/stresstest/utils"
)

type Report struct {
	requestCounter int32
	startedAt      time.Time
	operationTime  time.Duration
	topics         []topic
}

type topic struct {
	title          string
	titleColorFunc func(string) string
	subtitle       string
	items          []string
}

func (r *Report) GetRequestCounter() int32 {
	return r.requestCounter
}

func (r *Report) GetOperationTime() time.Duration {
	return r.operationTime
}

func (r *Report) RequestCounterIncrease() {
	atomic.AddInt32(&r.requestCounter, 1)
}

func (r *Report) Start() {
	r.startedAt = time.Now()
}

func (r *Report) Stop() {
	r.operationTime = time.Since(r.startedAt)
}

func (r *Report) Print() {

	printLine()
	printLine()
	printLine()

	for _, t := range r.topics {
		printTitle(t.title, t.titleColorFunc)
		printSubTitle(t.subtitle)
		printItems(t.items)

		printLine()
		printLine()
	}

	fmt.Println(utils.BlueText("\n\nStress test completed\n"))

	printLine()
	printLine()
}

func (r *Report) AddGreenTopic(title string, subtitle string, items []string) {
	r.addTopic(title, subtitle, items, utils.GreenText)
}

func (r *Report) AddBlueTopic(title string, subtitle string, items []string) {
	r.addTopic(title, subtitle, items, utils.BlueText)
}

func (r *Report) AddYellowTopic(title string, subtitle string, items []string) {
	r.addTopic(title, subtitle, items, utils.YellowText)
}

func (r *Report) addTopic(title string, subtitle string, items []string, color func(string) string) {
	r.topics = append(r.topics, topic{
		title:          title,
		subtitle:       subtitle,
		items:          items,
		titleColorFunc: color,
	})
}

func printLine() {
	fmt.Println("")
}

func printTitle(title string, color func(string) string) {
	if title != "" {
		fmt.Println(color("***************************************"))
		fmt.Println(color(title))
		fmt.Println(color("***************************************"))
	}
}

func printSubTitle(subtitle string) {
	if subtitle != "" {
		fmt.Printf("%s", subtitle)
	}
}

func printItems(items []string) {
	for _, item := range items {
		printLine()
		fmt.Printf("- %s", item)
	}

}
