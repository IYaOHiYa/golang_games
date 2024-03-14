package main

import (
	collision_ball "my_game/collision_ball"
	_ "my_game/gobang"
	"time"
)

func main() {
	time.Local, _ = time.LoadLocation("Asia/Taipei")

	collision_ball.Run()
	// gobang.Run()
}
