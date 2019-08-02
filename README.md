# 生成占位图

一个用go写的生成占位图的开源小服务

## 使用在线服务

[http://placeholder.quanquan88.com/](http://placeholder.quanquan88.com/)

## 安装到自己的服务器

如果你想部署一个到自己的服务器，可以按照如下步骤完成部署。

```shell
# 手动编译
git clone https://github.com/airect/img_placeholder.git && cd img_placeholder
go build *.go
sudo chmod -R 776 images
nohup ./main &
```

可以使用nginx做一层代理，方便记录访问日志，平滑部署

默认端口 8001，可以通过修改main.go中的`Port`改变监听端口

## 使用
```
http://localhost:port/{width}/{height}?text={text}&bg={bg}&font-family={font-family}&font-color={font-color}
```

| 参数          | 作用                             | 默认值           | 备注           |
| ------------- | -------------------------------- | ---------------- | -------------- |
| {width}       | 图片宽度px                       | 200              | path中         |
| {height}      | 图片高度px                       | 100              | path中         |
| {bg}          | 图片背景颜色 支持rgba和6位16进制 | cccccc           | query_string中 |
| {font-color}  | 文字颜色 支持rgba和6位16进制     | 666666           | query_string中 |
| {text}        | 图片中的文本，居中展示           | {width}x{height} | query_string中 |
| {font-family} | 字体，见附录                     | 000 思源黑体     | query_string中 |

  

## 示例

### width 300px height 50px
```html
<img src="http://placeholder.quanquan88.com/300/50">
```
![](http://placeholder.quanquan88.com/300/50)

### width 300px height 100px 红色背景，白色文字

```html
<img src="http://placeholder.quanquan88.com/300/100?bg=255,0,0&font-color=ffffff">
```
![](http://placeholder.quanquan88.com/300/100?bg=255,0,0&font-color=ffffff)

### width 300px height 100px 红色背景，白色文字，带文字

```html
<img src="http://placeholder.quanquan88.com/300/100?bg=255,0,0&font-color=ffffff&text=中文">
```
![](http://placeholder.quanquan88.com/300/100?bg=255,0,0&font-color=ffffff&text=中文)

## 附录

### 字体支持

> 仅支持ttf字体，下面是挑选了几个免费字体
> 可以通过自己部署程序，添加字体到fonts文件夹，并修改placeholder.go中的`fontMap`来拓展字体

| 字体            | 编号     | 示例                                                         |
| --------------- | -------- | ------------------------------------------------------------ |
| 思源黑体        | 000 默认 | ![](http://placeholder.quanquan88.com/100/100?bg=255,255,0&font-family=000&text=GO) |
| 阿里巴巴普惠体M | 001      | ![](http://placeholder.quanquan88.com/100/100?font-family=001&text=中文) |
| 方正黑体简体    | 002      | ![](http://placeholder.quanquan88.com/100/100?bg=0,255,255&font-family=002&text=中文&font-color=99999) |
| 文泉驿等宽正黑  | 003      | ![](http://placeholder.quanquan88.com/100/100?bg=0,255,255&font-family=003&text=中文) |
| 黄引齐招牌体    | 004      | ![](http://placeholder.quanquan88.com/100/100?bg=0,255,255&font-family=004&text=中文) |

