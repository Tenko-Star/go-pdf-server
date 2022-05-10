package pdf

//type fontMap []string

type Content struct {
	Width   float64   `json:"width"`
	Height  float64   `json:"height"`
	Pages   []*Page   `json:"pages"`
	Modules moduleMap `json:"modules"`
}
