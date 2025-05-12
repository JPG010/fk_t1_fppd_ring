// Código exemplo para o trabaho de sistemas distribuidos (eleicao em anel)
// By Cesar De Rose - 2022

package main

import (
	"fmt"
	"sync"
	"time"
)

type mensagem struct {
	tipo  int    // 1: eleição, 2: falha, 3: recuperação, 4: novo líder, 5: fim
	corpo [3]int // [id_origem, id_atual, id_lider]
}

var (
	chans = []chan mensagem{
		make(chan mensagem, 1), // Adicionando buffer para evitar bloqueio
		make(chan mensagem, 1),
		make(chan mensagem, 1),
		make(chan mensagem, 1),
	}
	controle = make(chan int, 1) // Buffer para evitar bloqueio
	wg       sync.WaitGroup
)

func ElectionControler(in chan int) {
	defer wg.Done()

	var temp mensagem

	// Simula falha do processo 0
	temp.tipo = 2
	chans[3] <- temp
	fmt.Printf("Controle: processo 0 falhou\n")
	<-in

	// Simula falha do processo 1
	temp.tipo = 2
	chans[0] <- temp
	fmt.Printf("Controle: processo 1 falhou\n")
	<-in

	// Inicia eleição
	temp.tipo = 1
	temp.corpo = [3]int{2, 2, 0} // Processo 2 inicia eleição
	chans[1] <- temp
	fmt.Printf("Controle: iniciando eleição pelo processo 2\n")

	time.Sleep(1 * time.Second)

	// Simula recuperação do processo 0
	temp.tipo = 3
	chans[3] <- temp
	fmt.Printf("Controle: processo 0 recuperado\n")
	<-in

	time.Sleep(1 * time.Second)

	// Envia mensagem de fim para todos os processos
	temp.tipo = 5
	for i := 0; i < 4; i++ {
		chans[i] <- temp
		fmt.Printf("Controle: enviando mensagem de fim para processo %d\n", i)
	}

	fmt.Println("\n   Processo controlador concluído\n")
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	var actualLeader int = leader
	var bFailed bool = false

	for {
		select {
		case temp := <-in:
			fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])

			if temp.tipo == 5 {
				fmt.Printf("%2d: finalizando processo\n", TaskId)
				return
			}

			if bFailed {
				controle <- -5
				continue
			}

			switch temp.tipo {
			case 1: // Mensagem de eleição
				if temp.corpo[0] == TaskId {
					// Mensagem voltou ao iniciador
					actualLeader = temp.corpo[1]
					// Anuncia novo líder
					newMsg := mensagem{tipo: 4, corpo: [3]int{TaskId, actualLeader, actualLeader}}
					select {
					case out <- newMsg:
						fmt.Printf("%2d: novo líder eleito: %d\n", TaskId, actualLeader)
					default:
						fmt.Printf("%2d: não foi possível anunciar novo líder\n", TaskId)
					}
				} else {
					// Adiciona seu ID se for maior
					if TaskId > temp.corpo[1] {
						temp.corpo[1] = TaskId
					}
					select {
					case out <- temp:
						fmt.Printf("%2d: passando mensagem de eleição\n", TaskId)
					default:
						fmt.Printf("%2d: não foi possível passar mensagem de eleição\n", TaskId)
					}
				}

			case 2: // Falha
				bFailed = true
				fmt.Printf("%2d: falho %v\n", TaskId, bFailed)
				controle <- -5

			case 3: // Recuperação
				bFailed = false
				fmt.Printf("%2d: recuperado\n", TaskId, bFailed)
				controle <- -5

			case 4: // Novo líder
				actualLeader = temp.corpo[2]
				fmt.Printf("%2d: novo líder é %d\n", TaskId, actualLeader)
				if !bFailed {
					select {
					case out <- temp:
						fmt.Printf("%2d: propagando novo líder\n", TaskId)
					default:
						fmt.Printf("%2d: não foi possível propagar novo líder\n", TaskId)
					}
				}
			}

			fmt.Printf("%2d: líder atual %d\n", TaskId, actualLeader)

		case <-time.After(3 * time.Second):
			fmt.Printf("%2d: timeout - finalizando processo\n", TaskId)
			return
		}
	}
}

func main() {
	wg.Add(5)

	// Criar os processos do anel
	go ElectionStage(0, chans[3], chans[0], 0)
	go ElectionStage(1, chans[0], chans[1], 0)
	go ElectionStage(2, chans[1], chans[2], 0)
	go ElectionStage(3, chans[2], chans[3], 0)

	fmt.Println("\n   Anel de processos criado")

	// Criar o processo controlador
	go ElectionControler(controle)

	fmt.Println("\n   Processo controlador criado\n")

	wg.Wait()
}
