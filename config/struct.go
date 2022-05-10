package config

type Config struct {
	Roots    `json:"roots"`
	Default  `json:"default"`
	Server   `json:"server"`
	Callback `json:"callback"`

	Templates []Template `json:"template"`
	Fonts     []Font     `json:"fonts"`
	FontMap   map[string]string

	KeepTemp bool `json:"keep_temp"`

	MaxQueue  uint `json:"max_queue"`
	MaxThread uint `json:"max_thread"`

	Debug bool `json:"debug"`
}

type Roots struct {
	Output    string `json:"output"`
	Template  string `json:"template"`
	Fonts     string `json:"fonts"`
	Log       string `json:"log"`
	Temporary string `json:"temporary"`
}

type Default struct {
	PageWidth  float64 `json:"page_width"`
	PageHeight float64 `json:"page_height"`
	Font       string  `json:"font"`
	FontSize   uint    `json:"font_size"`
}

type Template struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Server struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	Charset string `json:"charset"`
}

type Callback struct {
	Success string `json:"success"`
	Failure string `json:"failure"`
	Upload  string `json:"upload"`
}

type Font struct {
	Family string `json:"family"`
	Path   string `json:"path"`
}
