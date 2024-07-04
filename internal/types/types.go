package types

type FloorListing struct {
	Id         string
	Collection string
	Tick       string
	Price      *Price
	PrevFloor  float64
}

type Price struct {
	Sats     float64
	BTC      float64
	Delta    float64
	DeltaUSD float64
}
