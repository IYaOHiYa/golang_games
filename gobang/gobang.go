package gobang

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	PLATE_ROW = 10
	PLATE_COL = 10

	PLAYER   = 1
	OPPONENT = 2
)

// 玩家圖標
var PlayerImg map[int]string = map[int]string{
	PLAYER:   "❍", // ⚪️
	OPPONENT: "✕", // 🔴
}

// 用戶輸入邏輯控制
func InputVal(typeStr string, v *int, limit int) {
	// fmt.Scan(row)
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Please input %s: ", typeStr)
		t, _ := r.ReadString('\n')
		t = strings.TrimSpace(t)
		num, err := strconv.Atoi(t)
		if err != nil {
			fmt.Println("Please input correct type!!")
			fmt.Println("==============================")
			continue
		}

		if num < 1 || num >= limit {
			fmt.Printf("The range value is: %d~%d\n", 1, limit-1)
			fmt.Println("==============================")
			continue
		}

		*v = num
		return
	}
}

// 初始化棋盤
func InitPlate(plateRow, plateCol int) [][]int {
	plate := make([][]int, plateRow)
	for i := 0; i < plateRow; i++ {
		plate[i] = make([]int, plateCol)
	}
	return plate
}

// 標出玩家位置
func MarkPoint(plate [][]int, row int, col int, player int) error {
	if plate[row][col] > 0 {
		return errors.New("this place has already been taken, please try again")
	}
	plate[row][col] = player
	return nil
}

// AI防守
func MarkAIPoint(plate [][]int, posX int, posY int) {
	// direction := map[string][2]int{
	// 	"top":       {1, 0},
	// 	"left":      {0, 1},
	// 	"left_top":  {-1, -1},
	// 	"left_down": {1, -1},
	// }

	/**
		TO DO... bot logic
	**/

	// if err := MarkPoint(plate, row+1, col+1, OPPONENT); err != nil {
	// 	panic("AI position error, reason: " + err.Error())
	// }
}

// 檢查五連線邏輯，針對四條線來檢查是否形成五連線(上下，左右，左上右下，左下右上)
func IsWin(plate [][]int, player int, posX int, posY int) bool {
	// 上下
	connectCnt := 0
	minX := math.Max(1, float64(posX-5))
	maxX := math.Min(PLATE_COL-1, float64(posX+5))
	for x := minX; x <= maxX; x++ {
		if plate[int(x)][posY] == player {
			connectCnt++
		}
		if connectCnt >= 5 {
			return true
		}
	}

	// 左右
	connectCnt = 0
	minY := math.Max(1, float64(posY-5))
	maxY := math.Min(PLATE_ROW-1, float64(posY+5))
	for y := minY; y <= maxY; y++ {
		if plate[posX][int(y)] == player {
			connectCnt++
		}
		if connectCnt >= 5 {
			return true
		}
	}

	connectCnt = 1
	for i := 1; i <= 5; i++ {
		// 左上
		if posX-i > 0 && posY-i > 0 && plate[posX-i][posY-i] == player {
			connectCnt++
		}

		// 右下
		if posX+i < PLATE_COL && posY+i < PLATE_ROW && plate[posX+i][posY+i] == player {
			connectCnt++
		}
	}
	if connectCnt >= 5 {
		return true
	}

	connectCnt = 1
	for i := 1; i <= 5; i++ {
		// 左下
		if posX+i < PLATE_COL && posY-i > 0 && plate[posX+i][posY-i] == player {
			connectCnt++
		}

		// 右上
		if posX-i > 0 && posY+i < PLATE_ROW && plate[posX-i][posY+i] == player {
			connectCnt++
		}
	}
	return connectCnt >= 5
}

// 檢查是否滿盤
func IsFull(plate [][]int) bool {
	cnt := 0
	for x := range plate {
		for y := range plate[x] {
			if plate[x][y] > 0 {
				cnt++
			}
		}
	}
	return cnt >= (PLATE_ROW-1)*(PLATE_COL-1)
}

// 繪製棋盤
func RenderPlate(plate [][]int) {
	buf := new(bytes.Buffer)
	for x := range plate {
		for y := range plate[x] {
			// 第一格第一個位置
			if x == 0 && y == 0 {
				buf.WriteString("   ")
				continue
			}

			// 第一行
			if x == 0 {
				buf.WriteString(strconv.Itoa(y) + "  ")
				continue
			}

			// 每一列第一個元素
			if x > 0 && y == 0 {
				buf.WriteString(strconv.Itoa(x))
				continue
			}

			player := plate[x][y]
			if player > 0 {
				buf.WriteString("  " + PlayerImg[player])
			} else {
				buf.WriteString("  " + "+")
			}

		}
		buf.WriteString("\n")
	}

	fmt.Fprint(os.Stdout, "\033[2J\033[1;1H") // 讓輸出不要一直往下寫
	fmt.Print(buf.String())
}

func Run() {
	plateBoard := InitPlate(PLATE_ROW, PLATE_COL)
	RenderPlate(plateBoard)

	var row int
	var col int
	for {
		InputVal("row", &row, PLATE_ROW)
		InputVal("col", &col, PLATE_COL)

		tt := time.Now()
		if err := MarkPoint(plateBoard, row, col, PLAYER); err != nil {
			fmt.Println(err)
			continue
		}
		// MarkAIPoint(plateBoard, row, col)
		RenderPlate(plateBoard)

		if IsWin(plateBoard, PLAYER, row, col) {
			fmt.Println("Congrats!! you win~")
			break
		}

		if IsFull(plateBoard) {
			fmt.Println("No winner，board is full")
			break
		}
		fmt.Println(time.Since(tt))
	}
}
