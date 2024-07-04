package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zeroaddresss/golang-unisat-monitor/internal/config"
	"github.com/zeroaddresss/golang-unisat-monitor/internal/logger"
	"github.com/zeroaddresss/golang-unisat-monitor/internal/monitor"
	"github.com/zeroaddresss/golang-unisat-monitor/internal/utils"
)

func main() {
	utils.PrintFancyText("Unisat Monitor")

	logsDirectory := "./logs/"
	logger.L = logger.NewLogger(logsDirectory)
	defer logger.L.LogFile.Close()

	configFilePath := "./internal/config/config.json"
	c, err := config.LoadConfig(configFilePath)
	if err != nil {
		logger.L.Error("%v", err)
		return
	}
	if err = config.ValidateConfig(c); err != nil {
		logger.L.Error("%v", err)
		return
	}

	// TODO: Runes support
	if c.Protocol == "brc20" {
		c.MonitoringURL = fmt.Sprintf("https://open-api.unisat.io/v3/market/%s/auction/list", c.Protocol)
		c.CollectionBaseURL = "https://unisat.io/market/brc20?tick="
		c.ListingBaseURL = "https://unisat.io/inscription/"
	} else {
		logger.L.Error("Only BRC20 supported for now: '%s' invalid", c.Protocol)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	for i := range c.Collections {
		wg.Add(1)
		// Each collection is monitored in its own goroutine
		go func(collectionIndex int) {
			defer wg.Done()

			reqData, _ := monitor.NewRequestData(c, collectionIndex)
			if err := monitor.StartMonitoring(ctx, reqData); err != nil {
				logger.L.Error("%v", err)
			}
		}(i)
	}

	// Handle SIGINT and SIGTERM for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.L.Info("Termination signal received, shutting down...")

	// Cancel the context to signal shutdown to the goroutines
	cancel()
	wg.Wait()
	logger.L.Info("All monitoring routines completed")
}
