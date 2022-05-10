package pdf

type name string
type value interface{}

type Text struct {
	OffsetX float64        `json:"offset_x"`
	OffsetY float64        `json:"offset_y"`
	Text    string         `json:"text"`
	Props   map[name]value `json:"props"`
}

func (t *Text) Set(x, y float64) {
	t.OffsetX = x
	t.OffsetY = y
}

func (t *Text) SetX(x float64) {
	t.OffsetX = x
}

func (t *Text) SetY(y float64) {
	t.OffsetY = y
}

func (t *Text) GetColor() (uint8, uint8, uint8) {
	return propFontColor(t)
}

func (t *Text) GetPosition(pageWidth float64) (float64, float64) {
	propBoxWidth(t, pageWidth)
	propAlign(t)

	return t.OffsetX, t.OffsetY
}

func (t *Text) GetFontConfig() (string, string, uint) {
	return propFontFamily(t), propFontStyle(t), propFontSize(t)
}
