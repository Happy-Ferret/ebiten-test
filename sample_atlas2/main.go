// Copyright (c) 2017 hirowaki https://github.com/hirowaki
// ebiten (https://github.com/hajimehoshi/ebiten) Copyright (c) 2013 Hajime Hoshi

package main

// https://golang.org/pkg/encoding/json/

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
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
	Frame texRect `json:"frame"`
}

type texMeta struct {
	Image string  `json:"image"`
	Size  texSize `json:"size"`
}

type jsonData struct {
	Frame map[string]*texFrame `json:"frames"`
	Meta  texMeta              `json:"meta"`

	// additional field.
	FileName  string
	FrameKeys []string
}

func (data *jsonData) ReadFile(filename string) error {
	var (
		err error
		fs  ebitenutil.ReadSeekCloser
		len int64
		bin []byte
	)

	// ebitenutil.OpenFile is great!!
	// It supports local files on JS!
	fs, err = ebitenutil.OpenFile(filename)
	defer fs.Close()
	if err != nil {
		return err
	}

	len, err = fs.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	_, err = fs.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	bin = make([]byte, len)
	_, err = fs.Read(bin)
	if err != nil {
		return err
	}

	// Unmarshal the binary data to the jsonData structure.
	json.Unmarshal(bin, data)

	data.PostProcess(filename)

	return nil
}

func (data *jsonData) PostProcess(filename string) {
	data.FileName = path.Join(path.Dir(filename), data.Meta.Image)
	data.FrameKeys = make([]string, len(data.Frame))

	index := 0
	for k, v := range data.Frame {
		data.FrameKeys[index] = k

		texRect := &v.Frame
		texRect.Rect = image.Rect(texRect.X, texRect.Y, texRect.X+texRect.W, texRect.Y+texRect.H)
		index++
	}
}

// end of schema.

const jsonFileName = "../assets/testjsonmap.json"

var (
	tick         int
	data         jsonData
	texFileName  string
	atlasTexture *ebiten.Image
)

func debugInfo(screen *ebiten.Image, targetFrame int) {
	str := fmt.Sprintf("Read and unmarshaled %s\nTexture: %s\n\n", jsonFileName, texFileName)
	for _, k := range data.FrameKeys {
		str += fmt.Sprintf(" name: %s\n rect(%s)\n\n", k, data.Frame[k].Frame.Rect.String())
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
	key := data.FrameKeys[targetFrame]
	op.SourceRect = &(data.Frame[key].Frame.Rect)
	op.GeoM.Translate(float64((WIDTH-(w>>2))>>1), y+float64(h))
	screen.DrawImage(atlasTexture, op)

	tick++
	return nil
}

func main() {
	var err error

	tick = 0

	if err = data.ReadFile(jsonFileName); err != nil {
		panic(err)
	}

	atlasTexture, _, err = ebitenutil.NewImageFromFile(data.FileName, ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}

	if err := ebiten.Run(update, WIDTH, HEIGHT, SCALE, "Image"); err != nil {
		panic(err)
	}
}
