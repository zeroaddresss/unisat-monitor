package webhook

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
)

type Message struct {
	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string    `json:"title,omitempty"`
	Url         string    `json:"url,omitempty"`
	Description string    `json:"description,omitempty"`
	Color       int       `json:"color,omitempty"`
	Timestamp   string    `json:"timestamp,omitempty"`
	Author      Author    `json:"author,omitempty"`
	Fields      []Field   `json:"fields,omitempty"`
	Thumbnail   Thumbnail `json:"thumbnail,omitempty"`
	Image       Image     `json:"image,omitempty"`
	Footer      Footer    `json:"footer,omitempty"`
}

type Author struct {
	Name    string `json:"name,omitempty"`
	Url     string `json:"url,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}

type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
	Hidden bool   `json:"-"`
}

type Thumbnail struct {
	Url string `json:"url,omitempty"`
}

type Image struct {
	Url string `json:"url,omitempty"`
}

type Footer struct {
	Text    string `json:"text,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}

/////////////////////
/* EMBED FUNCTIONS */
/////////////////////

func NewDiscordEmbed(title, description string) *Embed {
	return &Embed{
		Title:       title,
		Description: description,
		Timestamp:   timestamp(),
	}
}

func (e *Embed) AddEmbedField(name, value string, inline bool) {
	e.Fields = append(e.Fields, Field{Name: name, Value: value, Inline: inline})
}

func (e *Embed) SetFooter(text, iconURL string) {
	e.Footer = Footer{Text: text, IconUrl: iconURL}
}

func (e *Embed) SetThumbnail(url string) {
	e.Thumbnail = Thumbnail{Url: url}
}

func (e *Embed) SetImage(url string) {
	e.Image = Image{Url: url}
}

func (e *Embed) ToJSON() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func timestamp() string {
	return time.Now().Local().Format(time.RFC3339)
}

func BuildEmbed(fl *types.FloorListing, c *config.Config) *Embed {
	title := "Discounted Listing Detected 👀"
	description := fmt.Sprintf("**%.2f%% off floor** 🚀", fl.Price.Delta)
	embed := NewDiscordEmbed(title, description)
	embed.Color = 16711680 // Red
	embed.AddEmbedField("Collection", fmt.Sprintf("[**%s**](%s%s)", fl.Tick, c.CollectionBaseURL, fl.Tick), false)

	embed.AddEmbedField("Old Floor 💸", fmt.Sprintf("**%s sats**", logger.FormatFloat(fl.PrevFloor)), true)
	embed.AddEmbedField("New Listing 🤑", fmt.Sprintf("**%s sats**", logger.FormatFloat(fl.Price.Sats)), true)
	embed.AddEmbedField("Price 💰", fmt.Sprintf("**%.5f BTC **", fl.Price.BTC), false)
	embed.AddEmbedField("Delta 📊", fmt.Sprintf("**%.2f USD**", fl.Price.DeltaUSD), true)
	embed.AddEmbedField(" ", fmt.Sprintf("[**Buy Now 🛒**](%s%s)", c.ListingBaseURL, fl.Id), false)
	embed.SetFooter("@zeroaddresss", "https://m.media-amazon.com/images/I/61nGO5El1BL._AC_UF894,1000_QL80_.jpg")
	return embed
}

func SendWebhook(ctx context.Context, embed *Embed, webhookUrl string) error {
	client := &http.Client{}
	message := &Message{Embeds: []Embed{*embed}}
	messageData, err := json.Marshal(message)
	if err != nil {
		logger.L.Error("Error marshalling webhook embed: %v", err)
	}
	payload := bytes.NewReader(messageData)

	req, err := http.NewRequestWithContext(
		ctx, "POST", webhookUrl, payload,
	)
	if err != nil {
		return errors.NewRequestError("Error creating webhook request: %v", err)
	}
	req.Header = http.Header{
		"content-type": {"application/json"},
	}
	res, err := client.Do(req)
	if err != nil {
		return errors.NewRequestError("Request error: %v", err)
	}
	if res.StatusCode != 204 {
		return errors.NewRequestError("[%d] Webhook request unsuccessful [%s]", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	return nil
}
