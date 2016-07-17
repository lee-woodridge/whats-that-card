package types

type CardSets map[string][]Card

type Mechanic struct {
	Name string
}

type Card struct {
	Mechanics    []Mechanic // needs it's own type [{"name": "Charge"},{"name": "Divine Shield"}]
	Durability   int
	Locale       string
	Text         string
	HowToGet     string
	ImgGold      string
	Cost         int
	Flavor       string
	PlayerClass  string
	Img          string
	Attack       int
	Health       int
	Type         string
	Collectible  bool
	Faction      string
	InPlayText   string
	Elite        bool
	HowToGetGold string
	CardSet      string
	Name         string
	Artist       string
	Rarity       string
	Race         string
	CardId       string
}
