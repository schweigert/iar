package main

import (
  "os"
  "os/exec"
  "fmt"
  "sync"
  "time"
  "math"
  "image"
  "strconv"
  "math/rand"
  "image/png"
  "image/color"
)

// Three:
// CONST_ALPHA = 1.2
// CONST_K1 = 0.3 // GET
// CONST_K2 = 1.6 // PUT

// Two:
// CONST_ALPHA = 8
// CONST_K1 = 1.0 // GET
// CONST_K2 = 1.0 // PUT

// One:
// CONST_ALPHA = 6.5
// CONST_K1 = 0.75 // GET
// CONST_K2 = 2.5 // PUT

const (
  CONST_R = 2
  CONST_GROUPS = 3
  CONST_MAP_H = 72
  CONST_MAP_L = 72
  CONST_ANTS = 50
  CONST_DEAD_ANTS = 400

  CONST_COLORS = 4

  CONST_ALPHA = 6.5
  CONST_K1 = 0.75 // GET
  CONST_K2 = 2.5 // PUT

  CONST_FIGS = 3

  CONST_DEBUG_ONE_ANT = false
  CONST_DEBUG_FINISH = true
  CONST_DEBUG_STEPS = false
)

func create_random() *rand.Rand {
  return rand.New(rand.NewSource(time.Now().UnixNano()))
}

type DeadAnt struct {
  data []float64
  class int
}

func NewDeadAnt() *DeadAnt {
  var data []float64
  data = make([]float64, 2)

  var class int

  fmt.Scanf("%f,%f,", &data[0], &data[1])
  fmt.Scanf("%d", &class)

  return &DeadAnt{ class: class, data: data }
}

func (d1 *DeadAnt) Distancy(d2 *DeadAnt) float64 {
  sum := float64(0)
  for i := range d1.data {
    sum += (d1.data[i] - d2.data[i]) * (d1.data[i] - d2.data[i])
  }

  return math.Sqrt(sum)
}

type Ant struct {
  l, h, r int
  density, max_density float32
  function_f float64
  random *rand.Rand
  environment *Map
  dead_ant *DeadAnt
  mutex *sync.Mutex
}

func (a *Ant) Think() {
  a.UpdateFunctionF()
  a.Garbage()
  a.Walk()
  if CONST_DEBUG_ONE_ANT && a.function_f != 0.0{
    a.PrintRadius()
  }
}

func (a *Ant) PrintRadius() {
  if a.GetProbability() >= 0 {
    return
  }
  a.environment.PrintRadius(a.l, a.h, 10)
  fmt.Println("F:\t", a.function_f)
  fmt.Println("PUT:\t", a.PutProbability() * 100)
  fmt.Println("GET:\t", a.GetProbability() * 100)
  fmt.Println("GAR:\t", a.IsGarbagging())
  fmt.Println("CLS:\t", string('a' - 1 + a.dead_ant.class))
  fmt.Println("----------------")
}

func (a *Ant) GetDeadAnt() {
  a.dead_ant = a.environment.GetDeadAnt(a.l, a.h)
}

func (a *Ant) Garbage() {
  if a.IsGarbagging() {
    if a.random.Float64() <= a.PutProbability() {
      if a.environment.PutDeadAnt(a.l, a.h, a.dead_ant) {
        a.dead_ant = nil
      }
    }
  } else {
    if a.random.Float64() <= a.GetProbability() {
      a.dead_ant = a.environment.GetDeadAnt(a.l, a.h)
    }
  }
}

func (a *Ant) PutProbability() float64 {
  value := (a.function_f / (CONST_K2 + a.function_f))
  return (value * value)
}

func (a *Ant) GetProbability() float64 {
  value := (CONST_K1/(CONST_K1 + a.function_f))

  return (value * value)
}

func (a *Ant) UpdateFunctionF() {
  a.function_f = a.environment.FunctionF(a.l, a.h, a.r, a.dead_ant)
}

func (a *Ant) IsGarbagging() bool {
  return a.dead_ant != nil
}

func (a *Ant) IsOverAnt() bool {
  return a.environment.has_dead_ant_at(a.l, a.h)
}

func (a *Ant) RelativeDensity() float32 {
  return a.density / a.max_density
}

