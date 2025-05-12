# Enunciado T1 - FPPD

Implementar o algoritmo distribuído de eleição (de um processo coordenador) baseado em anel lógico (Ring Algorithm) utilizando o código base fornecido na linguagem Go :

Baseado no uso de um anel lógico;
Cada processo conhece o anel inteiro, mas manda mensagens somente para o próximo processo ativo no sentido do anel;
Quando processo detecta que o coordenador não está mais ativo (isto pode ser simulado com uma mensagem de um processo externo), ele envia uma mensagem de eleição no anel contendo seu id (disputando esta coordenação que está vaga);
A cada passo, o processo que recebe a mensagem de eleição inclui seu id na mensagem e passa adiante no anel
No final de uma volta completa, o processo que iniciou a eleição recebe a mensagem e escolhe aquele que tem maior id como novo processo coordenador;
Uma nova mensagem é enviada através do anel para que todos conheçam o novo coordenador;
Modelar uma forma de mostrar que está funcionando (ex: processo externo que simula solicitações) e também de simular falhas para que a resposta não seja sempre a mesma (nodos do anel que param de funcionar).

A avaliação do trabalho será feita com base no acompanhamento do desenvolvimento em laboratório e da entrega do relatório, nos mesmos moldes do trabalho anterior.

