/* Pregunta 3: dos procesos con x compartida. La condición y la asignación
 * NO están en atomic ni d_step: son pasos separados. */
byte x = 2;

active proctype P() {
p1: do
    :: x < 3 ->            /* p1: evaluar condición (verdadera) -> pasa a p2 */
p2:    x = x + 1           /* p2: asignación; luego vuelve a p1 */
    :: else -> break       /* p1: condición falsa -> sale a p3 (fin) */
    od
    /* p3: fin = final del cuerpo del proctype (estado final válido) */
}

active proctype Q() {
q1: do
    :: x > 1 ->            /* q1: condición verdadera -> pasa a q2 */
q2:    x = x - 1           /* q2: asignación; luego vuelve a q1 */
    :: else -> break       /* q1: condición falsa -> q3 (fin) */
    od
}
