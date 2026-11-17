//go:build ignore

// Pregunta 1 (a): versión INCORRECTA, sin protección.
// Ejecutar: go run -race pregunta1-a-incorrecta.go
package main

import (
	"fmt"
	"runtime"
	"sync"
)

const (
	creditosIniciales = 300
	solicitudes       = 800
	costo             = 2
)

var (
	saldo      = creditosIniciales
	aceptadas  = 0
	rechazadas = 0
)

func solicitar(wg *sync.WaitGroup) {
	defer wg.Done()
	if saldo >= costo { // consulta sin protección
		runtime.Gosched() // ensancha la ventana entre la consulta y la actualización
		saldo -= costo    // actualización sin protección
		aceptadas++
	} else {
		rechazadas++
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < solicitudes; i++ {
		wg.Add(1)
		go solicitar(&wg)
	}
	wg.Wait()
	fmt.Printf("[INCORRECTA] saldo=%d aceptadas=%d rechazadas=%d (esperado: 0/150/650)\n",
		saldo, aceptadas, rechazadas)
}
