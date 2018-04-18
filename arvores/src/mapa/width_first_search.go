package mapa

import (
  "fmt"
)

func (m *Mapa) BuscaEmLargura(start_x, start_y, end_x, end_y int) (int, *Estado) {
  e_inicial := NewEstado(uint(start_x), uint(start_y))
  e_final := NewEstado(uint(end_x), uint(end_y))

  return -1, m.sub_busca_em_largura([]*Estado{e_inicial}, e_final)
}

func (m *Mapa) sub_busca_em_largura(e_iniciais []*Estado, e_final *Estado) (*Estado) {
  if len(e_iniciais) == 0 {
    return nil
  }

  e_inicial := e_iniciais[0]

  for ; len(e_iniciais) == 0; {
    var estado_atual *Estado
    // pop
    estado_atual, e_iniciais = e_iniciais[0], e_iniciais[1:]

    // Norte
    if estado_atual.pos_x >= 0 || estado_atual.pos_y >= 0 || estado_atual.pos_x >= int(m.size_x) || estado_atual.pos_y >= int(m.size_y)

    fmt.Println(estado_atual)
  }

  return e_inicial
}
