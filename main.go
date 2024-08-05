package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const tamanhoTabuleiro = 3

type QuebraCabeca struct {
	Tabuleiro      [tamanhoTabuleiro][tamanhoTabuleiro]int
	XVazio, YVazio int
	Caminho        []([tamanhoTabuleiro][tamanhoTabuleiro]int)
}

// Estrutura para a fila de prioridade
type PrioridadeEstado struct {
	estado     QuebraCabeca
	heuristica int
}

type FilaPrioridade []PrioridadeEstado

func (p FilaPrioridade) Len() int            { return len(p) }
func (p FilaPrioridade) Less(i, j int) bool  { return p[i].heuristica < p[j].heuristica }
func (p FilaPrioridade) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *FilaPrioridade) Push(x interface{}) { *p = append(*p, x.(PrioridadeEstado)) }
func (p *FilaPrioridade) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}

func main() {
	var quebraCabeca QuebraCabeca
	for {
		quebraCabeca = gerarTabuleiro()
		if isTabuleiroSolucionavel(quebraCabeca.Tabuleiro) {
			break
		}
	}

	caminhoSolucao, encontrado := buscaGulosa(quebraCabeca)

	if encontrado {
		fmt.Println("Solução encontrada:")
		for _, estado := range caminhoSolucao {
			imprimirTabuleiro(estado)
			fmt.Println()
		}
	} else {
		imprimirTabuleiro(quebraCabeca.Tabuleiro)
		fmt.Println("Nenhuma solução encontrada.")
	}
}

func gerarTabuleiro() QuebraCabeca {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	numeros := rng.Perm(9)
	var tabuleiro [tamanhoTabuleiro][tamanhoTabuleiro]int

	for i := 0; i < tamanhoTabuleiro; i++ {
		for j := 0; j < tamanhoTabuleiro; j++ {
			tabuleiro[i][j] = numeros[i*tamanhoTabuleiro+j]
		}
	}

	var posX, posY int
	for i := 0; i < tamanhoTabuleiro; i++ {
		for j := 0; j < tamanhoTabuleiro; j++ {
			if tabuleiro[i][j] == 0 {
				posX, posY = i, j
				break
			}
		}
	}

	return QuebraCabeca{Tabuleiro: tabuleiro, XVazio: posX, YVazio: posY}
}

func isTabuleiroSolucionavel(tabuleiro [tamanhoTabuleiro][tamanhoTabuleiro]int) bool {
	inversoes := 0
	tabuleiroLinear := [9]int{}

	for i := 0; i < tamanhoTabuleiro; i++ {
		for j := 0; j < tamanhoTabuleiro; j++ {
			tabuleiroLinear[i*tamanhoTabuleiro+j] = tabuleiro[i][j]
		}
	}

	for i := 0; i < len(tabuleiroLinear); i++ {
		if tabuleiroLinear[i] == 0 {
			continue
		}
		for j := i + 1; j < len(tabuleiroLinear); j++ {
			if tabuleiroLinear[j] == 0 {
				continue
			}
			if tabuleiroLinear[i] > tabuleiroLinear[j] {
				inversoes++
			}
		}
	}
	return inversoes%2 == 0
}

func buscaGulosa(inicio QuebraCabeca) ([][tamanhoTabuleiro][tamanhoTabuleiro]int, bool) {
	fila := &FilaPrioridade{}
	heap.Push(fila, PrioridadeEstado{estado: inicio, heuristica: calcularHeuristica(inicio.Tabuleiro)})
	visitados := make(map[[tamanhoTabuleiro * tamanhoTabuleiro]int]struct{})

	for fila.Len() > 0 {
		p := heap.Pop(fila).(PrioridadeEstado)
		estadoAtual := p.estado

		visitados[criarChave(estadoAtual.Tabuleiro)] = struct{}{}

		if verificarEstadoFinal(estadoAtual.Tabuleiro) {
			caminhoSolucao := append(estadoAtual.Caminho, estadoAtual.Tabuleiro)
			return caminhoSolucao, true
		}

		for _, novoEstado := range obterEstadosAdjacentes(estadoAtual) {
			chave := criarChave(novoEstado.Tabuleiro)
			if _, existe := visitados[chave]; !existe {
				heuristica := calcularHeuristica(novoEstado.Tabuleiro)
				heap.Push(fila, PrioridadeEstado{estado: novoEstado, heuristica: heuristica})
			}
		}
	}
	return nil, false
}