func (a *Ant) UpdateDensity() {
  a.density = a.environment.density(a.l, a.h, a.r)
  if a.density > a.max_density {
    a.max_density = a.density
  }
}

func (a *Ant) Walk() {
  a.l = (a.l + a.random.Intn(3) - 1)
  a.h = (a.h + a.random.Intn(3) - 1)

  if a.l == a.environment.l {
    a.l = 0
  }

  if a.h == a.environment.h {
    a.h = 0
  }

  if a.l < 0 {
    a.l = a.environment.l -1
  }

  if a.h < 0  {
    a.h = a.environment.h - 1
  }
}

func NewAnt(l, h int, env *Map) *Ant {
  ant := &Ant{
    random: create_random(),
    mutex: &sync.Mutex{},
    environment: env,
    dead_ant: nil,
    r : CONST_R,
  }

  ant.l = ant.random.Intn(l)
  ant.h = ant.random.Intn(h)

  return ant
}

type Map struct {
  l, h int
  environment [][]*DeadAnt
  ants []*Ant
  random *rand.Rand
  mutex *sync.Mutex
}

func (m *Map) PutDeadAnt(l, h int, ant *DeadAnt) bool {
  l %= m.l
  if l < 0 {
    l *= -1
  }

  h %= m.h
  if h < 0 {
    h *= -1
  }

  if m.environment[h][l] == nil {
    m.environment[h][l] = ant
    return true
  }

  return false
}

func (m *Map) GetDeadAnt(l, h int) *DeadAnt {
  l %= m.l
  if l < 0 {
    l *= -1
  }

  h %= m.h
  if h < 0 {
    h *= -1
  }

  ant := m.environment[h][l]
  m.environment[h][l] = nil

  return ant
}

func (m *Map) PrintRadius(l, h, r int) {
  for i := -r; i <= r; i++ {
    for j := -r; j <= r; j++ {
      if i == 0 && j == 0 {
        fmt.Print("@ ")
      } else {
        if m.has_dead_ant_at(i + l, j + h) {
          dead := m.GetDeadAnt(i + l, j + h)
          fmt.Print(string(dead.class + 'a' - 1)," ")
          m.PutDeadAnt(i + l, j + h, dead)
        } else {
          fmt.Print("  ")
        }
      }
    }
    fmt.Println("")
  }
}

func (m *Map) FunctionF(l, h, r int, ant *DeadAnt) float64 {
  ants := 0
  var sum float64

  sum = 0.0
  var hehe bool

  hehe = false

  if ant == nil {
    ant = m.GetDeadAnt(l, h)
    hehe = true

    if ant == nil {
      return 0.0
    }
  }

  for i := -r; i <= r; i++ {
    for j := -r; j <= r; j++ {
      other_ant := m.GetDeadAnt(l + i, h + j);
      if other_ant == nil {
        continue
      }

      ants++
      distancy := ant.Distancy(other_ant)
      sum += 1 - (distancy/CONST_ALPHA)

      m.PutDeadAnt(l + i, h + j, other_ant)
    }
  }

  if hehe {
    m.PutDeadAnt(l, h, ant)
  }

  if ants == 0 {
    return 0.0
  }

  ret := (1.0/float64(r * r)) * sum

  if ret <= 0.0 {
    return 0.0
  }

  return ret
}

func (m *Map) has_dead_ant_at(l, h int) bool {
  l %= m.l
  if l < 0 {
    l *= -1
  }

  h %= m.h
  if h < 0 {
    h *= -1
  }

  if m.environment[h][l] == nil {
    return false
  }
  return true
}

func (m *Map) density(l, h, r int) float32 {
  sum := 0;
  total_sum := 0;

  for i := -r; i < r; i++ {
    for j := -r; j < r; j++ {
      total_sum += 1
      if m.has_dead_ant_at(l + i, h + j) {
        sum += 1
      }
    }
  }

  return float32(sum) / float32(total_sum)
}

func (m* Map) Interate(n int) {
  ma := -1
  for i := 0; i < n; i++ {

    if ((100 * i) / n) > ma {
      ma = ((100 * i) / n)
      fmt.Println("Walking... ", (100 * i) / n, "%")
    }

    if i % 1000 == 0 && CONST_DEBUG_STEPS {
      m.Print()
    }
    for a := range m.ants {
      m.ants[a].Think()
    }
  }
}

