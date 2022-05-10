package pdf

type Page struct {
	Values []*Value `json:"values"`
}

func (p *Page) Text(text string) *Text {
	t := &Text{
		Text: text,
	}

	p.Values = append(p.Values, &Value{
		T: TEXT,
		V: t,
	})

	return t
}

func (p *Page) Image(path string) *Image {
	i := &Image{
		Path: path,
	}

	p.Values = append(p.Values, &Value{
		T: IMAGE,
		V: i,
	})

	return i
}
