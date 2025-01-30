package internal

import "time"

type Coin struct {
	Title string
	Rate  float64
	Date  time.Time
}
