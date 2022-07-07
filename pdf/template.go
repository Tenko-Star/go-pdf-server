package pdf

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type replaceMap = map[string]string
type insertMap = map[string][]int
type moduleMap = map[string]string

func ReplaceTokenString(data string, m replaceMap) string {
	if m == nil {
		return data
	}

	for key, value := range m {
		key = "{" + key + "}"
		value = strings.ReplaceAll(value, "\r", "")
		value = strings.ReplaceAll(value, "\n", "")
		value = strings.ReplaceAll(value, "\t", "")
		data = strings.ReplaceAll(data, key, value)
	}

	return data
}

func ReplaceTokenByte(data []byte, m replaceMap) []byte {
	if m == nil {
		return data
	}

	return []byte(ReplaceTokenString(string(data), m))
}

func GetReplaceMaps(maps map[string]string) map[int]replaceMap {
	var (
		ok    bool
		key   string
		value string

		currentSection int
		currentKey     string
		section        replaceMap

		ret = make(map[int]replaceMap)
	)

	for key, value = range maps {
		currentKey, currentSection = ParseToKeyAndSection(key)
		if section, ok = ret[currentSection]; !ok {
			ret[currentSection] = make(replaceMap)
			section = ret[currentSection]
		}

		section[currentKey] = value
	}

	return ret
}

func ParseToKeyAndSection(s string) (string, int) {
	var (
		key string
		sec int
		err error
	)

	if pos := strings.LastIndexByte(s, '|'); pos == -1 {
		return s, 0
	} else {
		key = s[:pos]
		sec, err = strconv.Atoi(s[pos+1:])
		if err != nil {
			return s, -1
		}

		return key, sec
	}
}

func GetModules(mMap moduleMap, mRoot string) map[string][]byte {
	var (
		file     *os.File
		filePath string

		ret = make(map[string][]byte)

		err error
	)

	for name, path := range mMap {
		if len(path) == 0 {
			path = name + ".json"
		}

		filePath = fmt.Sprintf("%s/%s", mRoot, path)
		_, err = os.Stat(filePath)
		if err != nil {
			ret[name] = nil
			continue
		}

		file, err = os.Open(filePath)
		if err != nil {
			ret[name] = nil
			continue
		}

		ret[name], err = ioutil.ReadAll(file)
		if err != nil {
			ret[name] = nil
		}

		_ = file.Close()
	}

	return ret
}

func ReplacePageValues(page *Page, rMap replaceMap) {
	var (
		err  error
		temp string
	)

	if rMap == nil {
		return
	}

	for _, value := range page.Values {
		switch value.T {
		case TEXT:
			temp, err = replaceText(value.V.(*Text).Text, rMap)
			if err != nil {
				value.V.(*Text).Text = ""
				continue
			}

			value.V.(*Text).Text = temp
			// break
		case IMAGE:
			temp, err = replaceText(value.V.(*Image).Path, rMap)
			if err != nil {
				value.V.(*Image).Path = ""
				continue
			}

			value.V.(*Image).Path = temp
			// break
		}
	}
}

func replaceText(t string, rMap replaceMap) (string, error) {
	var (
		regex *regexp.Regexp
		err   error
	)

	regex, err = regexp.Compile(`\{(.*)}`)
	if err != nil {
		return "", err
	}
	subMatches := regex.FindAllStringSubmatch(t, -1)
	if subMatches == nil {
		return t, nil
	}

	for _, match := range subMatches {
		t = strings.Replace(t, match[0], rMap[match[1]], 1)
	}

	return t, nil
}

func InsertIntoPages(pages []*Page, page *Page, pos int) []*Page {
	var (
		temp []*Page
	)

	if pos >= len(pages) {
		// 添加到一个比现有长度还要大的位置时
		temp = make([]*Page, pos+1)
		copy(temp, pages)
		temp[pos] = page

	} else if pages[pos] == nil {
		// 正好存在空位时
		pages[pos] = page
		return pages

	} else {
		// 稠密排列或pos位正好被占用
		temp = make([]*Page, 0, len(pages)+1)
		for i := 0; i < pos; i++ {
			temp = append(temp, pages[i])
		}
		temp = append(temp, page)
		for i := pos; i < len(pages); i++ {
			temp = append(temp, pages[i])
		}

	}

	return temp
}

func CompressPages(pages []*Page) []*Page {
	var (
		length = len(pages)

		ret []*Page
	)

	for i := 0; i < length; i++ {
		if pages[i] == nil {
			continue
		}

		ret = append(ret, pages[i])
	}

	return ret
}
