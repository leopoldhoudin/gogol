package main

import (
  tm "github.com/buger/goterm"
  "time"
  "fmt"
  "io/ioutil"
  "strings"
  "flag"
  "runtime"
)

func main() {
	var filename string
  var speed int
  var maxIterations int
	flag.StringVar(&filename, "filename", "", "--filename myfile.gogol")
	flag.IntVar(&speed, "speed", 500, "--speed 500")
	flag.IntVar(&maxIterations, "max", 100, "--max 100")
	flag.Parse()

  state := loadInitialState(filename)

  runtime.GOMAXPROCS(2)

  c := make(chan [][]int)

  iteration := 1
  for {
    go getNextState(state, c)
    render(state, iteration, maxIterations)

    time.Sleep(time.Duration(speed) * time.Millisecond)

    state = <-c
    iteration += 1
    if iteration > maxIterations {
      break
    }
  }

  for {
    time.Sleep(time.Second)
  }
}

func loadInitialState(filename string) [][]int {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }

  lines := strings.Split(string(data), "\n")

  if lines[len(lines) - 1] == "" {
    lines = lines[:len(lines) - 1]
  }

  h := len(lines)
  w := len(lines[0])

  initialState := make([][]int, w)
  for i := range initialState {
    initialState[i] = make([]int, h)
  }

  for y, row := range lines {
    for x, char := range row {
      if char == '1' {
        initialState[x][y] = 1
      }
    }
  }

  return initialState
}

/* *** *** *** *** *** *** RENDERING *** *** *** *** *** *** */

func render(state [][]int, iteration int, maxIterations int) {
  tm.Clear()

  tm.MoveCursor(3, 1)
  tm.Print(fmt.Sprintf("GOGOL - Iteration: %d/%d", iteration, maxIterations))

  for x, row := range state {
    for y := range row {
      if state[x][y] > 0 {
        fillBlock(x, y)
      } else {
        emptyBlock(x, y)
      }
    }
  }

  tm.MoveCursor(1, 1)
  tm.Flush()
}

const block = "\u2588\u2588"
const empty = "\u2591\u2591"

func fillBlock(x int, y int) {
  tm.MoveCursor(2 * x + 1, y + 2)
  tm.Print(tm.Color(block, tm.WHITE))
}

func emptyBlock(x int, y int) {
  tm.MoveCursor(2 * x + 1, y + 2)
  tm.Print(empty)
}

/* *** *** *** *** *** *** STATE *** *** *** *** *** *** */

func getNextState(prevState [][]int, c chan [][]int) {
  w := len(prevState)
  h := len(prevState[0])

  nextState := make([][]int, w)
  for i := range nextState {
    nextState[i] = make([]int, h)
  }

  for x := 0; x < w; x++ {
    for y := 0; y < h; y++ {
      aliveNeighbors := getAliveNeighbors(prevState, x, y, w, h)

      if prevState[x][y] > 0 {
        if aliveNeighbors == 2 || aliveNeighbors == 3 {
          nextState[x][y] = 1
        }
      } else {
        if aliveNeighbors == 3 {
          nextState[x][y] = 1
        }
      }
    }
  }

  c <- nextState
}

var offsets = []int{ -1, 0, 1 }

func getAliveNeighbors(state [][]int, x int, y int, w int, h int) int {
  var aliveNeighbors int = 0

  for _, dx := range offsets {
    X := x + dx
    if X < 0 || X >= w {
      continue
    }

    for _, dy := range offsets {
      Y := y + dy
      if Y < 0 || Y >= h || (dx == 0 && dy == 0) {
        continue
      }

      aliveNeighbors += state[X][Y]
    }
  }

  return aliveNeighbors
}
