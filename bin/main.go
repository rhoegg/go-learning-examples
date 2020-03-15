package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
)

func main() {
	log.Print("Hello learner")
	if err := ebiten.Run(update, 640, 480, 1, "Hello, world!"); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.White)
	gameScene, _ := ebiten.NewImage(640, 480, ebiten.FilterDefault)
	gameScene.Fill(color.Black)
	cols, rows := 3, 3

	lineWidth := 3
	lineColor := color.RGBA{0xd0, 0, 0, 0xa0} // RED
	w, h := gameScene.Size()
	for c := 1; c < cols; c++ {
		center := c * w / cols
		ebitenutil.DrawRect(gameScene, float64(center-(lineWidth/2)), 0, float64(lineWidth), float64(h), lineColor)
	}
	for r := 1; r < rows; r++ {
		center := r * h / rows
		ebitenutil.DrawRect(gameScene, 0, float64(center-(lineWidth/2)), float64(w), float64(lineWidth), lineColor)
	}

	reader, err := os.Open("assets/tic-tac-toe-game.png")
	if err != nil {
		log.Printf("Error openiing: %v", err)
	}
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear

	img, _, err := image.Decode(reader)
	spriteMap, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	tileW, tileH := spriteMap.Size()

	vSlices := []int{121, 161, 297, 336}
	hSlices := []int{119, 159, 294, 334}

	board, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	hTile1, _ := ebiten.NewImageFromImage(spriteMap.SubImage(image.Rect(vSlices[0], 0, vSlices[1], h)), ebiten.FilterDefault)
	hTile2, _ := ebiten.NewImageFromImage(spriteMap.SubImage(image.Rect(vSlices[2], 0, vSlices[3], h)), ebiten.FilterDefault)
	vTile1, _ := ebiten.NewImageFromImage(spriteMap.SubImage(image.Rect(0, hSlices[0], w, hSlices[1])), ebiten.FilterDefault)
	vTile2, _ := ebiten.NewImageFromImage(spriteMap.SubImage(image.Rect(0, hSlices[2], w, hSlices[3])), ebiten.FilterDefault)

	xTile, _ := ebiten.NewImageFromImage(spriteMap.SubImage(image.Rect(0, 0, vSlices[0], hSlices[0])), ebiten.FilterDefault)
	oTile, _ := ebiten.NewImageFromImage(spriteMap.SubImage(image.Rect(vSlices[1], 0, vSlices[2], hSlices[0])), ebiten.FilterDefault)

	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(float64(hSlices[0]+2), 0)
	board.DrawImage(hTile1, op)
	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(float64(hSlices[2]+3), 0)
	board.DrawImage(hTile2, op)
	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(0, float64(vSlices[0])-2)
	board.DrawImage(vTile1, op)
	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(0, float64(vSlices[2])-3)
	board.DrawImage(vTile2, op)
	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(0, 0)
	board.DrawImage(xTile, op)
	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(float64(hSlices[1]+2), float64(vSlices[1]+7))
	board.DrawImage(oTile, op)

	backgroundColor := color.RGBA{0x60, 0x70, 0x80, 0xc0}
	gameScene.Fill(backgroundColor)

	//for _, x := range vSlices {
	//	ebitenutil.DrawLine(board, float64(x), 0, float64(x), float64(tileH), lineColor)
	//}
	//
	//for _, y := range hSlices {
	//	ebitenutil.DrawLine(board, 0, float64(y), float64(tileW), float64(y), lineColor)
	//}
	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(float64((w-tileW)/2), float64((h-tileH)/2))
	//gameScene.DrawImage(spriteMap, op)
	gameScene.DrawImage(board, op)

	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(gameScene, op)
	return nil
}
