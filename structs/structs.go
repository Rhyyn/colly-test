package structs

import "os"

// region ScrapedItem
type Item struct {
	ID                 int                `json:"id"`
	Title              Title              `json:"title"`
	Params             Params             `json:"params"`
	Stats              map[int]Stat       `json:"stats"`
	Droprates          map[int]Droprate   `json:"droprates"`
	SublimationDetails SublimationDetails `json:"sublimations_details,omitempty"`
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
	Format      string  `json:"format,omitempty"`
	IsNegative  bool    `json:"is_negative,omitempty"`
	Value       int     `json:"value,omitempty"`
	NumElements int     `json:"num_elements,omitempty"`
}

type Droprate struct {
	MonsterID   int     `json:"monster_id"`
	MonsterName Display `json:"monster_name"`
	DropChance  float64 `json:"drop_chance"`
}

type Recipe struct {
	JobID       int                `json:"job_id"`
	JobLevel    int                `json:"job_level"`
	JobName     Display            `json:"job_name"`
	RecipeId    int                `json:"recipe_id"`
	ResultId    int                `json:"result_id"`
	Ingredients map[int]Ingredient `json:"ingredients"`
}

type Ingredient struct {
	ID       int     `json:"id"`
	TypeID   int     `json:"type_id"`
	Quantity int     `json:"quantity"`
	IngName  Display `json:"ing_name"`
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

// region itemTypes
type TypesDefinition struct {
	ID       int `json:"id"`
	ParentID int `json:"parentId"`
}

type TypesTitle struct {
	Fr string `json:"fr"`
	En string `json:"en"`
}

type TypesItem struct {
	Definition TypesDefinition `json:"definition"`
	Title      TypesTitle      `json:"title"`
}

// endregion itemTypes

type ScrapingParameters struct {
	SingleItemURL  map[string]string
	IndexUrl       map[string]string
	ItemUrl        map[string]string
	MaxPage        int
	MaxItems       int
	SelectedId     int
	SingleItemMode bool
	SelectedType   string
}

// region Ressources

type SublimationDetails struct {
	IsEpic        bool   `json:"is_epic"`
	LevelInc      int    `json:"level_inc"`
	SocketsColors string `json:"sockets_colors"`
	MaxLevel      int    `json:"max_level"`
	Desc          Title  `json:"desc"`
	ID            int    `json:"id"`
}

// endregion Ressources

type ScrapedItem struct {
	ID int `json:"_id"`
}

type Type struct {
	Fr string `json:"fr"`
	En string `json:"en"`
}

type UrlArgs struct {
	Fr string `json:"fr"`
	En string `json:"en"`
}

type ItemTypes struct {
	ID      []int             `json:"id"`
	Title   map[string]string `json:"title"`
	Type    Type              `json:"type"`
	UrlArgs UrlArgs           `json:"url_args"`
}
type SubCategory struct {
	Title     map[string]string `json:"title"`
	ID        []int             `json:"id"`
	Index_url map[string]string `json:"index_url"`
	Item_url  map[string]string `json:"item_url"`
	MaxPage   int               `json:"max_page"`
	MaxItems  int               `json:"max_items"`
	ItemTypes []ItemTypes       `json:"items_types"`
}

// Move this to structs
type ItemInfo struct {
	Title         map[string]string `json:"title"`
	ID            []int             `json:"id"`
	Index_url     map[string]string `json:"index_url"`
	Item_url      map[string]string `json:"item_url"`
	MaxPage       int               `json:"max_page"`
	SubCategories []SubCategory     `json:"sub_categories"`
}

type EditFileOptions struct {
	IsSubCat bool
	ID       int
	SubCat   SubCategory
}

type FileResult struct {
	File *os.File
}
