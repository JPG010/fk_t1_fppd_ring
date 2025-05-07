package main

import (
	"fmt"
	"sync"
	"time"
)

type mensagem struct {
	tipo    int    // 0:eleição, 1:líder, 2:falha, 3:recuperação, 9:encerrar
	corpo   []int  // IDs dos processos
	sender  int    // Processo que enviou
	origem  int    // Quem iniciou a eleição
	message string // Descrição
}

var (
	chans = []chan mensagem{
		make(chan mensagem, 10),
		make(chan mensagem, 10),
		make(chan mensagem, 10),
		make(chan mensagem, 10),
		make(chan mensagem, 10),
	}
	controle = make(chan string, 10)
	wg       sync.WaitGroup
	done     = make(chan struct{})
)

func main() {
	wg.Add(6) // 5 processos + controlador

	// Iniciar processos (0-4 representam 1-5)
	go ElectionStage(0, chans[4], chans[0], 4)
	go ElectionStage(1, chans[0], chans[1], 4)
	go ElectionStage(2, chans[1], chans[2], 4)
	go ElectionStage(3, chans[2], chans[3], 4)
	go ElectionStage(4, chans[3], chans[4], 4)

	fmt.Println("\nAnel de processos criado (5 processos, líder inicial: 5)")

	go ElectionControler(controle)
	fmt.Println("\nProcesso controlador criado\n")

	wg.Wait()
	fmt.Println("Todos os processos foram encerrados")
}

func ElectionControler(in chan string) {
	defer wg.Done()

	fmt.Println("Controlador: Iniciando simulação...")
	time.Sleep(1 * time.Second)

	// 1. Definir líder inicial (processo 4 como 5)
	fmt.Println("\nControlador: Definindo processo 4 como líder inicial (5)")
	broadcast(mensagem{tipo: 1, corpo: []int{4}, message: "Novo líder definido: 4"})
	time.Sleep(2 * time.Second)

	// 2. Simular falha do processo 4 (5)
	fmt.Println("\nControlador: Simulando falha do processo 4 (5)")
	chans[4] <- mensagem{tipo: 2, message: "Eleição (falha 5)"}
	waitForElection()

	// 3. Simular falha do processo 3 (4)
	fmt.Println("\nControlador: Simulando falha do processo 3 (4)")
	chans[3] <- mensagem{tipo: 2, message: "Eleição (falha 4)"}
	waitForElection()

	// 4. Simular retornos
	fmt.Println("\nControlador: Simulando retorno do processo 3 (4)")
	chans[3] <- mensagem{tipo: 3, message: "Processo 4 recuperado"}
	time.Sleep(1 * time.Second)

	fmt.Println("\nControlador: Simulando retorno do processo 4 (5)")
	chans[4] <- mensagem{tipo: 3, message: "Processo 5 recuperado"}
	time.Sleep(1 * time.Second)

	// 5. Encerrar
	fmt.Println("\nControlador: Encerrando todos os processos")
	close(done)
	broadcast(mensagem{tipo: 9, message: "Encerramento"})
	fmt.Println("\nControlador: Simulação concluída")
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	var (
		actualLeader = leader
		bFailed      = false
		next         = (TaskId + 1) % 5 // Próximo no anel
	)

	fmt.Printf("%d: Iniciando (líder inicial: %d)\n", TaskId, actualLeader)

	for {
		select {
		case msg := <-in:
			if !processMessage(TaskId, msg, &actualLeader, &bFailed, out, next) {
				return
			}
		case <-done:
			fmt.Printf("%d: Encerrando\n", TaskId)
			return
		}
	}
}

func processMessage(TaskId int, msg mensagem, leader *int, failed *bool, out chan mensagem, next int) bool {
	if msg.tipo == 9 {
		fmt.Printf("%d: Encerrando\n", TaskId)
		return false
	}

	fmt.Printf("%d: Recebido: %s\n", TaskId, msg.message)

	switch msg.tipo {
	case 0: // Eleição
		if *failed {
			fmt.Printf("%d: Ignorando (falho)\n", TaskId)
			return true
		}

		if msg.origem == TaskId {
			newLeader := max(msg.corpo)
			*leader = newLeader
			fmt.Printf("%d: Eleição concluída. Novo líder: %d\n", TaskId, newLeader)
			out <- mensagem{
				tipo:    1,
				corpo:   []int{newLeader},
				sender:  TaskId,
				message: fmt.Sprintf("Eleição (volta %d)", newLeader),
			}
			if TaskId == 0 {
				controle <- fmt.Sprintf("Novo líder eleito: %d", newLeader)
			}
		} else {
			if !contains(msg.corpo, TaskId) {
				msg.corpo = append(msg.corpo, TaskId)
			}
			fmt.Printf("%d: Encaminhando para %d\n", TaskId, next)
			out <- msg
		}

	case 1: // Novo líder
		if *failed {
			fmt.Printf("%d: Ignorando (falho)\n", TaskId)
			return true
		}
		*leader = msg.corpo[0]
		if msg.sender != TaskId {
			out <- msg
		}

	case 2: // Falha
		*failed = true
		if TaskId == *leader {
			fmt.Printf("%d: Era líder, iniciando eleição\n", TaskId)
			out <- mensagem{
				tipo:    0,
				corpo:   []int{TaskId},
				sender:  TaskId,
				origem:  TaskId,
				message: fmt.Sprintf("Eleição (iniciada por %d)", TaskId),
			}
		}

	case 3: // Recuperação
		*failed = false
	}

	return true
}

// Funções auxiliares
func broadcast(msg mensagem) {
	for i := 0; i < 5; i++ {
		chans[i] <- msg
	}
}

func waitForElection() {
	fmt.Println("Controlador: Aguardando eleição...")
	select {
	case log := <-controle:
		fmt.Printf("Controlador: %s\n", log)
	case <-time.After(3 * time.Second):
		fmt.Println("Controlador: Timeout - continuando")
	}
}

func max(ids []int) int {
	m := ids[0]
	for _, id := range ids[1:] {
		if id > m {
			m = id
		}
	}
	return m
}

func contains(ids []int, id int) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}