func calcularHeuristica(tabuleiro [tamanhoTabuleiro][tamanhoTabuleiro]int) int {
	heuristica := 0
	estadoFinal := [tamanhoTabuleiro][tamanhoTabuleiro]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 0},
	}

	for i := 0; i < tamanhoTabuleiro; i++ {
		for j := 0; j < tamanhoTabuleiro; j++ {
			valor := tabuleiro[i][j]
			if valor == 0 {
				continue
			}
			finalX, finalY := encontrarPosicao(valor, estadoFinal)
			heuristica += abs(i-finalX) + abs(j-finalY)
		}
	}
	return heuristica
}

func encontrarPosicao(valor int, tabuleiro [tamanhoTabuleiro][tamanhoTabuleiro]int) (int, int) {
	for i := 0; i < tamanhoTabuleiro; i++ {
		for j := 0; j < tamanhoTabuleiro; j++ {
			if tabuleiro[i][j] == valor {
				return i, j
			}
		}
	}
	return -1, -1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func obterEstadosAdjacentes(estado QuebraCabeca) []QuebraCabeca {
	deltaX := []int{1, -1, 0, 0}
	deltaY := []int{0, 0, 1, -1}

	var vizinhos []QuebraCabeca

	for i := 0; i < 4; i++ {
		novoX := estado.XVazio + deltaX[i]
		novoY := estado.YVazio + deltaY[i]

		if isPosicaoValida(novoX, novoY) {
			novoTabuleiro := estado.Tabuleiro

			novoTabuleiro[estado.XVazio][estado.YVazio], novoTabuleiro[novoX][novoY] = novoTabuleiro[novoX][novoY], novoTabuleiro[estado.XVazio][estado.YVazio]
			novoCaminho := make([][tamanhoTabuleiro][tamanhoTabuleiro]int, len(estado.Caminho))

			copy(novoCaminho, estado.Caminho)

			novoCaminho = append(novoCaminho, novoTabuleiro)
			vizinhos = append(vizinhos, QuebraCabeca{Tabuleiro: novoTabuleiro, XVazio: novoX, YVazio: novoY, Caminho: novoCaminho})
		}
	}
	return vizinhos
}

func isPosicaoValida(x, y int) bool {
	return x >= 0 && x < tamanhoTabuleiro && y >= 0 && y < tamanhoTabuleiro
}

func verificarEstadoFinal(tabuleiro [tamanhoTabuleiro][tamanhoTabuleiro]int) bool {
	estadoFinal := [tamanhoTabuleiro][tamanhoTabuleiro]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 0},
	}
	return tabuleiro == estadoFinal
}

func criarChave(tabuleiro [tamanhoTabuleiro][tamanhoTabuleiro]int) [tamanhoTabuleiro * tamanhoTabuleiro]int {
	var chave [tamanhoTabuleiro * tamanhoTabuleiro]int
	for i := 0; i < tamanhoTabuleiro; i++ {
		for j := 0; j < tamanhoTabuleiro; j++ {
			chave[i*tamanhoTabuleiro+j] = tabuleiro[i][j]
		}
	}

	return chave
}

func imprimirTabuleiro(tabuleiro [tamanhoTabuleiro][tamanhoTabuleiro]int) {
	fmt.Println("+" + strings.Repeat("-", tamanhoTabuleiro*2+1) + "+")

	for _, linha := range tabuleiro {
		fmt.Print("|")
		for _, valor := range linha {
			fmt.Printf(" %d", valor)
		}
		fmt.Println(" |")
	}

	fmt.Println("+" + strings.Repeat("-", tamanhoTabuleiro*2+1) + "+")
}
