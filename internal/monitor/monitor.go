package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zeroaddresss/golang-unisat-monitor/internal/config"
	"github.com/zeroaddresss/golang-unisat-monitor/internal/errors"
	"github.com/zeroaddresss/golang-unisat-monitor/internal/logger"
	"github.com/zeroaddresss/golang-unisat-monitor/internal/types"
	"github.com/zeroaddresss/golang-unisat-monitor/internal/utils"
	"github.com/zeroaddresss/golang-unisat-monitor/internal/webhook"
)

type RequestPayload struct {
	Filter struct {
		NftType  string `json:"nftType"`
		Tick     string `json:"tick"`
		MinPrice int    `json:"minPrice"`
		IsEnd    bool   `json:"isEnd"`
	} `json:"filter"`
	Sort struct {
		UnitPrice int `json:"unitPrice"`
	} `json:"sort"`
	Start int `json:"start"`
	Limit int `json:"limit"`
}

type RequestData struct {
	Config         *config.Config
	FloorListing   *types.FloorListing
	RequestPayload *RequestPayload
}

// Great tool discovered in the process: https://mholt.github.io/json-to-go/
type ApiResponse struct {
	Code any    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		List []struct {
			AuctionID         string  `json:"auctionId"`
			Amount            float64 `json:"amount"`
			AtomicalID        any     `json:"atomicalId"`
			AtomicalIndex     any     `json:"atomicalIndex"`
			AtomicalTxid      any     `json:"atomicalTxid"`
			Index             any     `json:"index"`
			InscriptionID     string  `json:"inscriptionId"`
			InscriptionNumber int     `json:"inscriptionNumber"`
			Limit             int     `json:"limit"`
			MarketType        string  `json:"marketType"`
			NftType           string  `json:"nftType"`
			Status            string  `json:"status"`
			Tick              string  `json:"tick"`
			Txid              any     `json:"txid"`
			UnitPrice         float64 `json:"unitPrice"`
			Price             float64 `json:"price"`
			Address           string  `json:"address"`
		} `json:"list"`
		Total     int   `json:"total"`
		Timestamp int64 `json:"timestamp"`
	} `json:"data"`
}

func NewRequestData(c *config.Config, idx int) (*RequestData, error) {

	rp := &RequestPayload{
		Filter: struct {
			NftType  string `json:"nftType"`
			Tick     string `json:"tick"`
			MinPrice int    `json:"minPrice"`
			IsEnd    bool   `json:"isEnd"`
		}{
			NftType:  c.Protocol,
			Tick:     c.Collections[idx],
			MinPrice: 0,
			IsEnd:    false,
		},
		Sort: struct {
			UnitPrice int `json:"unitPrice"`
		}{
			UnitPrice: 1,
		},
		Start: 0,
		Limit: 5,
	}

	return &RequestData{
			Config:         c,
			RequestPayload: rp,
			FloorListing: &types.FloorListing{
				Price: &types.Price{},
			},
		},
		nil
}

