# Ponderada Rinha de Backend

## Descrição

API RESTful desenvolvida em Go com banco de dados PostgreSQL e com Docker, com objetivo de simular movimentações financeiras para múltiplos usuários simultaneamente. O sistema foi projetado para suportar alta carga de requisições, com balanceamento de carga, concorrência segura e estrutura escalável.

[Rinha de Backend 2024 Q1](https://github.com/zanfranceschi/rinha-de-backend-2024-q1)

## Funcionalidades Implementadas

- **POST `/clientes/{id}/transacoes`**
  - Realiza transações de crédito ou débito
  - Validação de saldo e limite
  - Resposta com saldo atualizado

- **GET `/clientes/{id}/extrato`**
  - Exibe as últimas 10 transações
  - Inclui saldo atual, limite e data da consulta

## Tecnologias Utilizadas

### Backend

- **Go** – linguagem leve, compilada e altamente performática
- `net/http` – servidor HTTP nativo da linguagem
- `chi` – roteador minimalista e eficiente, ideal para APIs REST
- `pgx/v5` – driver PostgreSQL com suporte avançado a pool de conexões
- `Context API` – gerenciamento de timeout e controle transacional
- `Transações SQL (**FOR UPDATE**)` – controle de concorrência no acesso ao saldo

### Banco de Dados (PostgreSQL)

- **PostgreSQL**
- `Schema relacional` – tabelas *clientes* e *transacoes* com constraints
- `Integridade garantida` – uso de foreign keys e validações no banco
- `Comandos SQL otimizados` – *INSERT*, *SELECT FOR UPDATE*, *UPDATE* com foco em concorrência segura

### DevOps & Testes

- **Docker & Docker Compose**
- `NGINX` – balanceador de carga configurado com `least_conn`
- `k6` – ferramenta para testes de carga, performance e stress testing
- `Pool de conexões` – até 100 conexões simultâneas otimizadas para PostgreSQL
- `Escalabilidade horizontal` – duas instâncias da API Go balanceadas por NGINX

## Resposta das Perguntas Enunciado

### O que você fez para garantir a segurança do sistema?
Como o desafio da Rinha não exigia autenticação, foquei em validar bem o que chega na API. Todo payload é checado: se o tipo da transação é válido, se a descrição tem o tamanho correto e se o valor faz sentido. Também isolei as operações por cliente.

### O que você fez para garantir a integridade dos dados?
Usei transações com `SELECT` `FOR UPDATE` no PostgreSQL para evitar condições de corrida ao atualizar saldo. Isso garante que duas transações simultâneas não causem erro no saldo final. Também defini constraints no banco (como tipo ser só 'c' ou 'd') e foreign keys para garantir que cada transação esteja ligada a um cliente válido.

### O que você fez para garantir a disponibilidade do sistema?
Configurei o NGINX para fazer balanceamento de carga entre duas instâncias da API. Isso já melhora bastante a disponibilidade. Além disso, com o Docker Compose, é fácil subir tudo novamente se algo falhar. A ideia é que o sistema sempre tenha ao menos uma instância ativa respondendo.

### O que você fez para garantir a escalabilidade do sistema?
A escalabilidade foi feita horizontalmente: duas instâncias da API Go, balanceadas pelo NGINX com `least_conn`, que distribui melhor quando há picos. Também deixei o banco com `max_connections` e configurei o pool de conexões da aplicação com `MaxConns`, o que ajuda bastante quando temos muitos usuários simultâneos.

### O que você fez para garantir a performance do sistema?
Go já é rápido por padrão, mas mesmo assim evitei usar bibliotecas pesadas. Os acesso ao banco é direto com SQL, configurei um pool de conexões para evitar abrir uma nova a cada requisição.

### O que você fez para garantir a manutenibilidade do sistema?
Organizei o projeto em pastas separadas por responsabilidade, o que ajuda muito a manter e entender o sistema. Por exemplo:
- A pasta `api/` contém os handlers da API REST — é onde estão as rotas e a lógica para transações e extratos.
- A pasta `database/` cuida da conexão com o PostgreSQL e do script de inicialização do banco (`init.sql`).
- Em `k6-test/`, deixei os testes de carga escritos com `k6`, facilitando a validação de performance sempre que algo for alterado.
- O main.go serve só para orquestrar tudo: carrega a conexão, monta as rotas e sobe o servidor.

### O que você fez para garantir a testabilidade do sistema?
Comecei focando no teste principal da atividade, que é o de carga. Criei um script com o k6, que simula usuários reais fazendo transações e consultas de extrato, em diferentes volumes. Esse teste mede tempo de resposta, taxa de erros e quantidade de requisições válidas.

Pra facilitar isso, organizei o código da API de forma que cada parte tem sua responsabilidade. Por exemplo, o arquivo handlers.go cuida só das rotas e da lógica da API, e a conexão com o banco está separada em `database/connection.go`. Isso ajuda muito caso eu quisesse testar as funções de forma isolada no futuro, como com testes unitários ou mocks.
