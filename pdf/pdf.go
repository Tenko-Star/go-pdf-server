package pdf

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/golog"
	"github.com/signintech/gopdf"
	"io/ioutil"
	"os"
	"pdf-server/config"
)

const (
	GRID  = 0
	TEXT  = 1
	IMAGE = 2
)

var _logger *golog.Logger
var _config *config.Config

func New(logger *golog.Logger) *Pdf {
	if _config == nil {
		_config = config.Get()
	}

	if _logger == nil {
		_logger = logger
	}

	return &Pdf{
		driver: &gopdf.GoPdf{},
	}
}

func (p *Pdf) SetOutput(target string) {
	_logger.Debugf("Set new pdf struct with param file: %s", target)
	p.target = target
}

func (p *Pdf) FromJson(data string) error {
	_logger.Debugf("Json data: %s", data)
	c := &Content{}
	err := json.Unmarshal([]byte(data), c)
	if err != nil {
		return err
	}

	p.content = c

	return nil
}

func (p *Pdf) FromData(data *Content) {
	_logger.Debugf("Raw data")
	p.content = data
}

func (p *Pdf) FromTemplate(template string, texts replaceMap, insert insertMap) error {
	_logger.Debugf("Template from %s\n\t\treplace map: \t%v\n\t\tinsert  map: \t%v", template, texts, insert)

	var (
		err      error
		fileInfo os.FileInfo
		tempFile string

		tRoot = _config.Roots.Template
		mRoot string

		jsonContent []byte
		content     = &Content{}
		modules     map[string][]byte

		rMap map[int]replaceMap
	)

	tempFile = fmt.Sprintf("%s/%s", tRoot, template)
	fileInfo, err = os.Stat(tempFile)
	if err != nil {
		_logger.Error("Could not find template.")
		return err
	}

	if fileInfo.IsDir() {
		p.root = tempFile
		mRoot = tempFile
		tempFile = fmt.Sprintf("%s/%s/template.json", tRoot, template)
		_, err = os.Stat(tempFile)
		if err != nil {
			_logger.Error("Could not find template.")
			return err
		}
	} else {
		p.root = _config.Roots.Template
		mRoot = ""
	}

	file, err := os.Open(tempFile)
	defer file.Close()

	jsonContent, err = ioutil.ReadAll(file)
	if err != nil {
		_logger.Error("Could not find template.")
		return err
	}
	err = json.Unmarshal(jsonContent, content)
	if err != nil {
		_logger.Error("Could not parse template.")
		return err
	}

	modules = GetModules(content.Modules, mRoot)

	rMap = GetReplaceMaps(texts)
	_logger.Debugf("rMap: %v", rMap)
	// 记录不正确key的值
	if unknown, ok := rMap[-1]; ok {
		for key, value := range unknown {
			_logger.Warnf("Unknown replace key: %s, value: %s.", key, value)
		}
	}

	// 替换全局的部分
	for _, page := range content.Pages {
		ReplacePageValues(page, rMap[0])
	}

	// 插入新页面
	for moduleName, pages := range insert {
		_logger.Debugf("pages: %v", pages)
		var (
			ok     bool
			module []byte
			pg     *Page
		)

		if module, ok = modules[moduleName]; !ok {
			_logger.Warnf("Could not found this module was named %s.", moduleName)
			continue
		}

		for _, page := range pages {
			pg = &Page{}
			_logger.Debugf("rMap[%v]: %v", page, rMap[page])
			err = json.Unmarshal(ReplaceTokenByte(module, rMap[page]), pg)
			if err != nil {
				_logger.Warnf("Could not parse this module which named %s and will be inserted into page %d. Because %s", moduleName, page, err.Error())
				continue
			}

			content.Pages = InsertIntoPages(content.Pages, pg, page)
			_logger.Debugf("pages: %v", content.Pages)
			pg = nil
		}
	}

	// 移除空页面
	content.Pages = CompressPages(content.Pages)

	p.content = content

	return nil
}

func (p *Pdf) Save() error {
	var (
		err      error
		fileInfo os.FileInfo

		pages = p.content.Pages
	)

	p.driver.Start(gopdf.Config{
		PageSize: gopdf.Rect{
			W: p.content.Width,
			H: p.content.Height,
		},
	})

	// 处理字体文件
	for family, path := range _config.FontMap {
		err = p.driver.AddTTFFont(family, path)
		if err != nil {
			_logger.Warnf("Could not add font, because %s", err.Error())
		}
	}

	// 处理页面
	for _, page := range pages {
		p.driver.AddPage()

		// 处理页面内容
		for _, value := range page.Values {
			switch value.T {
			case TEXT:
				t := value.V.(*Text)

				// 设置坐标
				x, y := t.GetPosition(p.content.Width)
				p.driver.SetX(x)
				p.driver.SetY(y)

				// 设置字体
				err = p.driver.SetFont(t.GetFontConfig())
				if err != nil {
					_logger.Errorf("Could not set font, because %s", err.Error())
					return err
				}

				// 设置颜色
				p.driver.SetTextColor(t.GetColor())

				// 打印字符
				err = p.driver.Cell(nil, t.Text)
				if err != nil {
					_logger.Errorf("Could not print, because %s", err.Error())
					return err
				}

				p.resetText()
				// break

			case IMAGE:
				i := value.V.(*Image)

				// 获取图片文件
				image, err := i.GetFile(p.root)
				if err != nil {
					_logger.Warnf("Could not load photo, because %s", err.Error())
					break
				}
				p.temp = append(p.temp, image)

				// 设置坐标
				x, y := i.GetPosition()
				p.driver.SetX(x)
				p.driver.SetY(y)
				// 设置图片
				err = p.driver.Image(image.path, x, y, i.GetSize())
				if err != nil {
					_logger.Warnf("Could not print photo, because %s", err.Error())
				}

				p.resetPos()
			}
		}
	}

	fileInfo, err = os.Stat(_config.Output)
	if err != nil {
		err = os.MkdirAll(_config.Output, os.ModePerm)
	} else if !fileInfo.IsDir() {
		_logger.Errorf("Could write pdf to target, because output is not a dir")
		return err
	}

	if err != nil {
		_logger.Errorf("Could write pdf to target, because %s", err.Error())
		return err
	}

	fileInfo, err = os.Stat(p.target)
	if err == nil && !fileInfo.IsDir() {
		_ = os.Remove(p.target)
	}

	err = p.driver.WritePdf(p.target)
	if err != nil {
		_logger.Warnf("Could write pdf to target, because %s", err.Error())
		return err
	}

	// 关闭临时文件句柄
	p.closeDriver()
	//for _, file := range p.temp {
	//	file.Close()
	//}

	return nil
}

func (p *Pdf) closeDriver() {
	err := p.driver.Close()
	if err != nil {
		_logger.Warnf("Could not to close driver, because %s", err.Error())
	}
}

func (p *Pdf) resetText() {
	p.resetPos()
	p.driver.SetTextColor(255, 255, 255)
}

func (p *Pdf) resetPos() {
	p.driver.SetX(0)
	p.driver.SetY(0)
}