func StartMonitoring(ctx context.Context, rd *RequestData) error {

	// shuffle api keys
	apiKeyIndex := 0
	var retries int
	tick := rd.RequestPayload.Filter.Tick

	jsonData, err := json.Marshal(rd.RequestPayload)
	if err != nil {
		return errors.NewRequestError("[%s] Error marshalling request data: %v", tick, err)
	}
	client := &http.Client{
		Timeout: time.Duration(rd.Config.Timeout) * time.Millisecond,
	}
	for {
		select {
		case <-ctx.Done():
			// context cancelled
			return errors.NewRequestError("Stopping monitoring for collection: %s", tick)

		default:
			time.Sleep(time.Duration(rd.Config.Delay) * time.Millisecond)

			// each iteration requires a new request body because it's consumed by the request
			payload := bytes.NewReader(jsonData)
			if payload == nil {
				return errors.NewRequestError("[%s] Error marshalling request data: %v", tick, err)
			}

			req, err := http.NewRequest("POST", rd.Config.MonitoringURL, payload)
			if err != nil {
				return errors.NewRequestError("[%s] Error creating request: %v", tick, err)
			}

			apiKeyIndex++
			if apiKeyIndex >= len(rd.Config.ApiKeys) {
				apiKeyIndex = 0
			}
			req.Header = http.Header{
				"Authorization": []string{"Bearer " + rd.Config.ApiKeys[apiKeyIndex]},
				"Accept":        []string{"application/json"},
				"Content-Type":  []string{"application/json"},
			}
			res, err := client.Do(req)
			if err != nil {
				return errors.NewRequestError("[%s] Error sending request: %s [%s]", tick, err, utils.Red(res.StatusCode))
			}

			if res.StatusCode != 200 {
				// most likely because of API free usage exceeded (10,000 requests per day allowed)
				logger.L.Error("[%s] [%s] Request unsuccesful: %s", utils.Red(res.StatusCode), tick, res.Status)
				if retries < rd.Config.MaxRetries {
					retries++
					continue
				}
				return errors.NewRequestError("[%s] [%s] Max retries exceeded: %s", tick, utils.Red(res.StatusCode), res.Status)
			}

			var responseData ApiResponse

			if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
				logger.L.Error("[%s] [%s] Error decoding response: %v", tick, utils.Red(res.StatusCode), err)
				if retries < rd.Config.MaxRetries {
					retries++
					continue
				}
				return errors.NewRequestError("[%s] [%s] Max retries exceeded (couldn't decode response): %s", tick, utils.Red(res.StatusCode), res.Status)
			}
			res.Body.Close()
			listings := responseData.Data.List

			if len(listings) == 0 {
				logger.L.Error("[%s] Response returned no listings [%s]", tick, utils.Red(res.StatusCode))
				continue
			}
			cheapestListing := listings[0]
			precision := float64(1e8) // 1 satoshi
			prevFloor := rd.FloorListing.Price.Sats
			prevFloorBTC := rd.FloorListing.Price.Sats / precision
			currentFloor := cheapestListing.UnitPrice
			currentFloorBTC := cheapestListing.Price / precision

			// initialize floor values (first iteration)
			if rd.FloorListing.Price.Sats == 0 {
				rd.FloorListing.Id = cheapestListing.InscriptionID
				rd.FloorListing.Collection = cheapestListing.Address
				rd.FloorListing.Price.Sats = cheapestListing.UnitPrice
				rd.FloorListing.Price.BTC = cheapestListing.Price / precision
				rd.FloorListing.PrevFloor = cheapestListing.UnitPrice
				rd.FloorListing.Tick = cheapestListing.Tick
				logger.L.Info("[%s] Current floor %s sats [%s]", tick, logger.FormatFloat(rd.FloorListing.Price.Sats), utils.Green(res.StatusCode))
				continue
			}

			if currentFloor == prevFloor {
				// no changes
				logger.L.Info("[%s] Current floor %s sats [%s]", tick, logger.FormatFloat(currentFloor), utils.Green(res.StatusCode))
				continue
			}

			// floor has changed
			delta := ((prevFloor - currentFloor) / prevFloor) * float64(100)
			rd.FloorListing.Id = cheapestListing.InscriptionID
			rd.FloorListing.Collection = cheapestListing.Address
			rd.FloorListing.Tick = cheapestListing.Tick
			rd.FloorListing.Price.Sats = cheapestListing.UnitPrice // Sats per unit
			rd.FloorListing.Price.BTC = cheapestListing.Price / precision
			rd.FloorListing.PrevFloor = prevFloor
			rd.FloorListing.Price.Delta = delta

			if delta < 5 {
				// doesn't notify for small changes (less than 5% price drop)
				var deltaStr string
				if delta < 0 {
					// floor price has increased
					deltaStr = fmt.Sprintf("+%.3f %%", -delta)
				} else {
					deltaStr = fmt.Sprintf("-%.3f %%", delta)
				}
				logger.L.Info("[%s] New floor %s sats (%s). Skipping webhook [%s]", tick, logger.FormatFloat(rd.FloorListing.Price.Sats), deltaStr, utils.Green(res.StatusCode))
				continue
			}
			logger.L.Info("[%s] New floor %s sats (-%.3f%%). Sending webhook [%s]", tick, logger.FormatFloat(rd.FloorListing.Price.Sats), delta, utils.Green(res.StatusCode))

			rd.FloorListing.Price.DeltaUSD = (prevFloorBTC - currentFloorBTC) * 63000.0 // convert BTC to USD
			logger.L.Info("[%s] Price diff: %.2f USD", tick, rd.FloorListing.Price.DeltaUSD)

			// new floor listing detected, notify discord
			embed := webhook.BuildEmbed(rd.FloorListing, rd.Config)
			for name, url := range rd.Config.Webhooks {
				if err := webhook.SendWebhook(ctx, embed, url); err != nil {
					logger.L.Error("[%s] Error sending webhook: %v", tick, err)
					continue
				}
				logger.L.Info("[%s] Webhook sent to %s", tick, name)
			}
		}
	}
}
