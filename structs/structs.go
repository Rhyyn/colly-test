package structs

// region ScrapedItem
type Item struct {
	ID        int              `json:"id"`
	Title     Title            `json:"title"`
	Params    Params           `json:"params"`
	Stats     map[int]Stat     `json:"stats"`
	Droprates map[int]Droprate `json:"droprates"`
	Recipes   []Recipe         `json:"recipes"`
}

type Params struct {
	TypeId int `json:"type_id"`
	Level  int `json:"level"`
	Rarity int `json:"rarity"`
}

type Title struct {
	Fr string `json:"fr"`
	En string `json:"en"`
}

type Stat struct {
	Display     Display `json:"display"`
	ID          int     `json:"id"`
	Format      string  `json:"format"`
	Value       int     `json:"value"`
	NumElements int     `json:"num_elements,omitempty"`
}

type Droprate struct {
	MonsterID   int     `json:"monster_id"`
	MonsterName Display `json:"monster_name"`
	DropChance  float64 `json:"drop_chance"`
}

type Recipe struct {
	JobID       int          `json:"job_id"`
	JobLevel    int          `json:"job_level"`
	Ingredients []Ingredient `json:"ingredients"`
}

type Ingredient struct {
	ID       int `json:"id"`
	TypeID   int `json:"type_id"`
	Quantity int `json:"quantity"`
	GfxID    int `json:"gfx_id"`
}

// endregion

// region Stat
type Display struct {
	Fr string `json:"fr"`
	En string `json:"en"`
}

// endregion Stat

type StatProperties struct {
	Fr string `json:"fr"`
	En string `json:"en"`
}

type ParamsStatsProperties struct {
	AllPositivesStats map[string]StatProperties
	AllNegativesStats map[string]StatProperties
}
