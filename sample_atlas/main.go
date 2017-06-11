package main

// https://golang.org/pkg/encoding/json/

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"path"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	WIDTH  = 640
	HEIGHT = 480
	SCALE  = 1
)

// schema.
type texRect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`

	// additional field.
	Rect image.Rectangle
}

type texSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type texFrame struct {
	Filename string  `json:"filename"`
	Frame    texRect `json:"frame"`
}

type texMeta struct {
	Image string  `json:"image"`
	Size  texSize `json:"size"`
}

type jsonData struct {
	Frame []texFrame `json:"frames"`
	Meta  texMeta    `json:"meta"`
}

func (data *jsonData) PostProcess() {
	for i := 0; i < len(data.Frame); i++ {
		texRect := &data.Frame[i].Frame
		texRect.Rect = image.Rect(texRect.X, texRect.Y, texRect.X+texRect.W, texRect.Y+texRect.H)
	}
}

// end of schema.

const jsonFileName = "../assets/testjson.json"

var (
	tick         int
	data         jsonData
	texFileName  string
	atlasTexture *ebiten.Image
)

func debugInfo(screen *ebiten.Image, targetFrame int) {
	str := fmt.Sprintf("Read and unmarshaled %s\nTexture: %s\n\n", jsonFileName, texFileName)
	for _, frame := range data.Frame {
		str += fmt.Sprintf(" name: %s\n rect(%s)\n\n", frame.Filename, frame.Frame.Rect.String())
	}
	str += fmt.Sprintf("\n Target Frame: %d", targetFrame)

	ebitenutil.DebugPrint(screen, str)
}

func update(screen *ebiten.Image) error {
	// Fill the screen with #FF0000 color
	screen.Fill(color.NRGBA{0x40, 0x40, 0x40, 0xff})

	targetFrame := (tick % 60) / 15

	debugInfo(screen, targetFrame)

	// the whole texture.
	op := &ebiten.DrawImageOptions{}

	w, h := atlasTexture.Size()
	x := float64((WIDTH - w) >> 1)
	y := float64((HEIGHT >> 1) - h)
	op.GeoM.Translate(x, y)
	screen.DrawImage(atlasTexture, op)

	op.GeoM.Reset()

	// only the texture atlas.
	op.SourceRect = &(data.Frame[targetFrame].Frame.Rect)
	op.GeoM.Translate(float64((WIDTH-(w>>2))>>1), y+float64(h))
	screen.DrawImage(atlasTexture, op)

	tick++
	return nil
}

func main() {
	var err error
	var jsonFile []byte

	tick = 0

	jsonFile, err = ioutil.ReadFile("../assets/testjson.json")
	if err != nil {
		panic(err)
	}

	// Unmarshal the binary data to the jsonData structure.
	json.Unmarshal(jsonFile, &data)
	data.PostProcess()

	texFileName = path.Join(path.Dir(jsonFileName), data.Meta.Image)
	atlasTexture, _, err = ebitenutil.NewImageFromFile(texFileName, ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}

	if err := ebiten.Run(update, WIDTH, HEIGHT, SCALE, "Image"); err != nil {
		panic(err)
	}
}
