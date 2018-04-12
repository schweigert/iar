package mapa

import (
  "fmt"
)

type Mapa struct {
  size_x, size_y uint
  grid [][]uint
}

func NewMapa(size_x, size_y uint) *Mapa {
  return &Mapa{
    size_x: size_x,
    size_y: size_y,
    grid: read_map_from_terminal(size_x, size_y),
  }
}

func read_map_from_terminal(size_x, size_y uint) [][]uint {
  var value int

  grid := make([][]uint, size_y)

  for y := range(grid) {
    grid[y] = make([]uint, size_x)
  }

  for y := uint(0); y < size_y; y++ {
    for x := uint(0); x < size_x; x++ {
      fmt.Scanf("%d", &value)
      grid[y][x] = uint(value)
    }
  }

  return grid
}
