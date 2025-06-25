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

### O que você fez para garantir a testabilidade do sistema?
