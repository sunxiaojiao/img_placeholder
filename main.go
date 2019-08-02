package main

import (
    "net/http"
    "fmt"
    "strings"
    "log"
    "image/color"
    // "os"
    // "io"
)

const (
    // 监听端口
    Port = 8001
)


func main() {
    http.Handle("/", &CurrentHandler{})
    log.Println("server creating...")
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))
}

type CurrentHandler struct {
}

func (h * CurrentHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    path := req.URL.Path
    log.Println(path)
    if path == "" || path == "/" {
        http.ServeFile(w, req, "index.html")
        return
    }
    if path == "/favicon.ico" {
        http.ServeFile(w, req, "favicon.ico")
        return 
    }

    err := req.ParseForm()
    if err != nil {
        log.Println("request parseform fail");
    }

    params := []string{"", "200", "100"}
    pathParams := strings.Split(path, "/")

    if len(params) > len(pathParams) {
        for i := 0; i < len(pathParams); i++ {
            params[i] = pathParams[i]
        }
    } else {
        params = pathParams
    }
    log.Println(params)
    
    width := ConvInt(params[1], 200)
    height := ConvInt(params[2], 100)
    bg := req.Form.Get("bg")
    text := req.Form.Get("text")
    font := req.Form.Get("font-family")
    fontColor := req.Form.Get("font-color")
    
    if text == "" {
        text = fmt.Sprintf("%dx%d", width, height)
    }
    
    var bgColor, fColor *color.RGBA
    if bg == "" {
        bgColor = DefaultBg
    } else if strings.Contains(bg, ",") {
        bgColor, _ = ParseRGBA(bg)
    } else {
        bgColor = HexToRGBA(bg)
    }
    if fontColor == "" {
        fColor = DefaultFontColor
    } else if strings.Contains(fontColor, ",") {
        fColor, _ = ParseRGBA(fontColor)
    } else {
        fColor = HexToRGBA(fontColor)
    }
    
    img := NewImg(width, height, bgColor, text, font, fColor)
    
    // 生成图片
    imgName, imgPath, _ := img.Gen()
    imgPath = fmt.Sprintf("%s/%s", ".", imgPath)

    imgFileAddr := fmt.Sprintf("%s/%s", imgPath, imgName)
    log.Println(imgFileAddr)
    http.ServeFile(w, req, imgFileAddr)
}