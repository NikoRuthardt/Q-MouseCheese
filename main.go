package main

import (
	"log"
  "fmt"
	"github.com/hajimehoshi/ebiten"
  "github.com/hajimehoshi/ebiten/ebitenutil"
  "github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 600
	screenHeight = 600
  tileSize = 100
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
  row,col int
  dir Dir
}

func (m *mouse) move(d Dir){
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


func (m *mouse) input() (Dir, bool){
    if inpututil.IsKeyJustPressed(ebiten.KeyUp){
      if m.col == 0 {
        return none,true
      }
      return up,true
    }
    if inpututil.IsKeyJustPressed(ebiten.KeyDown){
      if m.col == screenHeight / tileSize -1 {
        return none,true
      }
      return down,true
    }
    if inpututil.IsKeyJustPressed(ebiten.KeyLeft){
      if m.row == 0 {
        return none,true
      }
      return left,true
    }
    if inpututil.IsKeyJustPressed(ebiten.KeyRight){
      if m.row == screenWidth / tileSize -1 {
        return none,true
      }
      return right,true
    }

    return none, false
}

func (m *mouse) updateState(grid Grid) {
    nextTile := grid.getIndex(m.row, m.col)
    switch grid.tiles[nextTile].value {
    case 0:
      return
    case 1:
      grid.tiles[nextTile].value = 0
    case 2:
      m.reset()
    case 3:
      m.reset()
    }
}

func (m *mouse) reset(){
  m.col = 0
  m.row = 0
  initGrid()
}

// GameBoard

// Grid the game board
type Grid struct{
 tiles []Tile
}

func (g *Grid) getIndex(row, col int) int{
    for i, t := range g.tiles{
        if t.cols == col && t.rows == row {
          return i
        }
      }
      return -1
}

// Tile :each point on the grid is a tile
type Tile struct {
  rows,cols,value int
}


var(
  tileImage, cheeseImage, trapImage, trippleCheeseImage, mouseImage *ebiten.Image
  grid Grid
  player mouse
)


// init Pictures
func init(){
  var err error

  tileImage, _, _ = ebitenutil.NewImageFromFile("assets/tile_1x.png", ebiten.FilterDefault)
  cheeseImage, _, _ = ebitenutil.NewImageFromFile("assets/cheese_1x.png", ebiten.FilterDefault)
  trippleCheeseImage, _, _ = ebitenutil.NewImageFromFile("assets/trippleCheese_1x.png", ebiten.FilterDefault)
  trapImage, _, _ = ebitenutil.NewImageFromFile("assets/trap_1x.png", ebiten.FilterDefault)

  mouseImage, _, _ = ebitenutil.NewImageFromFile("assets/mouse_1x.png", ebiten.FilterDefault)

  if err != nil {
    log.Fatal(err)
  }
}

// gameLoop
func update(screen *ebiten.Image) error {

  if dir,input := player.input(); input {
    player.move(dir)
    player.updateState(grid)
  }

  w, h := tileImage.Size()

  // Draw Grid
  for _, t := range grid.tiles{
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(float64(t.rows * h),float64(t.cols*w))
    screen.DrawImage(tileImage, op)

    switch t.value {
    case 1:
      screen.DrawImage(cheeseImage, op)
    case 2:
      screen.DrawImage(trapImage, op)
    case 3:
      screen.DrawImage(trippleCheeseImage, op)
    }
    // ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", i), t.rows * h + 20, t.cols*w +20)
  }
  // draw Player
  plOp := &ebiten.DrawImageOptions{}
  plOp.GeoM.Translate(float64(player.row * h),float64(player.col*w))
  screen.DrawImage(mouseImage,plOp)

  ebitenutil.DebugPrintAt(screen, fmt.Sprintf("row:%v col:%v", player.row, player.col), 0,0)

  if ebiten.IsDrawingSkipped() {
  	return nil
  }

  return nil
}

func initGrid(){
  grid = Grid{}
  //init grid
  for r:=0; r < screenHeight/tileSize; r++ {
    for c:=0; c < screenWidth/tileSize; c++{
      grid.tiles = append(grid.tiles, Tile{r,c,0})
    }
  }

  // simple chesse -> reward
  grid.tiles[2].value = 1
  grid.tiles[10].value = 1
  grid.tiles[14].value = 1
  grid.tiles[12].value = 1

  // mouse Trap -> punish
  grid.tiles[3].value = 2
  grid.tiles[5].value = 2
  grid.tiles[7].value = 2
  grid.tiles[11].value = 2
  grid.tiles[22].value = 2
  grid.tiles[27].value = 2
  grid.tiles[32].value = 2

  // tripple Cheese -> end of episode
  grid.tiles[35].value = 3

}

func main() {

  player = mouse{0,0,none}
  player.reset()
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "MouseCheese AI"); err != nil {
		log.Fatal(err)
	}
}
