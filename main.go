package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"log"
)

const (
	screenWidth  = 600
	screenHeight = 600
	tileSize     = 100
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
		return 0
	case 1:
		grid.tiles[nextTile].value = 0
    return 1
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
  episode int
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
  dir,_ := player.input(action)
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
		op.GeoM.Translate(float64(t.rows*h), float64(t.cols*w))
		screen.DrawImage(tileImage, op)

		switch t.value {
		case 1:
			screen.DrawImage(cheeseImage, op)
		case 2:
			screen.DrawImage(trapImage, op)
		case 3:
			screen.DrawImage(trippleCheeseImage, op)
		}

	}

	// draw Player
	opPlayer := &ebiten.DrawImageOptions{}
	opPlayer.GeoM.Translate(float64(player.row*h), float64(player.col*w))
	screen.DrawImage(mouseImage, opPlayer)

  // draw episode and epsillon
  ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Episode: %v", episode), 0, 20)
  ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Epsillon: %.2f", eps), 0, 0)

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
  grid.tiles[15].value = 1
  grid.tiles[17].value = 1
  grid.tiles[19].value = 1
  grid.tiles[23].value = 1

	// mouse Trap -> punish
	grid.tiles[3].value = 2
	grid.tiles[5].value = 2
	grid.tiles[7].value = 2
	grid.tiles[22].value = 2
	grid.tiles[32].value = 2

	// tripple Cheese -> end of episode
	grid.tiles[35].value = 3

}

func main() {
	agent = NewAgent(4)
	player = mouse{0, 0, none}
	player.reset()
  ebiten.SetMaxTPS(30)
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "MouseCheese AI"); err != nil {
		log.Fatal(err)
	}
}
