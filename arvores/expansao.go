package main

import (
  "mapa"
)

const (
  MAPA_X = 42
  MAPA_Y = 42
  START_X = 0
  START_Y = 0
  END_X = 10
  END_Y = 10
)

func main() {
  mapa_atual := mapa.NewMapa(MAPA_X, MAPA_Y)
  mapa_atual.Imprimir()
  mapa_atual.DesenharImagem("Teste")

  mapa_atual.BuscaEmProfundidade(START_X, START_Y, END_X, END_Y)
}
