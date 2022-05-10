package pdf

import (
	"encoding/json"
	"errors"
	"github.com/signintech/gopdf"
)

type Pdf struct {
	target  string
	content *Content
	driver  *gopdf.GoPdf
	temp    []*ImageFile

	root string
}

type Value struct {
	T int     `json:"type"`
	V Movable `json:"value"`
}

type Movable interface {
	Set(x, y float64)
	SetX(x float64)
	SetY(y float64)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	t := struct {
		T int `json:"type"`
	}{}

	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	switch t.T {
	case IMAGE:
		i := struct {
			T int   `json:"type"`
			I Image `json:"value"`
		}{}
		err := json.Unmarshal(data, &i)
		if err != nil {
			return err
		}

		v.T = IMAGE
		v.V = &i.I
		return nil
	case TEXT:
		t := struct {
			T int  `json:"type"`
			I Text `json:"value"`
		}{}
		err := json.Unmarshal(data, &t)
		if err != nil {
			return err
		}

		v.T = TEXT
		v.V = &t.I
		return nil
	}

	return errors.New("unknown type")
}
