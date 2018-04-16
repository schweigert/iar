package mapa

import (
  "fmt"
)

type Mapa struct {
  size_x, size_y uint
  grid [][]uint
  visited [][]bool
  edge [][]bool
}

func NewMapa(size_x, size_y uint) *Mapa {
  return &Mapa{
    size_x: size_x,
    size_y: size_y,
    grid: read_map_from_terminal(size_x, size_y),
    visited: create_map_mask(size_x, size_y, false),
    edge: create_map_mask(size_x, size_y, false),
  }
}

func (m *Mapa) Custo(x, y int) int {
  custos := [4]int{1, 5, 10, 15}

  return custos[m.grid[y][x]]
}

/*
 * Private
 * Methods
*/

func create_map_mask(size_x, size_y uint, def bool) [][]bool {
  grid := make([][]bool, size_y)

  for y := range(grid) {
    grid[y] = make([]bool, size_x)

    for x := range(grid[y]) {
      grid[y][x] = def
    }
  }

  return grid
}

func read_map_from_terminal(size_x, size_y uint) [][]uint {
  var value int

  grid := make([][]uint, size_y)

  for y := range(grid) {
    grid[y] = make([]uint, size_x)

    for x := range(grid[y]) {
      fmt.Scanf("%d", &value)
      grid[y][x] = uint(value)
    }
  }

  return grid
}