func (m* Map) InterateMod(n, groups, group_id int, wg *sync.WaitGroup) {
  for i := 0; i < n; i++ {
    if i % groups != group_id {
      continue
    }

    for a := range m.ants {
      m.ants[a].Think()
    }
  }

  wg.Done()
}

func (m* Map) ParallelInterate(n, groups int) {
  var wg sync.WaitGroup

  for i := 0; i < groups; i++ {
    wg.Add(1)
    go m.InterateMod(n, groups, i, &wg)
  }

  wg.Wait()
}

func (m* Map) CreateAnts(n int) {
  ants := make([]*Ant, n)
  for i := range ants {
    ants[i] = NewAnt(m.l, m.h, m)
  }

  m.ants = ants
}

func (m* Map) CreateDeadAnts(n int) {
  if n == 0 {
    return
  }

  h := m.random.Intn(m.h)
  l := m.random.Intn(m.l)

  if m.environment[h][l] == nil {
    m.environment[h][l] = NewDeadAnt()

    m.CreateDeadAnts(n-1)
  } else {
    m.CreateDeadAnts(n)
  }
}

func (m* Map) Finish() {
  for i := range m.ants {
    fmt.Println("Finishing ant: ", i)
    for ; m.ants[i].IsGarbagging(); {
      if CONST_DEBUG_FINISH {
        m.ants[i].PrintRadius()
      }
      m.ants[i].Think()
    }
  }
}

func (m* Map) Drawn(iteration int) {
  img := image.NewRGBA(image.Rect(0, 0, m.l * 8, m.h * 8))

  for i := range m.environment {
    for j := range m.environment[i] {
      if m.environment[i][j] == nil {
        for h := 0; h < 8; h++ {
          for l := 0; l < 8; l++ {
            img.Set(i*8 +h, j*8 +l, color.RGBA{0, 0, 0, 255})
          }
        }
      } else {
        for h := 0; h < 8; h++ {
          for l := 0; l < 8; l++ {
            img.Set(i*8 +h, j*8 +l, color.RGBA{uint8((255/CONST_COLORS)* m.environment[i][j].class), 100, 255, 255})
          }
        }

      }
    }
  }

  f, _ := os.OpenFile(strconv.Itoa(CONST_R) + "-" + strconv.Itoa(iteration + 1) + ".png", os.O_WRONLY|os.O_CREATE, 0600)
  defer f.Close()

  png.Encode(f, img)
}

func (m* Map) Print() {
  cmd := exec.Command("clear")
  cmd.Stdout = os.Stdout
  cmd.Run()
  fmt.Println("--------------------------------------------------")
  for i := range m.environment {
    for j := range m.environment[i] {
      if m.environment[i][j] == nil {
        fmt.Print(" ")
      } else {
        fmt.Print(string(-1 + 'a' + m.environment[i][j].class))
      }
    }
    fmt.Print("\n")
  }
  return
  for i := range m.ants {
    fmt.Println(m.ants[i].l, m.ants[i].h, m.ants[i].density)
  }
}

func NewMap(l, h int) *Map {
  create_dead_ants_matrix := func (row, col int) [][]*DeadAnt {
    m := make([][]*DeadAnt, row)
    for i := range m {
      m[i] = make([]*DeadAnt, col)

      for j := range m[i] {
        m[i][j] = nil
      }
    }

    return m
  }

  m := &Map{
    l: l,
    h: h,
    environment: create_dead_ants_matrix(h,l),
    random: create_random(),
    mutex: &sync.Mutex{},
  }

  return m
}

func main() {
  m := NewMap(CONST_MAP_L, CONST_MAP_H)
  m.CreateDeadAnts(CONST_DEAD_ANTS)
  m.CreateAnts(CONST_ANTS)


  m.Drawn(-1)
  for i := 0; i < CONST_FIGS; i++ {
  //for {
    fmt.Println("Walking... ", 1000000 * i, " steps")
    m.Interate(1000000)
    m.Finish()
    fmt.Println("Drawning")
    m.Drawn(i)
  }
  m.Print()
}
