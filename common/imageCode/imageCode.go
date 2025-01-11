package imageCode

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math/rand"
	"time"
)

type CreateImageCode struct {
	width     int
	height    int
	codeCount int
	lineCount int
	Code      string
	buffImg   *image.RGBA
	random    *rand.Rand
}

func NewCreateImageCode() *CreateImageCode {
	c := &CreateImageCode{
		width:     130,
		height:    38,
		codeCount: 4,
		lineCount: 5,
		random:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	c.createImage()
	return c
}

func NewCreateImageCodeWithSize(width, height int) *CreateImageCode {
	c := &CreateImageCode{
		width:     width,
		height:    height,
		codeCount: 4,
		lineCount: 5,
		random:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	c.createImage()
	return c
}

func (c *CreateImageCode) createImage() {
	c.buffImg = image.NewRGBA(image.Rect(0, 0, c.width, c.height))

	// 填充背景色
	bgColor := c.getRandColor(200, 250)
	draw.Draw(c.buffImg, c.buffImg.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Over)

	// 绘制干扰线
	for i := 0; i < c.lineCount; i++ {
		x1 := c.random.Intn(c.width)
		y1 := c.random.Intn(c.height)
		x2 := x1 + c.random.Intn(c.width/2)
		y2 := y1 + c.random.Intn(c.height/2)
		c.drawLine(x1, y1, x2, y2, c.getRandColor(1, 255))
	}

	// 添加噪点
	yawpRate := 0.01
	area := int(yawpRate * float64(c.width*c.height))
	for i := 0; i < area; i++ {
		x := c.random.Intn(c.width)
		y := c.random.Intn(c.height)
		c.buffImg.Set(x, y, color.RGBA{
			R: uint8(c.random.Intn(255)),
			G: uint8(c.random.Intn(255)),
			B: uint8(c.random.Intn(255)),
			A: 255,
		})
	}

	// 生成验证码文字
	c.Code = c.randomStr(c.codeCount)
	fontWidth := c.width / c.codeCount

	for i := 0; i < c.codeCount; i++ {
		x := i*fontWidth + 3
		y := c.height - 8
		c.drawString(string(c.Code[i]), x, y, c.getRandColor(1, 255))
	}
}

func (c *CreateImageCode) drawLine(x1, y1, x2, y2 int, col color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	steep := dy > dx

	if steep {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
	}
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	dx = x2 - x1
	dy = abs(y2 - y1)
	err := dx / 2
	ystep := -1
	if y1 < y2 {
		ystep = 1
	}

	for ; x1 <= x2; x1++ {
		if steep {
			c.buffImg.Set(y1, x1, col)
		} else {
			c.buffImg.Set(x1, y1, col)
		}
		err -= dy
		if err < 0 {
			y1 += ystep
			err += dx
		}
	}
}

func (c *CreateImageCode) drawString(s string, x, y int, col color.Color) {
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  c.buffImg,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(s)
}

func (c *CreateImageCode) randomStr(n int) string {
	str := "ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz123456789"
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = str[c.random.Intn(len(str))]
	}
	return string(bytes)
}

func (c *CreateImageCode) getRandColor(fc, bc int) color.Color {
	if fc > 255 {
		fc = 255
	}
	if bc > 255 {
		bc = 255
	}
	r := fc + c.random.Intn(bc-fc)
	g := fc + c.random.Intn(bc-fc)
	b := fc + c.random.Intn(bc-fc)
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

func (c *CreateImageCode) Write(w io.Writer) error {
	return png.Encode(w, c.buffImg)
}

func (c *CreateImageCode) GetCode() string {
	return c.Code
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
