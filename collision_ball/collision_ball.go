package collision_ball

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	WALL          = 1
	PLAYER        = 1 << 3
	SYMBOL_WALL   = "◼︎"
	SYMBOL_PLAYER = "●"
	FPS           = 20
)

var directionGrp [][]int = [][]int{
	{-1, -1}, // 左上
	{1, -1},  // 左下
	{-1, 1},  // 右上
	{1, 1},   // 右下
}

type Player struct {
	x         int
	y         int
	direction []int
}

func (p *Player) getDirection() ([]int, error) {
	prob := []int{25, 25, 25, 25}
	pick := -1
	total := func() int {
		sum := 0
		for _, val := range prob {
			sum += val
		}
		return sum
	}()

	seed := rand.Intn(total + 1)
	sum := 0
	for idx, val := range prob {
		sum += val
		if sum > seed {
			pick = idx
			break
		}
	}

	if pick == -1 {
		return []int{}, errors.New("random pick error")
	}
	return directionGrp[pick], nil
}

func (p *Player) move(gameMap [][]byte, maxX int, maxY int) error {
	if p.x == 0 && p.y == 0 {
		p.x = rand.Intn(maxX-2) + 1
		p.y = rand.Intn(maxY-2) + 1
		gameMap[p.y][p.x] = PLAYER
		return nil
	}

	if len(p.direction) == 0 {
		p.direction, _ = p.getDirection()
	}

	// 確保玩家位置不會超出地圖邊界
	if p.x+p.direction[1] < 1 || p.x+p.direction[1] >= maxX-1 {
		p.direction[1] *= -1 // 如果準備超出左右邊界，則調整水平方向
	}
	if p.y+p.direction[0] < 1 || p.y+p.direction[0] >= maxY-1 {
		p.direction[0] *= -1 // 如果準備超出上下邊界，則調整垂直方向
	}

	p.x += p.direction[1]
	p.y += p.direction[0]
	gameMap[p.y][p.x] = PLAYER
	return nil
}

type Game struct {
	height  int
	width   int
	gameMap [][]byte
	player  Player
}

func (g *Game) drawPlayer() error {
	if g.height == 0 || g.width == 0 {
		return errors.New("height or width is 0")
	}

	if err := g.player.move(g.gameMap, g.width, g.height); err != nil {
		return err
	}

	return nil
}

func (g *Game) render() {
	g.height = 15
	g.width = 30

	// init slice
	g.gameMap = make([][]byte, g.height)
	for h := 0; h < g.height; h++ {
		for w := 0; w < g.width; w++ {
			g.gameMap[h] = make([]byte, g.width)
		}
	}
	for h := 0; h < g.height; h++ {
		for w := 0; w < g.width; w++ {
			if h == 0 || h == g.height-1 || w == 0 || w == g.width-1 {
				g.gameMap[h][w] = WALL
			}
		}
	}

	// init player
	g.drawPlayer()

	// init wall
	buf := new(bytes.Buffer)
	for _, h := range g.gameMap {
		for _, w := range h {
			switch w {
			case WALL:
				buf.WriteString(SYMBOL_WALL)
			case PLAYER:
				buf.WriteString(SYMBOL_PLAYER)
			default:
				buf.WriteString(" ")
			}
		}
		buf.WriteString("\n")
	}

	fmt.Fprint(os.Stdout, "\033[2J\033[1;1H") // 讓輸出不要一直往下寫
	fmt.Fprint(os.Stdout, buf.String())
	fmt.Fprint(os.Stdout, g.player.y, g.player.x)
}

func (g *Game) NewGame() {
	frameDuration := time.Duration(1000/FPS) * time.Millisecond
	for {
		st := time.Now()
		time.Sleep(frameDuration - time.Since(st))
		g.render()
	}
}

// ======================================================================================================

func Run() {
	myGame := &Game{}
	myGame.NewGame()
}
