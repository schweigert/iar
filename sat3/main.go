package main

import (
  "os"
  "fmt"
  "time"
  "strconv"
  "strings"
  "math/rand"
  "github.com/Arafatk/glot"
)

const (
  CONST_FILE_NAME = "20-91.cnf"
  CONST_CLAUSES = 91
  CONST_VARS = 20
  CONST_INTERATIONS = 100000
  CONST_TESTS = 10
)

var (
  glot_spaces [CONST_INTERATIONS]float64
  glot_temperature [CONST_INTERATIONS]float64
  glot_energy [CONST_INTERATIONS]float64
  glot_i int
)

type SAT3Instance struct {
  clauses int
  expression [CONST_CLAUSES][3]int
  initial_set [CONST_VARS]bool
  final_set [CONST_VARS]bool
  work_set [CONST_VARS]bool
  randomizer *rand.Rand
}

func NewSAT3Instance() *SAT3Instance {
  seed := rand.NewSource(time.Now().UnixNano())
  sat3 := &SAT3Instance {
    clauses: CONST_CLAUSES,
    randomizer: rand.New(seed),
  }

  sat3.ReadExpression()
  sat3.RandomizeInitialInstance()

  return sat3
}

func (sat3 *SAT3Instance) RandomizeInitialInstance() {
  seed := rand.NewSource(time.Now().UnixNano())
  randomizer := rand.New(seed)

  fmt.Println("Initial Set:")

  for i := range(sat3.initial_set) {
    sat3.initial_set[i] = randomizer.Float64() > 0.50
    fmt.Printf("Set[%d]: %t\n", i, sat3.initial_set[i])
  }
}

func (sat3 *SAT3Instance) ReadExpression() {
  file, err := os.Open(CONST_FILE_NAME)

  if err != nil {
      panic(err)
  }

  defer file.Close()

  fmt.Println("Expression:")

  for i := range(sat3.expression) {
    var a, b, c, e int
    fmt.Fscanf(file, "%d %d %d %d",&a, &b, &c, &e)

    sat3.expression[i][0] = a
    sat3.expression[i][1] = b
    sat3.expression[i][2] = c

    if i == 0 {
      fmt.Printf(
        "  (\t%s\t v \t%s\t v \t%s\t)\n",
        strings.Replace(strconv.Itoa(a), "-", "¬", -1),
        strings.Replace(strconv.Itoa(b), "-", "¬", -1),
        strings.Replace(strconv.Itoa(c), "-", "¬", -1),
      )
    } else {
      fmt.Printf(
        "^ (\t%s\t v \t%s\t v \t%s\t)\n",
        strings.Replace(strconv.Itoa(a), "-", "¬", -1),
        strings.Replace(strconv.Itoa(b), "-", "¬", -1),
        strings.Replace(strconv.Itoa(c), "-", "¬", -1),
      )
    }
  }
}

func (sat3 *SAT3Instance) WorkScore() int {
  var score int

  for i := range(sat3.expression){
    a := sat3.expression[i][0]
    b := sat3.expression[i][1]
    c := sat3.expression[i][2]

    a_negative := a < 0
    b_negative := b < 0
    c_negative := c < 0

    if a_negative {
      a *= -1
    }

    if b_negative {
      b *= -1
    }

    if c_negative {
      c *= -1
    }

    a--
    b--
    c--

    a_bool := sat3.work_set[a]
    b_bool := sat3.work_set[b]
    c_bool := sat3.work_set[c]

    if a_negative {
      a_bool = !a_bool
    }

    if b_negative {
      b_bool = !b_bool
    }

    if c_negative {
      c_bool = !c_bool
    }

    if a_bool || b_bool || c_bool {
      score++
    }
  }

  return score
}

func (sat3 *SAT3Instance) CopyToWork() {
  for i := range(sat3.initial_set) {
    sat3.work_set[i] = sat3.initial_set[i]
  }
}

func (sat3 *SAT3Instance) CopyToFinal() {
  for i := range(sat3.work_set) {
    sat3.final_set[i] = sat3.work_set[i]
  }
}

