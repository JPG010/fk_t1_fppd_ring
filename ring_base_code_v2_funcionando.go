
package main

import (
	"fmt"
	"sync"
)

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo [5]int // conteudo da mensagem para colocar os ids (usar um tamanho ocmpativel com o numero de processos no anel)
}

var (
	chans = []chan mensagem{ // vetor de canias para formar o anel de eleicao - chan[0], chan[1] and chan[2] ...
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
	}
	controle = make(chan int)
	wg       sync.WaitGroup // wg is used to wait for the program to finish
)

func ElectionControler(in chan int) {
	defer wg.Done()

	var msg mensagem

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Printf("PROCESSO CONTROLE: mudar o processo 4 para falho\n")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
	msg.tipo = 2
	msg.corpo = [5]int{-1, -1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar pra nao substituir o lider atual
	chans[4] <- msg
	fmt.Printf("\nPROCESSO CONTROLE: confirmacao = %d\n", <-in) // confirmacao do canal controle

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Printf("PROCESSO CONTROLE: solicitar ao processo 0 para iniciar eleicao pois detectou o processo 4 como falho\n") 
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
	msg.tipo = 4
	msg.corpo = [5]int{4, 4, 4, 4, 4} // mandar qual falhou
	chans[0] <- msg
	fmt.Printf("\nPROCESSO CONTROLE: confirmacao = %d\n", <-in) // confirmacao do canal controle

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Printf("PROCESSO CONTROLE: mudar o processo 3 para falho\n")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
	msg.tipo = 2
	msg.corpo = [5]int{-1, -1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar pra nao substituir o lider atual
	chans[3] <- msg
	fmt.Printf("\nPROCESSO CONTROLE: confirmacao = %d\n", <-in) // confirmacao do canal controle

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Printf("PROCESSO CONTROLE: solicitar ao processo 0 para iniciar eleicao pois detectou o processo 3 como falho\n")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
	msg.tipo = 4
	msg.corpo = [5]int{3, 3, 3, 3, 3} // mandar qual falhou
	chans[0] <- msg
	fmt.Printf("\nPROCESSO CONTROLE: confirmacao = %d\n", <-in) // confirmacao do canal controle
	
	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Printf("PROCESSO CONTROLE: reativar o processo 4\n") 
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
	msg.tipo = 3
	msg.corpo = [5]int{-1, -1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar pra nao substituir o lider atual
	chans[3] <- msg
	fmt.Printf("\nPROCESSO CONTROLE: confirmacao = %d\n", <-in) // confirmacao do canal controle

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Printf("PROCESSO CONTROLE: solicitar ao processo 3 para iniciar eleicao pois detectou o processo 4 como reativado\n") 
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
	msg.tipo = 4
	msg.corpo = [5]int{-1, -1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar
	chans[2] <- msg
	fmt.Printf("\nPROCESSO CONTROLE: confirmacao = %d\n", <-in) // confirmacao do canal controle

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Printf("PROCESSO CONTROLE: reativar o processo 5\n") 
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
	msg.tipo = 3
	msg.corpo = [5]int{-1, -1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar pra nao substituir o lider atual
	chans[4] <- msg
	fmt.Printf("\nPROCESSO CONTROLE: confirmacao = %d\n", <-in) // confirmacao do canal controle

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Printf("PROCESSO CONTROLE: solicitar ao processo 4 para iniciar eleicao pois detectou o processo 5 como reativado\n") 
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
	msg.tipo = 4
	msg.corpo = [5]int{-1, -1, -1, -1, -1} //colocar "lixo" no corpo da mensagem para enviar
	chans[3] <- msg
	fmt.Printf("\nPROCESSO CONTROLE: confirmacao = %d\n", <-in) // confirmacao do canal controle

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Println("PROCESSO CONTROLE: matar todos os processos")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
	msg.tipo = 5

	for i, c := range chans {
		c <- msg
		fmt.Printf("PROCESSO CONTROLE: processo %2d foi morto, confirmacao = %d\n", i, <-in)
	}
	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Println("PROCESSO CONTROLE FINALIZADO")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")

}
func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	var actualLeader int = leader
	var bFailed bool = false // todos inciam sem falha

	var stopRing bool = false // variável para parar o processo quando chegar no sinal de término de processo

	for !stopRing {	

		temp := <-in // ler mensagem
		fmt.Printf("\nELEIÇÃO: %2d: recebi mensagem %d, [ %d, %d, %d, %d, %d ]\n", 
					TaskId, temp.tipo, 
					temp.corpo[0], 
					temp.corpo[1], 
					temp.corpo[2], 
					temp.corpo[3], 
					temp.corpo[4])

		switch temp.tipo {

		case 0: // confirmação de final de eleição, indica o novo líder elegido
			{
			if bFailed { //nao faz nada se o processo já estiver falho
				fmt.Printf("%2d: processo falho, votação não realizada\n", TaskId)					
			} else { 
				actualLeader = temp.corpo[0] // atualiza o lider a partir da mensagem que chegou até ali
				fmt.Printf("%2d: processo eleito como líder é: %d\n", TaskId, actualLeader)
			}
			out <- temp // envia a mensagem independente se o processo está falhado
			// nao precisa de controle pq nao sai do anel pra chegar no controle
			}
		case 1: // controle indica para este processo que um processo falhou = faz uma requisicao para iniciar o processo de eleição
			{
				if bFailed { // processo falhado, nao faz nada e passa pra frente no anel
					fmt.Printf("%2d: processo falho, não participa da eleição\n", TaskId)
				} else {
					temp.corpo[TaskId] = TaskId
					fmt.Printf("%2d: processo votou\n", TaskId)
				} 
				out <- temp //passa a mensagem
				// nao precisa de controle pq nao sai do anel pra chegar no controle
			}
		case 2: // falha manualmente um processo
			{
				bFailed = true
				fmt.Printf("%2d: processo falho %v \n", TaskId, bFailed)
				controle <- 1
			}
		case 3: // reativa manualmente um processo
			{
				bFailed = false
				fmt.Printf("%2d: processo reativado %v \n", TaskId, bFailed)
				controle <- 1
			}
		case 4: // processo de eleicao (recebe mensagem para iniciar a eleicao)
			{	
				if !bFailed { //processo não falho, então inicia a eleição
					fmt.Printf("\n%2d: falha do processo %2d identificada\n", TaskId, temp.corpo[0])
					fmt.Printf("%2d: iniciando eleição\n", TaskId)
					
					//eleição
					msg := mensagem{tipo: 1, corpo: [5]int{-1, -1, -1, -1, -1}}
					msg.corpo[TaskId] = TaskId
					out <- msg
					fmt.Printf("\n%2d: processo participou da votação da eleição", TaskId)

					result := <-in
					if result.tipo != 1 {
						fmt.Printf("%2d: mensagem não esperada\n", TaskId)
						controle <- 0
						return
					} 
					//recebendo os votos, determina o vencedor 
					fmt.Printf("\n%2d: votos recebidos [ %d, %d, %d, %d, %d ]\n", TaskId, result.corpo[0], result.corpo[1], result.corpo[2], result.corpo[3], result.corpo[4])
					var winner int = result.corpo[0]
					for i:=1; i<len(result.corpo); i++ {
						if result.corpo[i] > winner {
							winner = result.corpo[i]
						}
					}
					confirm := mensagem {tipo: 0, corpo: [5]int{winner, winner, winner, winner, winner}}
					actualLeader = winner
					fmt.Printf("%2d: novo processo eleito para líder: %d\n", TaskId, actualLeader)
					fmt.Printf("%2d: mensagem de confirmação de novo líder sendo enviada para o anel\n", TaskId)
					out <- confirm

					finalResult := <-in
					if finalResult.tipo != 0 {
						fmt.Printf("%2d: mensagem não esperada\n", TaskId)
						controle <- 0
						return
					} 
					controle <- 1

				} else {
					fmt.Printf("%2d: processo falho, eleição não iniciada\n", TaskId)
					controle <- 0 //processo falho, entao nao tem como iniciar a eleicao, retorna pro controle um sinal de erro
				}
			}

		case 5: // finaliza os processos
			{
				fmt.Println("Terminação de processo")
				controle <- 1
				stopRing = true // matar o processo
			}
		default:
			{
				fmt.Printf("%2d: não conheço este tipo de mensagem\n", TaskId)
				fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
			}
		}
	}

	fmt.Printf("Processo %2d terminado\n", TaskId)
}	

func main() {

	wg.Add(6) 

	go ElectionStage(0, chans[0], chans[1], 4) 
	go ElectionStage(1, chans[1], chans[2], 4) 
	go ElectionStage(2, chans[2], chans[3], 4) 
	go ElectionStage(3, chans[3], chans[4], 4) 
	go ElectionStage(4, chans[4], chans[0], 4) 

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Println("ANEL DE PROCESSOS CRIADO")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")

	go ElectionControler(controle)

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Println("PROCESSO CONTROLE CRIADO")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")

	wg.Wait() 

	fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - ")
	fmt.Println("PROCESSOS ENCERRADOS")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - \n")
}
