package main

import (
    "fmt"
    "sync"
)

// Estrutura de mensagem usada para comunicação entre os processos
type mensagem struct {
    tipo  int    // Tipo da mensagem (ex.: 2 para falha, 3 para recuperação, etc.)
    corpo [3]int // Conteúdo da mensagem (ex.: IDs dos processos envolvidos)
}

const (
    MSG_FAILURE     = 2 // Mensagem indicando falha
    MSG_RECOVERY    = 3 // Mensagem indicando recuperação
    MSG_ELECTION    = 4 // Mensagem indicando início de eleição
    MSG_COORDINATOR = 5 // Mensagem indicando novo coordenador
)

var (
    // Vetor de canais para formar o anel de eleição
    chans = []chan mensagem{
        make(chan mensagem), // Canal do processo 0
        make(chan mensagem), // Canal do processo 1
        make(chan mensagem), // Canal do processo 2
        make(chan mensagem), // Canal do processo 3
        make(chan mensagem), // Canal do processo 4
        make(chan mensagem), // Canal do processo 5
    }
    controle = make(chan int)       // Canal de controle para comunicação com o controlador
    monitor  = make(chan mensagem)  // Canal para comunicação com o processo 0 (monitoramento)
    wg       sync.WaitGroup         // WaitGroup para sincronizar a execução das goroutines
)

// Processo 0 (Monitoramento)
func MonitorProcess() {
    defer wg.Done()

    for {
        msg := <-monitor
        switch msg.tipo {
        case MSG_FAILURE:
            fmt.Printf("Monitor: Processo %d falhou, detectado por %d\n", msg.corpo[1], msg.corpo[0])
        case MSG_RECOVERY:
            fmt.Printf("Monitor: Processo %d recuperado\n", msg.corpo[1])
        case MSG_ELECTION:
            fmt.Printf("Monitor: Eleição iniciada por %d\n", msg.corpo[0])
        case MSG_COORDINATOR:
            fmt.Printf("Monitor: Novo coordenador é o processo %d\n", msg.corpo[0])
        default:
            fmt.Printf("Monitor: Mensagem desconhecida recebida\n")
        }
    }
}

// Função que representa cada estágio do anel de eleição
func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
    defer wg.Done()

    // Variáveis locais que indicam se este processo é o líder e se está ativo
    var actualLeader int = leader
    var bFailed bool = false // Todos iniciam sem falha

    for {
        temp := <-in
        switch temp.tipo {
        case MSG_FAILURE: // Tipo 2: mensagem indicando falha
            bFailed = true
            fmt.Printf("%2d: falho %v \n", TaskId, bFailed)
            monitor <- mensagem{tipo: MSG_FAILURE, corpo: [3]int{TaskId, TaskId, actualLeader}}
        case MSG_RECOVERY: // Tipo 3: mensagem indicando recuperação
            bFailed = false
            fmt.Printf("%2d: falho %v \n", TaskId, bFailed)
            monitor <- mensagem{tipo: MSG_RECOVERY, corpo: [3]int{TaskId, TaskId, actualLeader}}
        case MSG_ELECTION: // Tipo 4: mensagem indicando início de eleição
            if !bFailed {
                if temp.corpo[0] == TaskId {
                    newLeader := temp.corpo[2]
                    monitor <- mensagem{tipo: MSG_COORDINATOR, corpo: [3]int{newLeader, newLeader, newLeader}}
                } else {
                    if TaskId > temp.corpo[2] {
                        temp.corpo[2] = TaskId
                    }
                    out <- temp
                }
            }
        default: // Tipo desconhecido
            fmt.Printf("%2d: não conheço este tipo de mensagem\n", TaskId)
        }
    }
}

// Função principal
func main() {
    wg.Add(6) // Adiciona uma contagem de 6 ao WaitGroup (5 processos + 1 monitor)

    // Criar o processo 0 (monitoramento)
    go MonitorProcess()

    // Criar os processos do anel de eleição
    go ElectionStage(1, chans[0], chans[1], 1)
    go ElectionStage(2, chans[1], chans[2], 1)
    go ElectionStage(3, chans[2], chans[3], 1)
    go ElectionStage(4, chans[3], chans[4], 1)
    go ElectionStage(5, chans[4], chans[0], 1)

    fmt.Println("\n   Anel de processos criado")

    // Simular falhas e recuperações
    go func() {
        // Simular falha no processo 2
        fmt.Println("Simulação: Processo 2 falhou")
        chans[1] <- mensagem{tipo: MSG_FAILURE, corpo: [3]int{0, 2, 0}}

        // Simular recuperação do processo 2
        fmt.Println("Simulação: Processo 2 recuperado")
        chans[1] <- mensagem{tipo: MSG_RECOVERY, corpo: [3]int{0, 2, 0}}

        // Iniciar eleição
        fmt.Println("Simulação: Iniciando eleição")
        chans[1] <- mensagem{tipo: MSG_ELECTION, corpo: [3]int{1, 0, 1}}
    }()

    wg.Wait() // Aguarda todas as goroutines terminarem
    fmt.Println("Execução concluída")
}