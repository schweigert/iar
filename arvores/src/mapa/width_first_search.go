package mapa

import (
  "strconv"
)


func (m *Mapa) BuscaEmLargura(start_x, start_y, end_x, end_y int) (int, *Estado) {
  e_inicial := NewEstado(uint(start_x), uint(start_y))
  e_final := NewEstado(uint(end_x), uint(end_y))

  return m.sub_busca_em_largura([]*Estado{e_inicial}, e_final)
}

func (m *Mapa) sub_busca_em_largura(e_iniciais []*Estado, e_final *Estado) (int, *Estado) {

  estados_alcancados := []*Estado{}

  if len(e_iniciais) == 0 {
    return 0, nil
  }

  for i := range e_iniciais {
    e_atual := e_iniciais[i]

    x := int(e_atual.pos_x)
    y := int(e_atual.pos_y)

    if m.pode_visitar(x - 1, y) {
      m.visited[y][x] = true
      e_norte := NewEstado(x - 1, y)
      e_atual.norte = e_norte

      if e_norte.e_final(e_final) {
        return
      }

      estados_alcancados = append(estados_alcancados, e_norte)
    }
  }

  return 0, nil
}

func (m *Mapa) pode_visitar(x, y int) bool {
  if x >= int(m.size_x) || y >= int(m.size_y) || x < 0 || y < 0 {
    return false
  }

  if m.visited[y][x] {
    return false
  }

  return true
}
