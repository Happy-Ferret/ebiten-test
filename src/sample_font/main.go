// ebiten official sample is here.
// https://github.com/hajimehoshi/ebiten/blob/master/examples/font/main.go

// This sample is for a system that each TTFFont has its texture (individual canvas).
// We'll need smarter system to manage each letter.

package main

// https://github.com/golang/freetype
// https://godoc.org/github.com/golang/freetype/truetype
import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	WIDTH  = 640
	HEIGHT = 480
	SCALE  = 1
)

type FaceInfo struct {
	face  font.Face
	sw    int // source width
	sh    int // source height
	bline int // baseline
}

type TTFManager struct {
	ttf *FaceInfo
}

type TTFText struct {
	man    *TTFManager   // reference to textureManager to ask.
	canvas *ebiten.Image // canvas(texture) to write.
}

func (man *TTFManager) Setup(path string, size int) error {
	f, err := ebitenutil.OpenFile(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	// read.
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// parse.
	tt, err := truetype.Parse(b)
	if err != nil {
		return err
	}

	man.ttf = &FaceInfo{}
	man.ttf.face = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	met := man.ttf.face.Metrics()

	man.ttf.sw = size
	man.ttf.sh = met.Ascent.Ceil() + met.Descent.Ceil()
	man.ttf.bline = met.Ascent.Ceil()

	return nil
}

func (man *TTFManager) CreateText(path string) *TTFText {
	text := &TTFText{man, nil}

	return text.SetText(path)
}

func (text *TTFText) SetText(s string) *TTFText {
	// 1. calc boundary.
	// 2. render the s on image.
	// 3. render the image on the canvas.

	if text.canvas != nil {
		text.canvas.Dispose()
		text.canvas = nil
	}

	if len(s) == 0 {
		return text
	}

	para := strings.Split(s, "\n")
	var bounds fixed.Rectangle26_6
	for _, p := range para {
		b, _ := font.BoundString(text.man.ttf.face, p)
		bounds = bounds.Union(b)
	}
	w := bounds.Max.X.Ceil()
	h := text.man.ttf.sh * len(para)

	canvas, err := ebiten.NewImage(w, h, ebiten.FilterNearest)
	if err != nil {
		return nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: text.man.ttf.face,
	}
	for i, p := range para {
		d.Dot = fixed.P(0, i*text.man.ttf.sh+text.man.ttf.bline)
		d.DrawString(p)
	}

	canvas.ReplacePixels(dst.Pix)
	text.canvas = canvas

	return text
}

func (text *TTFText) Size() (width, height int) {
	if text.canvas == nil {
		return 0, 0
	}
	return text.canvas.Size()
}

var (
	ttfManager *TTFManager
	mainBody   *TTFText
	pressNext  *TTFText
	statsText  *TTFText
	index      int
	tick       int
)

var instruction = "Press [0] for English. [1] for Japanese. [2] for French."

var constitution = []string{
	"ACT I\n\n" +
		"SCENE I. A desert place.\n" +
		" Thunder and Lightning. Enter three witches\n\n" +
		"FIRST WITCH\n" +
		"  When shall we three meet again?\n" +
		"  In thunder, lightning, or in rain?\n\n" +
		"SECOND WITCH\n" +
		"  When the hurly-burly's done,\n" +
		"  When the battle's lost and won.\n\n" +
		"THIRD WITCH\n" +
		"  That will be ere the set of sun.",

	"春はあけぼの。\n" +
		"　やうやう白くなり行く山ぎは、少しあかりて、\n" +
		"　紫だちたる雲の細くたなびきたる。\n" +
		"夏は夜。\n" +
		"　月のころはさらなり、やみもなほ、ほたるの多く飛びちがひたる。\n" +
		"　また、ただ一つ、二つなど、ほのかにうち光るて行くもをかし。\n" +
		"　雨など降るもをかし。\n" +
		"秋は夕暮れ。\n" +
		"　夕日のさして山の端いと近うなりたるに、からすの寝どころへ行くとて、\n" +
		"　三つ四つ、二つ三つなど飛び急ぐさへあはれなり。\n" +
		"　まいてかりなどの連ねたるが、いと小さく見ゆるはいとをかし。\n" +
		"　日入り果てて、風の音、虫の音など、はた言ふべきにあらず。\n" +
		"冬はつとめて。\n" +
		"　雪の降りたるは言ふべきにもあらず、霜のいと白きも、またさらでも、\n" +
		"　いと寒きに、火など急ぎおこして、炭持て渡るもいとつきづきし。\n" +
		"　昼になりて、ぬるくゆるびもていけば、\n" +
		"　火桶の火も白き灰がちになりてわろし。",
	"M. Myriel.\n\n" +
		"En 1815, M. Charles-François-Bienvenu Myriel était évêque de Digne.\n" +
		"C’était un vieillard d’environ soixante-quinze ans ; il occupait le\n" +
		"siège de Digne depuis 1806.",
}

func SwitchText(i int) {
	if index != i {
		index = i
		mainBody.SetText(constitution[index])
	}

}

func update(screen *ebiten.Image) error {
	tick++

	if ebiten.IsKeyPressed(ebiten.Key0) {
		SwitchText(0)
	} else if ebiten.IsKeyPressed(ebiten.Key1) {
		SwitchText(1)
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		SwitchText(2)
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	// print at left-top.
	screen.DrawImage(mainBody.canvas, &ebiten.DrawImageOptions{})

	op := &ebiten.DrawImageOptions{}
	{ // print at bottom-right.
		w, h := pressNext.Size()
		op.GeoM.Translate(float64(WIDTH-w), float64(HEIGHT-h))
		theta := math.Pi * float64(tick%(ebiten.FPS*2)) / (ebiten.FPS * 2)
		sin := math.Min(1.0, math.Sin(theta)*1.5)
		op.ColorM.Scale(0.0, 1.0, 1.0, sin)
		screen.DrawImage(pressNext.canvas, op)
	}

	op.GeoM.Reset()
	op.ColorM.Reset()
	{ // print at top-right.
		sw, sh := mainBody.Size()
		statsText.SetText(fmt.Sprintf("canvas allocate %dx%d", sw, sh))

		w, _ := statsText.Size()
		op.GeoM.Translate(float64(WIDTH-w), 0)
		op.ColorM.Scale(0.6, 0.6, 0.6, 1.0)
		screen.DrawImage(statsText.canvas, op)
	}
	return nil
}

func main() {
	index = 0
	tick = 0

	ttfManager = &TTFManager{}
	if err := ttfManager.Setup("../assets/mplus-1p-regular.ttf", 18); err != nil {
		log.Fatal(err)
	}

	mainBody = ttfManager.CreateText(constitution[index])
	pressNext = ttfManager.CreateText(instruction)
	statsText = ttfManager.CreateText("")
	if err := ebiten.Run(update, WIDTH, HEIGHT, SCALE, "Font"); err != nil {
		log.Fatal(err)
	}
}
