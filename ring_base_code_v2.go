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
		make(chan mensagem, 1), // P0 (controlador)
		make(chan mensagem, 1), // P1
		make(chan mensagem, 1), // P2
		make(chan mensagem, 1), // P3
		make(chan mensagem, 1), // P4
		make(chan mensagem, 1), // P5
	}
	controle = make(chan int, 1)
	wg       sync.WaitGroup
	done     = make(chan bool)
)

func ElectionControler(in chan int) {
	defer wg.Done()

	var temp mensagem

	// 1. Falha em P5 e eleição
	fmt.Println("\n=== 1. Falha em P5 e eleição ===")
	temp.tipo = 2
	chans[5] <- temp
	fmt.Printf("P0: P5 falhou\n")
	<-in

	temp.tipo = 1
	temp.corpo = [3]int{2, 2, 0}
	chans[2] <- temp
	fmt.Printf("P0: P2 inicia eleição\n")
	time.Sleep(2 * time.Second)

	// 2. Falha em P4 e eleição
	fmt.Println("\n=== 2. Falha em P4 e eleição ===")
	temp.tipo = 2
	chans[4] <- temp
	fmt.Printf("P0: P4 falhou\n")
	<-in

	temp.tipo = 1
	temp.corpo = [3]int{3, 3, 0}
	chans[3] <- temp
	fmt.Printf("P0: P3 inicia eleição\n")
	time.Sleep(2 * time.Second)

	// 3. Recuperação de P4 e eleição
	fmt.Println("\n=== 3. Recuperação de P4 e eleição ===")
	temp.tipo = 3
	chans[4] <- temp
	fmt.Printf("P0: P4 recuperado\n")
	<-in

	temp.tipo = 1
	temp.corpo = [3]int{4, 4, 0}
	chans[4] <- temp
	fmt.Printf("P0: P4 inicia eleição\n")
	time.Sleep(2 * time.Second)

	// 4. Recuperação de P5 e eleição
	fmt.Println("\n=== 4. Recuperação de P5 e eleição ===")
	temp.tipo = 3
	chans[5] <- temp
	fmt.Printf("P0: P5 recuperado\n")
	<-in

	temp.tipo = 1
	temp.corpo = [3]int{5, 5, 0}
	chans[5] <- temp
	fmt.Printf("P0: P5 inicia eleição\n")
	time.Sleep(2 * time.Second)

	// Finaliza todos os processos
	fmt.Println("\n=== Finalizando processos ===")
	temp.tipo = 5
	for i := 1; i < 6; i++ {
		select {
		case chans[i] <- temp:
			fmt.Printf("P0: enviando mensagem de fim para P%d\n", i)
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("P0: timeout ao enviar mensagem de fim para P%d\n", i)
		}
	}
	close(done)
	fmt.Println("\n   Processo controlador concluído\n")
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	var actualLeader int = leader
	var bFailed bool = false

	for {
		select {
		case temp := <-in:
			fmt.Printf("P%d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])

			if temp.tipo == 5 {
				fmt.Printf("P%d: finalizando processo\n", TaskId)
				return
			}

			if bFailed {
				select {
				case controle <- -5:
				case <-time.After(100 * time.Millisecond):
				}
				continue
			}

			switch temp.tipo {
			case 1: // Mensagem de eleição
				if temp.corpo[0] == TaskId {
					actualLeader = temp.corpo[1]
					newMsg := mensagem{tipo: 4, corpo: [3]int{TaskId, actualLeader, actualLeader}}
					select {
					case out <- newMsg:
						fmt.Printf("P%d: novo líder eleito: P%d\n", TaskId, actualLeader)
					case <-time.After(100 * time.Millisecond):
						fmt.Printf("P%d: não foi possível anunciar novo líder\n", TaskId)
					}
					fmt.Printf("P%d: LÍDER ELEITO NA RODADA: P%d\n", TaskId, actualLeader)
				} else {
					if TaskId > temp.corpo[1] {
						temp.corpo[1] = TaskId
					}
					select {
					case out <- temp:
						fmt.Printf("P%d: passando mensagem de eleição\n", TaskId)
					case <-time.After(100 * time.Millisecond):
						fmt.Printf("P%d: não foi possível passar mensagem de eleição\n", TaskId)
					}
				}

			case 2: // Falha
				bFailed = true
				fmt.Printf("P%d: falho %v\n", TaskId, bFailed)
				select {
				case controle <- -5:
				case <-time.After(100 * time.Millisecond):
				}

			case 3: // Recuperação
				bFailed = false
				fmt.Printf("P%d: recuperado\n", TaskId)
				select {
				case controle <- -5:
				case <-time.After(100 * time.Millisecond):
				}

			case 4: // Novo líder
				actualLeader = temp.corpo[2]
				fmt.Printf("P%d: novo líder é P%d\n", TaskId, actualLeader)
				if !bFailed {
					select {
					case out <- temp:
						fmt.Printf("P%d: propagando novo líder\n", TaskId)
					case <-time.After(100 * time.Millisecond):
						fmt.Printf("P%d: não foi possível propagar novo líder\n", TaskId)
					}
				}
			}

			fmt.Printf("P%d: líder atual P%d\n", TaskId, actualLeader)

		case <-time.After(1 * time.Second):
			select {
			case <-done:
				fmt.Printf("P%d: finalizando processo por sinal de término\n", TaskId)
				return
			default:
			}
		}
	}
}

func main() {
	wg.Add(6) // 5 processos + 1 controlador

	go ElectionStage(1, chans[5], chans[1], 0) // P1
	go ElectionStage(2, chans[1], chans[2], 0) // P2
	go ElectionStage(3, chans[2], chans[3], 0) // P3
	go ElectionStage(4, chans[3], chans[4], 0) // P4
	go ElectionStage(5, chans[4], chans[5], 0) // P5

	fmt.Println("\n   Anel de processos criado")
	go ElectionControler(controle)
	fmt.Println("\n   Processo controlador criado\n")
	wg.Wait()
}
