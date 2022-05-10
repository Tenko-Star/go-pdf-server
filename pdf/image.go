package pdf

import (
	"errors"
	"fmt"
	"github.com/signintech/gopdf"
	"os"
	"pdf-server/http"
	"regexp"
)

var regexHttp = regexp.MustCompile(`https?://.*\.(png|jpg|jpeg|bmp)`)

type Image struct {
	Path    string  `json:"path"`
	OffsetX float64 `json:"offset_x"`
	OffsetY float64 `json:"offset_y"`
	Rect    *gopdf.Rect
	Props   map[name]value `json:"props"`
}

type ImageFile struct {
	path   string
	isTemp bool
}

func (i *Image) Set(x, y float64) {
	i.OffsetX = x
	i.OffsetY = y
}

func (i *Image) SetX(x float64) {
	i.OffsetX = x
}

func (i *Image) SetY(y float64) {
	i.OffsetY = y
}

func (i *Image) GetSize() *gopdf.Rect {
	propHeight(i)
	propWidth(i)

	return i.Rect
}

func (i *Image) GetPosition() (float64, float64) {
	return i.OffsetX, i.OffsetY
}

func (i *Image) GetFile(tRoot string) (*ImageFile, error) {
	var (
		err  error
		info os.FileInfo

		ret = &ImageFile{}
	)

	if regexHttp.MatchString(i.Path) {
		_logger.Debugf("Start to download from %s", i.Path)
		ret.path, err = http.DownloadToTemp(i.Path, _config.Roots.Temporary)
		if err != nil {
			return nil, err
		}

		ret.isTemp = true
	} else {
		path := fmt.Sprintf("%s/%s", tRoot, i.Path)
		info, err = os.Stat(path)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			return nil, errors.New("is not a correct file")
		}

		ret.path = path

		ret.isTemp = false
	}

	return ret, nil
}

//func (i *ImageFile) Close() {
//	var (
//		err error
//
//		name = i.handle.Name()
//	)
//
//	if i.isTemp {
//		err = i.handle.Close()
//		if err != nil {
//			_logger.Warnf("Could not to close file which named %s, because %s", i.handle.Name(), err.Error())
//		}
//
//		err = os.Remove(i.handle.Name())
//		if err != nil {
//			_logger.Warnf("Could not to delete file which named %s, because %s", i.handle.Name(), err.Error())
//		}
//	} else {
//		err = i.handle.Close()
//		if err != nil {
//			_logger.Warnf("Could not to close file which named %s, because %s", i.handle.Name(), err.Error())
//		}
//	}
//
//	if i.isTemp {
//		_logger.Debugf("Success close and delete temp file: %s", name)
//	} else {
//		_logger.Debugf("Success close file: %s", name)
//	}
//}
