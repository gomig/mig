package questions

// Wizard questions
type Wizard struct {
	qs  []Question
	ans map[string]string
}

// Init initialized wizard
func (w *Wizard) Init() {
	w.qs = make([]Question, 0)
	w.ans = make(map[string]string)
}

// AddQuestion append question
func (w *Wizard) AddQuestion(q Question) {
	w.qs = append(w.qs, q)
}

// Start wizard
func (w *Wizard) Start() {
	for _, q := range w.qs {
		w.ans[q.Name] = q.Ask()
	}
}

// Result Get wizard result
func (w *Wizard) Result(name string) string {
	return w.ans[name]
}
