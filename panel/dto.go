package panel

type Panel struct {
	Phone    Phone    `json:"phone"`
	Internet Internet `json:"internet"`
	Account  Account  `json:"account"`
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
