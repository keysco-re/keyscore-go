package keyscore

import (
	"encoding/json"
	"io"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type SourceInfo struct {
	Key          string            `json:"Key"`
	DisplayName  string            `json:"DisplayName"`
	AllowedTypes []string          `json:"AllowedTypes"`
	SubSources   map[string]string `json:"SubSources"`
	CompositeOf  []string          `json:"CompositeOf"`
}

type SourcesResponse struct {
	Sources map[string]SourceInfo `json:"sources"`
}

type HashLookupRequest struct {
	Terms []string `json:"terms"`
}

type HashRecord struct {
	Hash      string `json:"hash"`
	Type      string `json:"type"`
	Plaintext string `json:"plaintext"`
	Source    string `json:"source"`
	FirstSeen string `json:"first_seen"`
}

type HashLookupResponse struct {
	Took    int                   `json:"took"`
	Size    int                   `json:"size"`
	Results map[string]HashRecord `json:"results"`
}

type IPLookupRequest struct {
	Terms []string `json:"terms"`
}

type IPInfo struct {
	AS          string  `json:"as"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	ISP         string  `json:"isp"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Org         string  `json:"org"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	Status      string  `json:"status"`
	Timezone    string  `json:"timezone"`
	ZIP         string  `json:"zip"`
}

type IPLookupResponse struct {
	Took    int               `json:"took"`
	Size    int               `json:"size"`
	Results map[string]IPInfo `json:"results"`
	Errors  map[string]string `json:"errors"`
}

type CountRequest struct {
	Terms    []string `json:"terms"`
	Types    []string `json:"types"`
	Source   string   `json:"source,omitempty"`
	Wildcard bool     `json:"wildcard,omitempty"`
	Regex    bool     `json:"regex,omitempty"`
	Operator string   `json:"operator,omitempty"`
	DateFrom string   `json:"dateFrom,omitempty"`
	DateTo   string   `json:"dateTo,omitempty"`
}

type CountResponse struct {
	Count int64 `json:"count"`
}

type DetailedCountResponse struct {
	Counts     map[string]int64 `json:"counts"`
	TotalCount int64            `json:"total_count"`
	Took       int              `json:"took"`
}

type SearchRequest struct {
	Terms    []string `json:"terms"`
	Types    []string `json:"types"`
	Source   string   `json:"source"`
	Wildcard bool     `json:"wildcard,omitempty"`
	Regex    bool     `json:"regex,omitempty"`
	Operator string   `json:"operator,omitempty"`
	DateFrom string   `json:"dateFrom,omitempty"`
	DateTo   string   `json:"dateTo,omitempty"`
	Page     int      `json:"page,omitempty"`
	Pages    any      `json:"pages,omitempty"`
	PageSize int      `json:"pagesize,omitempty"`
}

type SearchResponse struct {
	Results map[string][]map[string]any `json:"results,omitempty"`
	Pages   map[string]any              `json:"pages,omitempty"`
	Size    int                         `json:"size"`
	Took    int                         `json:"took"`
}

type MachineInfo struct {
	BuildID           string   `json:"buildId"`
	IP                string   `json:"ip"`
	UserName          string   `json:"userName"`
	ComputerName      string   `json:"computerName"`
	OperationSystem   string   `json:"operationSystem"`
	Processor         string   `json:"processor"`
	InstalledRAM      string   `json:"installedRAM"`
	GraphicsCard      string   `json:"graphicsCard"`
	Country           string   `json:"country"`
	SystemLanguage    string   `json:"systemLanguage"`
	TimeZone          string   `json:"timeZone"`
	DisplayResolution string   `json:"displayResolution"`
	FileType          string   `json:"fileType"`
	FileTree          []string `json:"fileTree"`
}

// UnmarshalJSON handles the API response which wraps data in a "data" field
// and supports multiple field name variants.
func (m *MachineInfo) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as a wrapper with "data" field
	type wrapper struct {
		Data json.RawMessage `json:"data"`
	}
	var w wrapper
	if err := json.Unmarshal(data, &w); err == nil && len(w.Data) > 0 {
		// If we found a "data" field, unmarshal that instead
		data = w.Data
	}

	// Define a wire struct that includes multiple field name variants
	type wire struct {
		// Build ID variants
		BuildID1 string `json:"buildId"`
		BuildID2 string `json:"BuildID"`
		BuildID3 string `json:"buildid"`

		// IP variants
		IP1 string `json:"ip"`
		IP2 string `json:"IP"`
		IP3 string `json:"ipAddress"`

		// User name variants
		UserName1 string `json:"userName"`
		UserName2 string `json:"UserName"`
		UserName3 string `json:"username"`

		// Computer name variants
		ComputerName1 string `json:"computerName"`
		ComputerName2 string `json:"ComputerName"`
		ComputerName3 string `json:"computername"`

		// Operating system variants
		OperationSystem1 string `json:"operationSystem"`
		OperationSystem2 string `json:"OperationSystem"`
		OperatingSystem1 string `json:"operatingSystem"`
		OperatingSystem2 string `json:"OperatingSystem"`
		OSVersion        string `json:"osVersion"`

		// Processor variants
		Processor1 string `json:"processor"`
		Processor2 string `json:"Processor"`
		CPUName    string `json:"cpuName"`

		// RAM variants
		InstalledRAM1 string `json:"installedRAM"`
		InstalledRAM2 string `json:"InstalledRAM"`
		RAMSize       string `json:"ramSize"`

		// Graphics card variants
		GraphicsCard1 string   `json:"graphicsCard"`
		GraphicsCard2 string   `json:"GraphicsCard"`
		GPUs          []string `json:"gpus"`

		// Country variants
		Country1 string `json:"country"`
		Country2 string `json:"Country"`

		// Language variants
		SystemLanguage1 string `json:"systemLanguage"`
		SystemLanguage2 string `json:"SystemLanguage"`
		Language        string `json:"language"`

		// Timezone variants
		TimeZone1 string `json:"timeZone"`
		TimeZone2 string `json:"TimeZone"`
		TimeZone3 string `json:"timezone"`

		// Display resolution variants
		DisplayResolution1 string `json:"displayResolution"`
		DisplayResolution2 string `json:"DisplayResolution"`
		ScreenResolution   string `json:"screenResolution"`

		// File type variants
		FileType1 string `json:"fileType"`
		FileType2 string `json:"FileType"`

		// File tree
		FileTree []string `json:"fileTree"`
	}

	var wireData wire
	if err := json.Unmarshal(data, &wireData); err != nil {
		return err
	}

	pick := func(vals ...string) string {
		for _, v := range vals {
			if v != "" {
				return v
			}
		}
		return ""
	}

	m.BuildID = pick(wireData.BuildID1, wireData.BuildID2, wireData.BuildID3)
	m.IP = pick(wireData.IP1, wireData.IP2, wireData.IP3)
	m.UserName = pick(wireData.UserName1, wireData.UserName2, wireData.UserName3)
	m.ComputerName = pick(wireData.ComputerName1, wireData.ComputerName2, wireData.ComputerName3)
	m.OperationSystem = pick(wireData.OperationSystem1, wireData.OperationSystem2, wireData.OperatingSystem1, wireData.OperatingSystem2, wireData.OSVersion)
	m.Processor = pick(wireData.Processor1, wireData.Processor2, wireData.CPUName)
	m.InstalledRAM = pick(wireData.InstalledRAM1, wireData.InstalledRAM2, wireData.RAMSize)
	m.Country = pick(wireData.Country1, wireData.Country2)
	m.SystemLanguage = pick(wireData.SystemLanguage1, wireData.SystemLanguage2, wireData.Language)
	m.TimeZone = pick(wireData.TimeZone1, wireData.TimeZone2, wireData.TimeZone3)
	m.DisplayResolution = pick(wireData.DisplayResolution1, wireData.DisplayResolution2, wireData.ScreenResolution)
	m.FileType = pick(wireData.FileType1, wireData.FileType2)
	m.FileTree = wireData.FileTree

	// Handle graphics card - prefer single string, but join array if needed
	m.GraphicsCard = pick(wireData.GraphicsCard1, wireData.GraphicsCard2)
	if m.GraphicsCard == "" && len(wireData.GPUs) > 0 {
		m.GraphicsCard = wireData.GPUs[0] // Take first GPU
	}

	return nil
}

type DownloadResult struct {
	Body               io.ReadCloser
	ContentType        string
	ContentLength      int64
	ContentDisposition string
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type Option func(*Client)

type APIError struct {
	StatusCode int
	Message    string
}
