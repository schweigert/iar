package main

import (
  "mapa"
  "fmt"
)

const (
  MAPA_X = 42
  MAPA_Y = 42
  START_X = 10
  START_Y = 10
  END_X = 41
  END_Y = 41
)

func main() {
  mapa_atual := mapa.NewMapa(MAPA_X, MAPA_Y)
  mapa_atual.Imprimir()
  mapa_atual.DesenharImagem("Teste")

  custo, estados := mapa_atual.BuscaEmLargura(START_X, START_Y, END_X, END_Y)
  fmt.Println(custo, estados)
}
