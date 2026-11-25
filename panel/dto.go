package panel

import (
	"fmt"
	"strconv"

	"github.com/thekhanj/tci"
)

type Panel struct {
	Phone    Phone    `json:"phone"`
	Internet Internet `json:"internet"`
	Account  Account  `json:"account"`
}

var _ tci.Table = (*Panel)(nil)

func (Panel) Headers(format tci.OutputFormat) []string {
	header := []string{
		"number",
		"status",
		"days",
		"traffic",
		"phone_bill",
	}

	if format == tci.OutputWide {
		header = []string{
			"number",
			"status",
			"ip",
			"days",
			"traffic",
			"expires",
			"phone_bill",
			"phone_bill_url",
			"subscription",
		}
	}

	return header
}

func (p Panel) Values(fields []string) []string {
	usedDays := p.Internet.DaysTotal - p.Internet.DaysRemaining

	daysPercent := 0
	if p.Internet.DaysTotal > 0 {
		daysPercent = usedDays * 100 / p.Internet.DaysTotal
	}

	days := fmt.Sprintf("%d(used)/%d (%d%%)",
		usedDays, p.Internet.DaysTotal, daysPercent)

	traffic := fmt.Sprintf("%g/%g",
		p.Internet.TrafficRemaining, p.Internet.TrafficTotal)

	status := "offline"
	if p.Account.Online {
		status = "online"
	}

	values := map[string]tci.FieldValue{
		"number":         tci.NewFieldValue(p.Account.Username),
		"status":         tci.NewFieldValue(status),
		"ip":             tci.NewFieldValue(p.Account.Connection.IP),
		"days":           tci.NewFieldValue(days),
		"traffic":        tci.NewFieldValue(traffic),
		"expires":        tci.NewFieldValue(p.Internet.ExpiresAt),
		"phone_bill":     tci.NewFieldValue(strconv.FormatInt(p.Phone.Amount, 10)),
		"phone_bill_url": tci.NewFieldValue(p.Phone.PaymentURL),
		"subscription":   tci.NewFieldValue(p.Account.Subscription),
	}

	return tci.SelectFields(fields, values)
}

type Phone struct {
	Amount     int64  `json:"amount"`
	PaymentURL string `json:"payment_url"`
}

type Internet struct {
	Service          string  `json:"service"`
	ReservedTraffic  float64 `json:"reserved_traffic"`
	BundledTraffic   float64 `json:"bundled_traffic"`
	DaysRemaining    int     `json:"days_remaining"`
	DaysTotal        int     `json:"days_total"`
	ExpiresAt        string  `json:"expires_at"`
	TrafficRemaining float64 `json:"traffic_remaining"`
	TrafficTotal     float64 `json:"traffic_total"`
}

type Account struct {
	Online       bool       `json:"online"`
	Phone        string     `json:"phone"`
	Name         string     `json:"name"`
	Location     Location   `json:"location"`
	Registration string     `json:"registration"`
	Subscription string     `json:"subscription"`
	Credit       string     `json:"credit"`
	Mobile       string     `json:"mobile"`
	Username     string     `json:"username"`
	Platform     string     `json:"platform"`
	Connection   Connection `json:"connection"`
}

type Location struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Exchange string `json:"exchange"`
}

type Connection struct {
	IP       string `json:"ip"`
	Started  string `json:"started"`
	Download string `json:"download"`
	Upload   string `json:"upload"`
}
