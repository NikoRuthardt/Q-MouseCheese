package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 600
	screenHeight = 600
	tileSize     = 100
	padding      = 100
)

// Player

// Dir represents a direction.
type Dir int

const (
	up Dir = iota
	down
	left
	right
	none
)

type mouse struct {
	row, col int
	dir      Dir
}

func (m *mouse) move(d Dir) {
	switch d {
	case up:
		m.col--
	case down:
		m.col++
	case left:
		m.row--
	case right:
		m.row++
	}
}

func (m *mouse) input(action int) (Dir, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || action == 0 {
		if m.col == 0 {
			return none, true
		}
		return up, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || action == 1 {
		if m.col == screenHeight/tileSize-1 {
			return none, true
		}
		return down, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || action == 2 {
		if m.row == 0 {
			return none, true
		}
		return left, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || action == 3 {
		if m.row == screenWidth/tileSize-1 {
			return none, true
		}
		return right, true
	}

	return none, false
}

// reset --> end of episode else return reward
func (m *mouse) updateState(grid Grid) int {
	nextTile := grid.getIndex(m.row, m.col)
	switch grid.tiles[nextTile].value {
	case 0:
		return -1
	case 1:
		grid.tiles[nextTile].value = 0
		return 5
	case 2:
		m.reset()
		return -100
	case 3:
		m.reset()
		return 100
	}
	return 0
}

func (m *mouse) reset() {
	m.col = 0
	m.row = 0
	episode++
	initGrid()
}

// GameBoard

// Grid the game board
type Grid struct {
	tiles []Tile
}

func (g *Grid) getIndex(row, col int) int {
	for i, t := range g.tiles {
		if t.cols == col && t.rows == row {
			return i
		}
	}
	return -1
}

// Tile :each point on the grid is a tile
type Tile struct {
	rows, cols, value int
}

var (
	tileImage, cheeseImage, trapImage, trippleCheeseImage, mouseImage *ebiten.Image
	grid                                                              Grid
	player                                                            mouse
	agent                                                             *Agent
	episode                                                           int
)

// load images
func init() {
	var err error

	tileImage, _, err = ebitenutil.NewImageFromFile("assets/tile_1x.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	cheeseImage, _, err = ebitenutil.NewImageFromFile("assets/cheese_1x.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	trippleCheeseImage, _, err = ebitenutil.NewImageFromFile("assets/trippleCheese_1x.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	trapImage, _, err = ebitenutil.NewImageFromFile("assets/trap_1x.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	mouseImage, _, err = ebitenutil.NewImageFromFile("assets/mouse_1x.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

// gameLoop
func update(screen *ebiten.Image) error {
	// actual State (State = tileIndex)
	state := grid.getIndex(player.row, player.col)

	// agent choose and make action (epsillon only for debug)
	action, eps := agent.chooseAction(state)
	dir, _ := player.input(action)
	player.move(dir)

	//measure Reward
	reward := player.updateState(grid)
	newState := grid.getIndex(player.row, player.col)

	// Update Q
	agent.updateQ(state, action, float64(reward), newState)

	// Draw Grid
	w, h := tileImage.Size()

	for _, t := range grid.tiles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(t.rows*h)+padding/2, float64(t.cols*w)+padding/2)
		screen.DrawImage(tileImage, op)

		switch t.value {
		case 1:
			screen.DrawImage(cheeseImage, op)
		case 2:
			screen.DrawImage(trapImage, op)
		case 3:
			screen.DrawImage(trippleCheeseImage, op)
		}

		if q := agent.QTable[grid.getIndex(t.rows, t.cols)]; q != nil {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", q[0]), tileSize*t.rows+35+padding/2, tileSize*t.cols+5+padding/2)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", q[1]), tileSize*t.rows+35+padding/2, tileSize*t.cols+75+padding/2)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", q[2]), tileSize*t.rows+6+padding/2, tileSize*t.cols+40+padding/2)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", q[3]), tileSize*t.rows+60+padding/2, tileSize*t.cols+40+padding/2)
		}

	}

	// draw Player
	opPlayer := &ebiten.DrawImageOptions{}
	opPlayer.GeoM.Translate(float64(player.row*h)+padding/2, float64(player.col*w)+padding/2)
	screen.DrawImage(mouseImage, opPlayer)

	// draw  epsillon and episode
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Epsillon: %.2f", eps), 20, 5)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Episode: %v", episode), 20, 20)


	if ebiten.IsDrawingSkipped() {
		return nil
	}

	return nil
}

func initGrid() {
	grid = Grid{}
	//init grid
	for r := 0; r < screenHeight/tileSize; r++ {
		for c := 0; c < screenWidth/tileSize; c++ {
			grid.tiles = append(grid.tiles, Tile{r, c, 0})
		}
	}

	//  chesse -> reward
	grid.tiles[2].value = 1
	grid.tiles[14].value = 1
	grid.tiles[17].value = 1

	// mouse Trap -> punish
	grid.tiles[3].value = 2
	grid.tiles[1].value = 2
	grid.tiles[9].value = 2
	grid.tiles[21].value = 2
	grid.tiles[27].value = 2
	grid.tiles[9].value = 2
	grid.tiles[7].value = 2
	grid.tiles[11].value = 2
	grid.tiles[20].value = 2
	grid.tiles[22].value = 2
	grid.tiles[32].value = 2

	// tripple Cheese -> end of episode
	grid.tiles[33].value = 3

}

func main() {
	agent = NewAgent(4)
	player = mouse{0, 0, none}
	player.reset()
	ebiten.SetMaxTPS(30)
	if err := ebiten.Run(update, screenWidth+padding, screenHeight+padding, 1, "MouseCheese AI"); err != nil {
		log.Fatal(err)
	}
}
