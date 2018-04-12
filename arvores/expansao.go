package main

import (
  "fmt"
  "mapa"
)

const (
  MAPA_X = 42
  MAPA_Y = 42
)

func main() {
  fmt.Println(mapa.NewMapa(MAPA_X, MAPA_Y))
}
