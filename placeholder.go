// 生成图片
// TOOD 淘汰旧图片

package main

import (
	"strconv"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"flag"
	"fmt"
	"path/filepath"
	"os"
	"log"
	"strings"
	"io/ioutil"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"math"
	"bufio"
	"crypto/md5"
	"encoding/hex"
)

const (
	ImagePath = "images"
	fileMode = 0777
)

var (
	dpi      = flag.Float64("dpi", 172, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "", "filename of the ttf font")
	size     = flag.Float64("size", 14, "font size in points")
	spacing  = flag.Float64("spacing", 2, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
	DefaultBg = HexToRGBA("cccccc")
	DefaultFontColor = HexToRGBA("666666")
)

var fontMap = map[string]string{
	"000": "SourceHanSansCN-Medium.ttf",
	"001": "阿里巴巴普惠体M.ttf",
	"002": "方正黑体简体.ttf",
	"003": "文泉驿等宽正黑.ttf",
	"004": "黄引齐招牌体.ttf",
}

type Img struct {
	width int
	height int
	bgColor *color.RGBA
	text string
	name string
	font string
	fontColor *color.RGBA
}

func NewImg (width int, height int, bgColor *color.RGBA, text string, font string, color *color.RGBA) *Img {

	if _, ok := fontMap[font]; ok {
		font = fontMap[font]
	} else {
		font = fontMap["000"]
	}

    return &Img {
		width,
		height,
		bgColor,
		text,
		"",
		font,
		color,
    }
}

// 生成
func (img *Img) Gen () (string, string, error) {
	imgFileName := img.ImgName(false)
	log.Println(imgFileName)

	// 先判断图片是否存在
	if ok, _ := img.FileExisted(); ok {
		log.Println("图片存在")
		return imgFileName, img.ImgPath(), nil
	}

	rgba := image.NewRGBA(image.Rect(0, 0, img.width, img.height))
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{*img.bgColor}, image.ZP, draw.Src)

	// 绘制文本
	fontBytes, err := ioutil.ReadFile(fmt.Sprintf("./fonts/%s", img.font))
	if err != nil {
		log.Fatal(err)
	}
	f , err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	textCount :=  strings.Count(img.text, "") - 1
	// 计算字体大小
	textSize1 := float64(Min(img.width, img.height) / 5)
	textSize2 := float64(img.width / textCount) / *spacing

	*size = math.Min(textSize1, textSize2)
	
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(&image.Uniform{*img.fontColor})
	c.SetHinting(font.HintingNone)

	// 计算文字位置
	// https://developer.apple.com/library/archive/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyph_metrics_2x.png
	opts := truetype.Options{}
	opts.Size = *size
	opts.DPI = *dpi
	opts.Hinting = font.HintingNone

	face := truetype.NewFace(f, &opts)
	fwidth := 0
	fheight := 0
	for _, x := range(img.text) {
		bounds, awidth, _ := face.GlyphBounds(rune(x))
		fwidth = fwidth + int(awidth) >> 6
		fheight = int(bounds.Max.Y - bounds.Min.Y) >> 6
	}
	
	fX := (img.width - fwidth) / 2
	fY := (img.height - fheight) / 2 + fheight
	pt := freetype.Pt(fX, fY)

	// 绘制文字
	_, err = c.DrawString(img.text, pt)
	// 图片存储
	imgName, imgPath, err := img.Store(rgba)

	if err != nil {
		return "", "", err
	}
    return imgName, imgPath, err
}

// 图片是否存在
func (img *Img) FileExisted () (bool, error) {
	imgPath := fmt.Sprintf("%s/%s", img.ImgPath(), img.ImgName(false))

	_, err := os.Stat(imgPath)
    if err == nil {
        return true, nil	
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

// 存储
func (img *Img) Store (rgba *image.RGBA) (string, string, error) {
	imgPath := img.ImgPath()
	log.Println(imgPath)
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		errDir := os.MkdirAll(imgPath, fileMode)
		if errDir != nil {
			log.Fatal(err)
		}
	} 
	
	pathName := fmt.Sprintf("%s/%s", img.ImgPath(), img.ImgName(false))
	outFile, err := os.Create(pathName)
	if err != nil {
		return "", "", err
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		return "", "", err
	}
	err = b.Flush()
	if err != nil {
		return "", "", err
	}

	return img.ImgName(false), img.ImgPath(), nil
}

// 图片文件名
func (img *Img) ImgName (force bool) string {
	if img.name == "" || force {
		name := fmt.Sprintf("%s.png", img.Hash())
		img.name = name
	}
	return img.name
}

// 图片存储路径
func (img *Img) ImgPath () string {
	imgName := img.ImgName(false)
	
	imgPath := fmt.Sprintf("%s/%s", ImagePath, imgName[0:2])
	return imgPath
}

// 摘要 用于图片文件名
func (img *Img) Hash () string {
	str := fmt.Sprintf("%d%d%v%s%s%v", img.width, img.height, img.bgColor, img.text, img.font, img.fontColor)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func GetAbsPath () string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func ConvInt(input string, defaultInt int) (int) {
	convError := errors.New("参数转换失败")

	widthInt, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		log.Println(convError)
		return defaultInt
	}
	return int(widthInt)
}

func ParseRGBA(input string) (*color.RGBA, error) {
	arr := strings.Split(input, ",")
	if len(arr) >= 4 {
		return &color.RGBA{
			uint8(ConvInt(arr[0], 0)),
			uint8(ConvInt(arr[1], 0)),
			uint8(ConvInt(arr[2], 0)),
			uint8(ConvInt(arr[3], 0))}, nil
	} else if len(arr) == 3 {
		return &color.RGBA{
			uint8(ConvInt(arr[0], 0)),
			uint8(ConvInt(arr[1], 0)),
			uint8(ConvInt(arr[2], 0)),
			255}, nil
	} else {
		return &color.RGBA{0,0,0,0}, errors.New("无法解析rgba")
	}
}

func HexToRGBA(hex string) *color.RGBA {
	h, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		h = 0
	}
	rgba := &color.RGBA{
		R: uint8(h >> 16),
		G: uint8((h & 0x00ff00) >> 8),
		B: uint8(h & 0x0000ff),
		A: 255,
	}

	return rgba
}

func Min(x, y int) int {
    if x < y {
        return x
    }
    return y
}