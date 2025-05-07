// Código exemplo para o trabalho de sistemas distribuídos (eleição em anel)
// By Cesar De Rose - 2022

package main

import (
	"fmt"
	"sync"
)

// Constantes para tipos de mensagem
const (
	MSG_ELECTION = 1 // Mensagem de eleição
	MSG_FAILURE = 2  // Mensagem de falha
	MSG_RECOVERY = 3 // Mensagem de recuperação
	MSG_COORDINATOR = 4 // Mensagem de novo coordenador
)

// Estrutura de mensagem usada para comunicação entre os processos
type mensagem struct {
	tipo  int    // Tipo da mensagem
	corpo [3]int // Conteúdo da mensagem [id_origem, id_atual, id_maior]
}

var (
	// Vetor de canais para formar o anel de eleição
	chans = []chan mensagem{
		make(chan mensagem), // Canal do processo 0
		make(chan mensagem), // Canal do processo 1
		make(chan mensagem), // Canal do processo 2
		make(chan mensagem), // Canal do processo 3
	}
	controle = make(chan int) // Canal de controle para comunicação com o controlador
	wg       sync.WaitGroup   // WaitGroup para sincronizar a execução das goroutines
)

// Função que controla o processo de eleição
func ElectionControler(in chan int) {
	defer wg.Done()

	var temp mensagem

	// Simular falha do coordenador (processo 0)
	temp.tipo = MSG_FAILURE
	chans[3] <- temp
	fmt.Printf("Controle: mudar o processo 0 para falho\n")
	fmt.Printf("Controle: confirmação %d\n", <-in)

	// Aguardar um tempo para ver a eleição
	fmt.Println("Controle: aguardando eleição...")
	fmt.Printf("Controle: confirmação %d\n", <-in)

	// Simular falha do processo 1
	temp.tipo = MSG_FAILURE
	chans[0] <- temp
	fmt.Printf("Controle: mudar o processo 1 para falho\n")
	fmt.Printf("Controle: confirmação %d\n", <-in)

	// Enviar mensagem de recuperação para o processo 0
	temp.tipo = MSG_RECOVERY
	chans[3] <- temp
	fmt.Printf("Controle: recuperar o processo 0\n")
	fmt.Printf("Controle: confirmação %d\n", <-in)

	fmt.Println("\n   Processo controlador concluído\n")
}

// Função que representa cada estágio do anel de eleição
func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	var actualLeader int = leader
	var bFailed bool = false
	var isCoordinator bool = (TaskId == leader)

	fmt.Printf("%2d: iniciando (líder inicial: %d, é coordenador: %v)\n", TaskId, actualLeader, isCoordinator)

	for {
		// Ler mensagem recebida no canal de entrada
		temp := <-in
		fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])

		// Processar a mensagem recebida
		switch temp.tipo {
		case MSG_FAILURE:
			{
				bFailed = true
				isCoordinator = false
				fmt.Printf("%2d: falho %v, é coordenador: %v\n", TaskId, bFailed, isCoordinator)
				fmt.Printf("%2d: líder atual %d\n", TaskId, actualLeader)
				controle <- TaskId

				// Se o processo que falhou era o coordenador, iniciar eleição
				if actualLeader == TaskId {
					// Iniciar eleição
					electionMsg := mensagem{
						tipo: MSG_ELECTION,
						corpo: [3]int{TaskId, TaskId, TaskId},
					}
					out <- electionMsg
					fmt.Printf("%2d: iniciando eleição após falha do coordenador\n", TaskId)
				}
			}
		case MSG_RECOVERY:
			{
				bFailed = false
				fmt.Printf("%2d: recuperado %v, é coordenador: %v\n", TaskId, bFailed, isCoordinator)
				fmt.Printf("%2d: líder atual %d\n", TaskId, actualLeader)
				controle <- TaskId
			}
		case MSG_ELECTION:
			{
				if !bFailed {
					// Se a mensagem voltou ao iniciador
					if temp.corpo[0] == TaskId {
						// Escolher o maior ID como novo coordenador
						newLeader := temp.corpo[2]
						fmt.Printf("%2d: eleição concluída, novo líder: %d\n", TaskId, newLeader)
						
						// Enviar mensagem de novo coordenador
						coordMsg := mensagem{
							tipo: MSG_COORDINATOR,
							corpo: [3]int{newLeader, newLeader, newLeader},
						}
						out <- coordMsg
						controle <- TaskId
					} else {
						// Adicionar ID à mensagem de eleição se for maior
						if TaskId > temp.corpo[2] {
							temp.corpo[2] = TaskId
						}
						out <- temp
					}
				}
			}
		case MSG_COORDINATOR:
			{
				if !bFailed {
					actualLeader = temp.corpo[1]
					isCoordinator = (TaskId == actualLeader)
					fmt.Printf("%2d: novo líder definido: %d, é coordenador: %v\n", TaskId, actualLeader, isCoordinator)
					
					// Repassar mensagem de coordenador
					if temp.corpo[0] != TaskId {
						out <- temp
					}
				}
			}
		default:
			{
				fmt.Printf("%2d: não conheço este tipo de mensagem\n", TaskId)
			}
		}

		// Se o processo falhou, não processa mais mensagens
		if bFailed {
			fmt.Printf("%2d: processo falho, encerrando (era coordenador: %v)\n", TaskId, isCoordinator)
			return
		}
	}
}

// Função principal
func main() {
	wg.Add(5) // 4 processos + 1 controlador

	// Criar os processos do anel de eleição
	go ElectionStage(0, chans[3], chans[0], 0) // Processo 0 (líder inicial)
	go ElectionStage(1, chans[0], chans[1], 0) // Processo 1
	go ElectionStage(2, chans[1], chans[2], 0) // Processo 2
	go ElectionStage(3, chans[2], chans[3], 0) // Processo 3

	fmt.Println("\n   Anel de processos criado")

	// Criar o processo controlador
	go ElectionControler(controle)

	fmt.Println("\n   Processo controlador criado\n")

	wg.Wait()
}
