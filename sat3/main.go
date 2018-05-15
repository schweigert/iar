package main

import (
  "os"
  "fmt"
  "time"
  "strconv"
  "strings"
  "math/rand"
)

const (
  CONST_FILE_NAME = "20-91.cnf"
  CONST_CLAUSES = 20
  CONST_VARS = 91
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
  sat3.CopyToWork()
  sat3.CopyToFinal()
  score := sat3.WorkScore()

  for i := 0; i < interations; i++ {
    energy := sat3.SimulatedAnnealingEnergy(i, interations)
    sat3.NewSimulatedAnnealingWorkSet(energy)
    new_score := sat3.WorkScore()

    if SimulatedAnnealingRange(energy, score, new_score) {
      sat3.CopyToFinal()
      score = new_score
      fmt.Println("[", i, "] New score: ", score)
    }

    sat3.CopyFinalToWork()
  }

  return score
}

func SimulatedAnnealingRange(energy float64, score, new_score int) bool {
  return float64(new_score) * (1 + energy * 0.3) > float64(score) || new_score > score
}

func (sat3 *SAT3Instance) NewSimulatedAnnealingWorkSet(energy float64) {
  set_len := len(sat3.work_set)
  edit := int(float64(set_len) * energy)

  if edit <= 1 {
    edit = 1
  }

  for i := 0; i < edit; i++ {
    local := sat3.randomizer.Intn(set_len)

    sat3.work_set[local] = !sat3.work_set[local]
  }
}

func (sat3 *SAT3Instance) SimulatedAnnealingEnergy(interation, interations int) float64 {
  return 1.0 - (float64(interation) / float64(interations)) * (float64(interation) / float64(interations))
}

func (sat3 *SAT3Instance) NewRandomWorkSet() {
  seed := rand.NewSource(time.Now().UnixNano())
  randomizer := rand.New(seed)

  for i := range(sat3.initial_set) {
    sat3.work_set[i] = randomizer.Float64() > 0.50
  }
}

func main() {
  sat3 := NewSAT3Instance()
  fmt.Println(sat3.SimulatedAnnealingSearch(250000))
  fmt.Println(sat3.final_set)
  fmt.Println(sat3.RandomSearch(250000))
  fmt.Println(sat3.final_set)
}
