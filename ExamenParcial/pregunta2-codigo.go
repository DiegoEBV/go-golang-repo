//go:build ignore

// Pregunta 2: procesamiento de lotes de imágenes.
// Ejecutar: go run -race pregunta2-codigo.go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const numLotes = 50

type Auditoria struct {
	registrado bool
	idLote     int
	cantidad   int
}

var (
	muAcum     sync.Mutex // protege acumulador
	acumulador = 0

	muAudit   sync.Mutex // protege auditoria (distinto del anterior)
	auditoria Auditoria
)

func duracion(tam int) time.Duration {
	switch tam {
	case 100:
		return 300 * time.Millisecond
	case 200:
		return 400 * time.Millisecond
	default:
		return 500 * time.Millisecond
	}
}

func procesar(id, tam int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(duracion(tam)) // trabajo simulado: sin ningún mutex tomado

	muAcum.Lock()
	acumulador += tam
	muAcum.Unlock()

	muAudit.Lock() // consulta de 'registrado' y escritura completa: una sola sección crítica
	if !auditoria.registrado {
		auditoria = Auditoria{registrado: true, idLote: id, cantidad: tam}
	}
	muAudit.Unlock()
}

func main() {
	rng := rand.New(rand.NewSource(2026))
	opciones := []int{100, 200, 300}
	tamanos := make([]int, numLotes) // se genera antes de lanzar goroutines y no se modifica
	esperado := 0
	for i := range tamanos {
		tamanos[i] = opciones[rng.Intn(3)]
		esperado += tamanos[i]
	}

	var wg sync.WaitGroup
	for id := 0; id < numLotes; id++ {
		wg.Add(1)
		go procesar(id, tamanos[id], &wg)
	}
	wg.Wait() // el reporte se imprime solo cuando todas terminaron

	fmt.Printf("Total esperado : %d\n", esperado)
	fmt.Printf("Total calculado: %d\n", acumulador)
	fmt.Printf("Registro ganador: idLote=%d cantidad=%d registrado=%v\n",
		auditoria.idLote, auditoria.cantidad, auditoria.registrado)
	fmt.Printf("¿Total coincide? %v | ¿cantidad == tamaño del lote %d? %v\n",
		esperado == acumulador, auditoria.idLote, auditoria.cantidad == tamanos[auditoria.idLote])
}
