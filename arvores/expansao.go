package main

import (
  "mapa"
)

const (
  MAPA_X = 42
  MAPA_Y = 42
)

func main() {
  mapa_atual := mapa.NewMapa(MAPA_X, MAPA_Y)
  mapa_atual.Imprimir()
  mapa_atual.DesenharImagem("Teste")
}
