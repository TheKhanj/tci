package panel

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/thekhanj/tci"
)

var testPanelStr = `
{
  "phone": {
    "amount": 1850000,
    "payment_url": "https://my.tci.ir/bill/1234567890"
  },
  "internet": {
    "service": "استاندارد کم حجم 50 گیگ ماهانه- قرارداد 12 ماهه- 8Mb پس‌پرداخت",
    "reserved_traffic": -0.01,
    "bundled_traffic": 0.00,
    "days_remaining": 187,
    "days_total": 365,
    "expires_at": "دوشنبه 15 اسفند 1405",
    "traffic_remaining": 42.5,
    "traffic_total": 50
  },
  "account": {
    "online": true,
    "phone": "2112345678",
    "name": "علی محمدی",
    "location": {
      "country": "ایران",
      "province": "اصفهان",
      "city": "اصفهان",
      "exchange": "مرکز مخابرات نمونه"
    },
    "registration": "شنبه 12 اردیبهشت 1404",
    "subscription": "پس پرداخت",
    "credit": "0",
    "mobile": "09123456789",
    "username": "2112345678@tct1",
    "platform": "زوج سیم ADSL",
    "connection": {
      "ip": "192.0.2.42",
      "started": "1405/08/20 09:30",
      "download": "125 MB",
      "upload": "12.5 MB"
    }
  }
}
`

func TestPanelOutput(t *testing.T) {
	var p Panel
	if err := json.Unmarshal([]byte(testPanelStr), &p); err != nil {
		t.Fatal(err)
	}

	t.Run("default", func(t *testing.T) {
		if err := tci.TableWriteTsv(
			os.Stdout,
			p.Columns(tci.ColumnsDefault),
			[]Panel{p},
		); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("wide", func(t *testing.T) {
		if err := tci.TableWriteTsv(
			os.Stdout,
			p.Columns(tci.ColumnsWide),
			[]Panel{p},
		); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("json", func(t *testing.T) {
		if err := tci.TableWriteJson(os.Stdout, []Panel{p}); err != nil {
			t.Fatal(err)
		}
	})
}
