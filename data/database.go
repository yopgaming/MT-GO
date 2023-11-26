// Package data contains all the data related code
package data

import (
	"MT-GO/tools"
	"sync"
)

const (
	databaseLibPath   = "assets"
	coreFilePath      = databaseLibPath + "/core"
	airdropFilePath   = coreFilePath + "/airdrop.json"
	playerScavPath    = coreFilePath + "/playerScav.json"
	MainSettingsPath  = coreFilePath + "/client.settings.json"
	globalsFilePath   = coreFilePath + "/globals.json"
	locationsFilePath = coreFilePath + "/locations.json"
	locationsPath     = databaseLibPath + "/locations"
	matchMetricsPath  = coreFilePath + "/matchMetrics.json"
	serverConfigPath  = coreFilePath + "/server.json"
	editionsDirPath   = databaseLibPath + "/editions/"
	itemsPath         = databaseLibPath + "/items.json"
	localesPath       = databaseLibPath + "/locales"
	handbookPath      = databaseLibPath + "/handbook.json"
	traderPath        = databaseLibPath + "/traders/"
	questsPath        = databaseLibPath + "/quests.json"
	hideoutPath       = databaseLibPath + "/hideout/"
	customizationPath = databaseLibPath + "/customization.json"
	botMainDir        = databaseLibPath + "/bot/"
	botsMainDir       = botMainDir + "bots/"
)

var db *database

type database struct {
	cache         *Cache
	core          *Core
	customization map[string]*Customization
	bot           *Bots
	edition       map[string]*Edition
	template      *Template
	hideout       *Hideout
	item          map[string]*DatabaseItem
	location      *Location
	locale        map[string]*Locale
	profile       map[string]*Profile
	trader        map[string]*Trader
	quest         *Quest
	weather       *Weather
}

func SetPrimaryDatabase() {
	db = &database{
		cache: &Cache{
			player: make(map[string]*PlayerCache),
		},
	}

	var wg sync.WaitGroup
	numWorkers := tools.CalculateWorkers() / 4

	tools.RunTasks(&wg,
		[]func(){
			setBots, setEditions, setHideout,
			setLocales, setTraders,
			setCore, setHandbook,
			setQuests, setItems, setWeather,
			setLocations, setCustomization,
		},
		numWorkers)
}

func SetCache() {
	var wg sync.WaitGroup
	numWorkers := tools.CalculateWorkers() / 4

	tools.RunTasks(&wg,
		[]func(){
			setProfiles, setQuestLookup, setTraderOfferLookup,
			setServerConfig, setHideoutAreaLookup, setHideoutRecipeLookup,
			setScavcaseRecipeLookup, setCachedResponses, setHandbookIndex,
		},
		numWorkers)
}

// #region Database setters

// #endregion

// #region Database structs

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Prefab struct {
	Path string `json:"path"`
	Rcid string `json:"rcid"`
}

type Value struct {
	Value int `json:"value"`
}

type PriceModifier struct {
	PriceModifier float64 `json:"PriceModifier"`
}

type RepairStrategy struct {
	BuffTypes []string `json:"BuffTypes"`
	Filter    []string `json:"Filter"`
}

// #endregion