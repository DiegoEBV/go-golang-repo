//go:build ignore

// Pregunta 1 (b): versión CORRECTA con un mutex.
// Ejecutar: go run -race pregunta1-b-correcta.go
package main

import (
	"fmt"
	"sync"
)

const (
	creditosIniciales = 300
	solicitudes       = 800
	costo             = 2
)

var (
	mu         sync.Mutex // protege saldo, aceptadas y rechazadas
	saldo      = creditosIniciales
	aceptadas  = 0
	rechazadas = 0
)

func solicitar(wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock() // comprobación + descuento + contadores = una sola sección crítica
	if saldo >= costo {
		saldo -= costo
		aceptadas++
	} else {
		rechazadas++
	}
	mu.Unlock()
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < solicitudes; i++ {
		wg.Add(1)
		go solicitar(&wg)
	}
	wg.Wait()
	// Tras Wait ya no hay goroutines concurrentes: la lectura es segura.
	fmt.Printf("[CORRECTA] saldo=%d aceptadas=%d rechazadas=%d total=%d\n",
		saldo, aceptadas, rechazadas, aceptadas+rechazadas)
	if saldo != 0 || aceptadas != 150 || rechazadas != 650 {
		fmt.Println("ERROR: resultado inesperado")
	} else {
		fmt.Println("OK: saldo 0, 150 aceptadas, 650 rechazadas")
	}
}
