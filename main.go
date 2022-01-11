package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	basePath  string
	dataPath  string
	batchSize int
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	basePath = filepath.Dir(ex)
	dataPath = fmt.Sprintf("%s/%s", basePath, "/data.json")
	batchSize = 50
}

func main() {
	images := parseJson(dataPath)

	l := len(images)

	chunks := chunk(images, batchSize)
	pages := len(chunks)

	for i, chunkList := range chunks {
		wg := sync.WaitGroup{}
		wg.Add(len(chunkList))

		for idx, item := range chunkList {
			go func(item ImageItem, idx int) {
				processRow(item)
				log("item done %d/%d", 1+idx+i*batchSize, l)
				// log("结束 %d/%d", idx+1, l)
				wg.Done()
			}(item, idx)
		}
		wg.Wait()
		log("page done %d/%d", i+1, pages)
	}

	log("parse done")
	// os.Exit(0)
}

func log(f string, a ...interface{}) {
	str := fmt.Sprintf(f, a...)
	fmt.Println(str)
}

type JsonData struct {
	Images []ImageItem `json:"images"`
}

type ImageItem struct {
	Path string `json:"path"` // 路径
	Url  string `json:"url"` // 图片
	Name string `json:"name"` // 图片名
}

func parseJson(path string) []ImageItem {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic("读取文件错误:" + err.Error())
	}

	var data []ImageItem

	err = json.Unmarshal(content, &data)
	if err != nil {
		panic("文件内容错误:" + err.Error())
	}

	return data
}

func makeDir(name string) bool {
	return os.MkdirAll(name, os.ModePerm) == nil
}

func saveImage(path, name string, url ...string) bool {
	action := func(path, name, url string, idx int) bool {
		if url == "" {
			log("%s url is empty ", name)
			return false
		}

		client := http.Client{
			Timeout: 10 * time.Minute,
		}
		res, err := client.Get(url)

		if res != nil {
			defer res.Body.Close()
		}

		if err != nil {
			log("client get error " + err.Error())
			return false
		}

		var fPath string
		if idx == 0 {
			fPath = fmt.Sprintf("%s/%s", path, name)
		} else {
			fPath = fmt.Sprintf("%s/%s%d", path, name, idx+1)
		}

		f, err := os.Create(fPath)

		if err != nil {
			log("create file error " + err.Error())
			return false
		}

		_, err = io.Copy(f, res.Body)

		if err != nil {
			log("copy error:%+v, url: %+v ", err.Error(), url)
			return true
		}

		// 加后缀
		ext := getImgType(f)
		newFPath := fmt.Sprintf("%s.%s", fPath, ext)

		// ！！！此处要先关闭文件才能重命名！！！
		f.Close()

		err = os.Rename(fPath, newFPath)
		if err != nil {
			log("rename error " + err.Error())
			return false
		}

		return false
	}

	var retry bool
	for idx, v := range url {
		retry = action(path, name, v, idx) || retry
	}

	return true
}

func processRow(i ImageItem) {
	makeDir(i.Path)
	saveImage(i.Path, i.Name, i.Url)
}

func getImgType(file *os.File) (ext string) {
	ext = "jpg"
	buff := make([]byte, 512)

	_, err := file.Read(buff)

	if err != nil {
		return
	}

	filetype := http.DetectContentType(buff)

	exts := []string{
		"ase",
		"art",
		"bmp",
		"blp",
		"cd5",
		"cit",
		"cpt",
		"cr2",
		"cut",
		"dds",
		"dib",
		"djvu",
		"egt",
		"exif",
		"gif",
		"gpl",
		"grf",
		"icns",
		"ico",
		"iff",
		"jng",
		"jpeg",
		"jpg",
		"jfif",
		"jp2",
		"jps",
		"lbm",
		"max",
		"miff",
		"mng",
		"msp",
		"nitf",
		"ota",
		"pbm",
		"pc1",
		"pc2",
		"pc3",
		"pcf",
		"pcx",
		"pdn",
		"pgm",
		"PI1",
		"PI2",
		"PI3",
		"pict",
		"pct",
		"pnm",
		"pns",
		"ppm",
		"psb",
		"psd",
		"pdd",
		"psp",
		"px",
		"pxm",
		"pxr",
		"qfx",
		"raw",
		"rle",
		"sct",
		"sgi",
		"rgb",
		"int",
		"bw",
		"tga",
		"tiff",
		"tif",
		"vtf",
		"xbm",
		"xcf",
		"xpm",
		"3dv",
		"amf",
		"ai",
		"awg",
		"cgm",
		"cdr",
		"cmx",
		"dxf",
		"e2d",
		"egt",
		"eps",
		"fs",
		"gbr",
		"odg",
		"svg",
		"stl",
		"vrml",
		"x3d",
		"sxd",
		"v2d",
		"vnd",
		"wmf",
		"emf",
		"art",
		"xar",
		"png",
		"webp",
		"jxr",
		"hdp",
		"wdp",
		"cur",
		"ecw",
		"iff",
		"lbm",
		"liff",
		"nrrd",
		"pam",
		"pcx",
		"pgf",
		"sgi",
		"rgb",
		"rgba",
		"bw",
		"int",
		"inta",
		"sid",
		"ras",
		"sun",
		"tga",
	}

	for i := 0; i < len(ext); i++ {
		if strings.Contains(exts[i], filetype[6:len(filetype)]) {
			ext = strings.Replace(filetype, "image/", "", -1)
			if ext == "jpeg" {
				ext = "jpg"
				return
			}
		}
	}

	return
}

func chunk(items []ImageItem, chunkSize int) (chunks [][]ImageItem) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}
