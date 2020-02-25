package main

import "fmt"
import "os"

type Coord struct {
	x int16
	y int16
}

type Gri [][]uint8

const FREE = 5
const TAKEN = 6

func main() {
	grid := make(Gri, 20, 20)
	for i := 0; i < 20; i++ {
		g := make([]uint8, 30, 30)
		for j := 0; j < 30; j++ {
			g[j] = FREE
		}
		grid[i] = g
	}

	for {
		cleanGrid(grid)
		var N, P uint8
		fmt.Scan(&N, &P)

		players := make([]Coord, N, N)

		for i := uint8(0); i < N; i++ {
			var X0, Y0, X1, Y1 int16
			fmt.Scan(&X0, &Y0, &X1, &Y1)

			if X1 != -1 {
				grid[Y1][X1] = TAKEN + i
				grid[Y0][X0] = TAKEN + i
			} else {
			    removePlayer(grid, i)
			}
			
			players[i] = Coord{X1, Y1}
		}

		move := "LEFT"
		me := players[P]
		maxPoints := 0
		//fmt.Fprintln(os.Stderr, "Curr:", me)

		potentialMoveScore := evalPoint(Coord{me.x, me.y - 1}, grid, players, P)
		// fmt.Fprintln(os.Stderr, potentialMoveScore, "U")
		if potentialMoveScore > maxPoints {
			maxPoints = potentialMoveScore
			move = "UP"
		}

		potentialMoveScore = evalPoint(Coord{me.x + 1, me.y}, grid, players, P)
		// fmt.Fprintln(os.Stderr, potentialMoveScore, "R")
		if potentialMoveScore > maxPoints {
			maxPoints = potentialMoveScore
			move = "RIGHT"
		}

		potentialMoveScore = evalPoint(Coord{me.x, me.y + 1}, grid, players, P)
		// fmt.Fprintln(os.Stderr, potentialMoveScore, "D")
		if potentialMoveScore > maxPoints {
			maxPoints = potentialMoveScore
			move = "DOWN"
		}

		potentialMoveScore = evalPoint(Coord{me.x - 1, me.y}, grid, players, P)
		// fmt.Fprintln(os.Stderr, potentialMoveScore, "L")
		if potentialMoveScore > maxPoints {
			maxPoints = potentialMoveScore
			move = "LEFT"
		}

		fmt.Println(move)
	}
}

func doMinMax(grid Gri, players []Coord, mn int) string {
    
}

func evalPoint(point Coord, grid [][]uint8, players []Coord, mn uint8) int {
	if point.x >= 0 && point.x < 30 && point.y >= 0 && point.y < 20 && grid[point.y][point.x] < TAKEN {
		gc := deepCopy(grid)
		var pc = make([]Coord, len(players), len(players))
		copy(pc, players)
		pc[mn] = Coord{point.x, point.y}

		// fmt.Fprintln(os.Stderr, "Curr:", players[mn])
		// fmt.Fprintln(os.Stderr, "Move to:", point)
		// fmt.Fprintln(os.Stderr, "cont:", grid[point.y][point.x])
		gc[point.y][point.x] = TAKEN + mn
		fillGrid(gc, pc)

// 		printGrid(gc)

		mc, _ := count(gc, mn) // minmax here
		return mc
	}

	return 0
}

func fillGrid(grid [][]uint8, players []Coord) {
	points := make([]Coord, 0, 100)

	for i, p := range players {
		free := getFreePoints(grid, p)

		for _, f := range free {
			grid[f.y][f.x] = uint8(i)
			points = append(points, f)
		}
	}

	for len(points) > 0 {
		p := points[0]
		free := getFreePoints(grid, p)
		if len(free) > 0 {
			player := grid[p.y][p.x]
			for _, f := range free {
				grid[f.y][f.x] = player
				points = append(points, f)
			}
		}
		points = points[1:]
	}
}

func getFreePoints(grid [][]uint8, point Coord) []Coord {
	res := make([]Coord, 0, 3)
	
	if point.x == -1 {
	    return res
	}
	
	if point.y > 0 && grid[point.y - 1][point.x] == FREE {
		res = append(res, Coord{point.x, point.y - 1})
	}

	if point.x < 29 && grid[point.y][point.x + 1] == FREE {
		res = append(res, Coord{point.x + 1, point.y})
	}

	if point.y < 19 && grid[point.y + 1][point.x] == FREE {
		res = append(res, Coord{point.x, point.y + 1})
	}

	if point.x > 0 && grid[point.y][point.x - 1] == FREE {
		res = append(res, Coord{point.x - 1, point.y})
	}

	return res
}

func count(grid [][]uint8, me uint8) (int, int) {
	mc, oc := 0, 0
	for i := 0; i < 20; i++ {
		for j := 0; j < 30; j++ {
			g := grid[i][j]

			if g == me {
				mc++
			} else if g != FREE && g < TAKEN {
				oc++
			}
		}
	}
	return mc, oc
}

func cleanGrid(grid Gri) {
	for i := 0; i < 20; i++ {
		for j := 0; j < 30; j++ {
			if grid[i][j] < TAKEN {
				grid[i][j] = FREE
			}
		}
	}
}

func deepCopy(grid Gri) Gri {
	res := make(Gri, 20, 20)
	for i := 0; i < 20; i++ {
		res[i] = make([]uint8, 30, 30)
		copy(res[i], grid[i])
	}
	return res
}

func printGrid(grid [][]uint8) {
	for i := 0; i < 20; i++ {
		fmt.Fprintln(os.Stderr, grid[i])
	}
}

func removePlayer(grid Gri, pn uint8) {
    for i := 0; i < 20; i++ {
		for j := 0; j < 30; j++ {
			if grid[i][j] == TAKEN + pn {
				grid[i][j] = FREE
			}
		}
	}
}
