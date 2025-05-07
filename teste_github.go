package main

import (
	"fmt"
	"sync"
)

/* TIPOS DE MENSAGENS NO CANAL INT (tipo)
	0 = election confirmation 
	1 = vote request 
	2 = falho 
	3 = ativo 
	4 = eleicao
	5 = finaliza
*/

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo [4]int // conteudo da mensagem para colocar os ids (usar um tamanho ocmpativel com o numero de processos no anel)
}

var (
	chans = []chan mensagem{ // vetor de canias para formar o anel de eleicao - chan[0], chan[1] and chan[2] ...
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
	} //como é um vetor formando o anel, basta fazer um for que percorre todos, criando a logica de anel (algoritmo)
	controle = make(chan int)
	wg       sync.WaitGroup // wg is used to wait for the program to finish
)

func ElectionControler(in chan int) {
	defer wg.Done()
	fmt.Println("Controler: iniciando rotina controle de eleição")

	var msg mensagem
	
	fmt.Printf("Controler: mudar o processo 2 para falho\n")
	msg.tipo = 2
	msg.corpo = [4]int{-1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar pra nao substituir o lider atual
	chans[2] <- msg
	fmt.Printf("Controler: confirmacao de resposta da ElectionStage = %d\n", <-in) // confirmacao do canal controle
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")

	fmt.Printf("Controler: solicitar ao processo 3 para iniciar eleicao pois detectou o processo 2 como falho\n") 
	msg.tipo = 4
	msg.corpo = [4]int{2, 2, 2, 2} // mandar qual falhou
	chans[3] <- msg
	fmt.Printf("Controler: confirmacao de resposta da ElectionStage = %d\n", <-in) // confirmacao do canal controle
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")

	fmt.Printf("Controler: mudar o processo 3 para falho\n")
	msg.tipo = 2
	msg.corpo = [4]int{-1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar pra nao substituir o lider atual
	chans[3] <- msg
	fmt.Printf("Controler: confirmacao de resposta da ElectionStage = %d\n", <-in) // confirmacao do canal controle
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")

	fmt.Printf("Controler: solicitar ao processo 0 para iniciar eleicao pois detectou o processo 3 como falho\n")
	msg.tipo = 4
	msg.corpo = [4]int{3, 3, 3, 3} // mandar qual falhou
	chans[0] <- msg
	fmt.Printf("Controler: confirmacao de resposta da ElectionStage = %d\n", <-in) // confirmacao do canal controle
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")

	fmt.Printf("Controler: reativar o processo 2\n") 
	msg.tipo = 3
	msg.corpo = [4]int{-1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar pra nao substituir o lider atual
	chans[2] <- msg
	fmt.Printf("Controler: confirmacao de resposta da ElectionStage %d\n", <-in) // confirmacao do canal controle
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")

	fmt.Printf("Controler: solicitar ao processo 1 para iniciar eleicao pois detectou o processo 2 como reativado\n") 
	msg.tipo = 4
	msg.corpo = [4]int{-1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar
	chans[1] <- msg
	fmt.Printf("Controler: confirmacao de resposta da ElectionStage %d\n", <-in) // confirmacao do canal controle
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")


	// matar todos os processos do anel iterando pelo vetor de canais
	fmt.Println("\nControler: matar todos os processos")
	msg.tipo = 5
	for i, c := range chans {
		c <- msg
		fmt.Printf("Controler: processo %2d foi morto, confirmacao = %d", i, <-in)
	}

	fmt.Println("\nControler: processo controlador concluído ")
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	// variaveis locais que indicam se este processo é o lider e se esta ativo
	var actualLeader int = leader // indicação do lider veio por parâmatro
	var bFailed bool = false // todos inciam sem falha

	var pare bool = false // variável para parar o processo quando chegar no sinal de término de processo
	for !pare {
		temp := <-in // ler mensagem do canal de entrada do anel
		fmt.Printf("\n\n\tElection: %2d: recebi mensagem %d, [ %d, %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3])

		switch temp.tipo {
		case 0: // confirmação de final de eleição, indica o novo líder elegido
			{
				if bFailed { //nao faz nada se o processo já estiver falho
					fmt.Printf("\t%2d: processo falho, líder não votado\n", TaskId)					
				} else { 
					actualLeader = temp.corpo[0] // atualiza o lider a partir da mensagem que chegou até ali
					fmt.Printf("\t%2d: processo eleito como NOVO LÍDER é: %d\n", TaskId, actualLeader)
				}
				out <- temp // envia a mensagem independente se o processo está falhado
				// nao precisa de controle pq nao sai do anel pra chegar no controle
			}
		case 1: // controle indica para este processo que um processo falhou = faz uma requisicao para iniciar o processo de eleição
			{
				if bFailed { // processo falhado, nao faz nada e passa pra frente no anel
					fmt.Printf("\t%2d: nao participou da votacao, processo está falho\n", TaskId)
				} else {
					temp.corpo[TaskId] = TaskId
					fmt.Printf("\t%2d: processo votou\n", TaskId)
				} 
				out <- temp //passa a mensagem
				// nao precisa de controle pq nao sai do anel pra chegar no controle

			}
		case 2: // falha manualmente um processo
			{
				bFailed = true
				fmt.Printf("\t%2d: processo falho %v \n", TaskId, bFailed)
				controle <- 1
			}
		case 3: // reativa manualmente um processo
			{
				bFailed = false
				fmt.Printf("\t%2d: processo retirado da falha (reativado) %v \n", TaskId, bFailed)
				controle <- 1
			}
		case 4: // processo de eleicao (recebe mensagem para iniciar a eleicao)
			{	
				if !bFailed { //processo não falho, então inicia a eleição
					fmt.Printf("\t%2d: processo %2d identificado como falho\n", TaskId, temp.corpo[0])
					fmt.Printf("\t%2d: iniciando eleicao\n", TaskId)
					
					//eleição
					msg := mensagem{tipo: 1, corpo: [4]int{-1, -1, -1, -1}}
					msg.corpo[TaskId] = TaskId
					out <- msg
					fmt.Printf("\t%2d: processo participou da votacao", TaskId)


					result := <-in
					if result.tipo != 1 {
						fmt.Printf("\t%2d: mensagem inesperada\n", TaskId)
						controle <- 0
						return
					} 

					//recebendo os votos, determina o vencedor 
					fmt.Printf("\n\t%2d: recebi os votos [ %d, %d, %d, %d ]\n", TaskId, result.corpo[0], result.corpo[1], result.corpo[2], result.corpo[3])
					var winner int = result.corpo[0]
					for i:=1; i<len(result.corpo); i++ {
						if result.corpo[i] > winner {
							winner = result.corpo[i]
						}
					}

					confirm := mensagem {tipo: 0, corpo: [4]int{winner, winner, winner, winner}}
					actualLeader = winner
					fmt.Printf("\t%2d: novo processo eleito para lider: %d\n", TaskId, actualLeader)
					fmt.Printf("\t%2d: enviando mensagem de confirmacao de lider para o anel, indicando o novo líder\n", TaskId)
					out <- confirm

					finalResult := <-in
					if finalResult.tipo != 0 {
						fmt.Printf("\t%2d: mensagem inesperada\n", TaskId)
						controle <- 0
						return
					} 
					controle <- 1

				} else {
					fmt.Printf("\t%2d: processo falho, eleicao nao iniciada\n", TaskId)
					controle <- 0 //processo falho, entao nao tem como iniciar a eleicao, retorna pro controle um sinal de erro
				}
				
			}
		case 5: // finaliza os processos
			{
				fmt.Println("\tChegou em processo com mensagem de termino (codigo 5)")
				controle <- 1
				pare = true // matar o processo
			}
		default:
			{
				fmt.Printf("\t%2d: não conheço este tipo de mensagem\n", TaskId)
			}
		}
	}
	
	fmt.Printf("\t%2d: terminei \n", TaskId)
}

func main() {

	wg.Add(5) // adiciona uma contagem para cada goroutine

	// cria os processo do anel de eleicao
	go ElectionStage(0, chans[0], chans[1], 3) // este é o lider
	go ElectionStage(1, chans[1], chans[2], 3) // não é lider, é o processo 0
	go ElectionStage(2, chans[2], chans[3], 3) // não é lider, é o processo 0
	go ElectionStage(3, chans[3], chans[0], 3) // não é lider, é o processo 0

	fmt.Println("\n		Main: Anel de processos criado")

	// cria o processo controlador
	go ElectionControler(controle)
	fmt.Println("\n		Main: Processo controlador criado")

	wg.Wait() // espera as goroutines para terminar
	fmt.Println("\n		Main: programa encerrado")
}