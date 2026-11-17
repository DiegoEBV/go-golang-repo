//go:build ignore

// Pregunta 4: barrera reutilizable (dos molinetes) con Mutex + semáforos propios.
// Ejecutar: go run -race pregunta4-codigo.go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ---------- Semáforo ----------
type Semaforo struct {
	mu   sync.Mutex
	val  int
	cond *sync.Cond
}

func NewSemaforo(n int) *Semaforo {
	s := &Semaforo{val: n}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *Semaforo) Wait() {
	s.mu.Lock()
	for s.val == 0 { // bucle: protege contra despertares espurios
		s.cond.Wait() // suspende (no es espera activa) y libera mu
	}
	s.val--
	s.mu.Unlock()
}

func (s *Semaforo) Signal() {
	s.mu.Lock()
	s.val++
	s.cond.Signal()
	s.mu.Unlock()
}

// ---------- Barrera reutilizable ----------
type Barrera struct {
	n         int
	mu        sync.Mutex // protege count
	count     int
	molinete1 *Semaforo // puerta de entrada: se abre cuando llegan los n
	molinete2 *Semaforo // puerta de salida: se abre cuando salieron de la 1ª los n
}

func NewBarrera(n int) *Barrera {
	return &Barrera{n: n, molinete1: NewSemaforo(0), molinete2: NewSemaforo(0)}
}

func (b *Barrera) Esperar() {
	// Fase de llegada
	b.mu.Lock()
	b.count++
	if b.count == b.n { // el último en llegar habilita al grupo
		for i := 0; i < b.n; i++ {
			b.molinete1.Signal()
		}
	}
	b.mu.Unlock()
	b.molinete1.Wait()

	// Fase de salida (rearme)
	b.mu.Lock()
	b.count--
	if b.count == 0 { // el último en salir de la 1ª puerta abre la 2ª
		for i := 0; i < b.n; i++ {
			b.molinete2.Signal()
		}
	}
	b.mu.Unlock()
	b.molinete2.Wait()
}

// ---------- Trabajadores ----------
const (
	N = 4
	R = 3
)

var (
	logMu sync.Mutex
	t0    = time.Now()
)

func registrar(w, ronda, fase int, ev string) {
	logMu.Lock()
	fmt.Printf("[%4dms] trabajador=%d ronda=%d fase=%d %s\n",
		time.Since(t0).Milliseconds(), w, ronda, fase, ev)
	logMu.Unlock()
}

func trabajador(id int, b *Barrera, dur [R][2]time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	for r := 1; r <= R; r++ {
		for fase := 1; fase <= 2; fase++ {
			registrar(id, r, fase, "INICIO")
			time.Sleep(dur[r-1][fase-1])
			registrar(id, r, fase, "FIN") // FIN antes de entrar a la barrera
			b.Esperar()                   // misma barrera, al final de cada fase
		}
	}
}

func main() {
	b := NewBarrera(N)
	var wg sync.WaitGroup
	for id := 0; id < N; id++ {
		rng := rand.New(rand.NewSource(int64(100 + id))) // generador local por trabajador
		var dur [R][2]time.Duration
		for r := 0; r < R; r++ {
			for f := 0; f < 2; f++ {
				dur[r][f] = time.Duration(50+rng.Intn(101)) * time.Millisecond
			}
		}
		wg.Add(1)
		go trabajador(id, b, dur, &wg)
	}
	wg.Wait()
	fmt.Println("Todos los trabajadores terminaron las 3 rondas.")
}
