# Payment Gateway

## 1. Membros do Grupo
1. Bruna Dias
2. Leonardo Borges
3. Victor Augusto

## 2. Explicação do Sistema
O Payment Gateway é um sistema de processamento de pagamentos que permite a realização segura de transações financeiras. Desenvolvido como parte de um projeto acadêmico, o sistema oferece uma API RESTful para gerenciar operações de pagamento, incluindo processamento de transações e consulta de status.

Principais funcionalidades:
- Criação de pagamentosç
- Processamento de pagamentos;
- Consulta do status dos pedidos;

Para que o Microserviço não fosse tão grande, a parte de criação de um pedido foi abstraída e tomado como pré-requisito para utilização do sistema de pagamento. 
## 3. Tecnologias Utilizadas
- **Linguagem de Programação**: Go (Golang)
  - Escolhida por sua performance e suporte nativo a concorrência

- **Framework Web**: Gin
  - Framework HTTP web rápido e produtivo para Go
  - Roteamento eficiente e middleware flexível
   
- **Banco de Dados**: MySQL
  - Base de dados relacional
   
- **Conteinerização**: Docker
  - Utilização de Docker-compose para gerenciamento de ambientes

- **Testes**: Testes unitários nativos do Go
  - Suporte a testes concorrentes
  - Ferramentas de benchmark integradas

- **Gerenciamento de Dependências**: Go Modules
  - Gerenciamento de versões de pacotes
  - Reprodução de builds consistente

- **Controle de Versão**: Git
  - Rastreamento de mudanças no código-fonte
  - Facilita o trabalho em equipe

- **Documentação**: Markdown
  - Formato leve para documentação
  - Facilidade de leitura e manutenção

## 4. Limitações do GitHub Actions
Não foi possível utilizar todos os ambientes para os testes E2E, por conta da utilização do docker-compose para gerenciar os ambientes. Como se trata de uma ferramenta baseada em linux, nos ambientes Windows e MacOS ocorreu uma incompatibilidade com o docker-compose, como é possível observar nas actions prévias, que tornou inviável a utilização dos demais ambientes nos testes E2E. Nesse sentido, os testes unitários estão sendo executados nos três ambientes e os testes E2E apenas no Linux.
