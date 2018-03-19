package main

import (
  "fmt"
)

type DeadAnt struct {}

type Map struct {
  l, h int
  environment *[][]*DeadAnt
}

func (m* Map) CreateDeadAnts (n int) {
  if n == 0 {
    return
  }
}

func NewMap(l, h int) *Map {
  m := &Map{
    l: l,
    h: h,
    environment: create_dead_ants_matrix(h, l),
  }

  return m
}

func create_dead_ants_matrix(row, col int) *[][]*DeadAnt {
  m := make([][]*DeadAnt, row)
  for i := range m {
    m[i] = make([]*DeadAnt, col)

    for j := range m[i] {
      m[i][j] = nil
    }
  }

  return &m
}

func main() {
  fmt.Println("Booting...")

  m := NewMap(30, 30)
}