func (sat3 *SAT3Instance) CopyFinalToWork() {
  for i := range(sat3.work_set) {
    sat3.work_set[i] = sat3.final_set[i]
  }
}

func (sat3 *SAT3Instance) RandomSearch(interations int) int {
  sat3.CopyToWork()
  sat3.CopyToFinal()
  score := sat3.WorkScore()

  for i := 0; i < interations; i++ {
    sat3.NewRandomWorkSet()
    new_score := sat3.WorkScore()

    if new_score > score {
      sat3.CopyToFinal()
      score = new_score
      fmt.Println("[", i, "] New best score: ", score)
    }
  }

  return score
}

func (sat3 *SAT3Instance) SimulatedAnnealingSearch(interations int) int {
  glot_i = 0

  sat3.CopyToWork()
  sat3.CopyToFinal()
  energy := sat3.WorkScore()

  for i := 0; i < interations; i++ {
    temperature := sat3.SimulatedAnnealingTemperature(i, interations)

    glot_spaces[glot_i] = float64(i) / float64(interations)
    glot_temperature[glot_i] = temperature

    sat3.NewSimulatedAnnealingWorkSet()
    new_energy := sat3.WorkScore()

    if sat3.SimulatedAnnealingRange(temperature, energy, new_energy) {
      sat3.CopyToFinal()
      if new_energy != energy {
        // fmt.Println("[", i, "] New energy: ", new_energy)
      }
      energy = new_energy
    }

    glot_energy[glot_i] += float64(energy) / CONST_TESTS


    sat3.CopyFinalToWork()

    glot_i++
  }

  return energy
}

func (sat3 *SAT3Instance) SimulatedAnnealingRange(energy float64, score, new_score int) bool {
  if new_score > score {
    return true
  }

  r := score - new_score

  if float64(r) <= float64(score) * energy {
    return true
  }
  return false
}

func (sat3 *SAT3Instance) NewSimulatedAnnealingWorkSet() {
  local := sat3.randomizer.Intn(len(sat3.work_set))

  sat3.work_set[local] = !sat3.work_set[local]
}

func (sat3 *SAT3Instance) SimulatedAnnealingTemperature(interation, interations int) float64 {
  t := (1 - float64(interation) / float64(interations))
  return t * t * t * t * t * t * t * t * t * t * t * t
}

func (sat3 *SAT3Instance) NewRandomWorkSet() {
  seed := rand.NewSource(time.Now().UnixNano())
  randomizer := rand.New(seed)

  for i := range(sat3.initial_set) {
    sat3.work_set[i] = randomizer.Float64() > 0.50
  }
}

func print_with_glot() {
  // Plot Energy
  dimensions := 2
  persist := false
  debug := false
  plot, _ := glot.NewPlot(dimensions, persist, debug)
  pointGroupName := "Temperature"
  style := "lines"
  points := [][]float64{glot_spaces[:], glot_temperature[:]}

  plot.AddPointGroup(pointGroupName, style, points)
  plot.SetTitle("Energy SAT3")
  plot.SetXLabel("Time")
  plot.SetYLabel("Temperature")
  plot.SavePlot("temperature.png")

  // Plot convergance
  plot, _ = glot.NewPlot(dimensions, persist, debug)
  pointGroupName = "Simulated Annealing"
  style = "lines"
  points = [][]float64{glot_spaces[:], glot_energy[:]}

  plot.AddPointGroup(pointGroupName, style, points)
  plot.SetTitle("Energy SAT3")
  plot.SetXLabel("Time")
  plot.SetYLabel("Energy")
  plot.SavePlot("convergence.png")
}

func main() {
  sat3 := NewSAT3Instance()
  for i := 0; i < 10; i++ {
    fmt.Println(sat3.SimulatedAnnealingSearch(CONST_INTERATIONS))
    // fmt.Println(sat3.RandomSearch(CONST_INTERATIONS))
  }

  print_with_glot()
}
