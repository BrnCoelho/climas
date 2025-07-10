Projeto Prático - Painel de Clima em Tempo Real
🎯 Objetivo do Projeto
Este projeto tem como objetivo a criação de uma aplicação funcional que integra Rust (coleta de dados), Go (orquestração/API) e Dart (interface Flutter), utilizando APIs públicas e dados em tempo real ou quase real. A aplicação exibe o clima atual e a previsão do tempo para diversas cidades, com uma arquitetura modular e comunicação via API RESTful.

🧩 Estrutura do Projeto
coleta_rust/: Módulo em Rust responsável por coletar dados climáticos de APIs externas e alimentar nosso banco de dados através da API Go.

api_go/: Servidor RESTful construído em Go. Contém toda a lógica de negócio e os endpoints, utilizando gorilla/mux para roteamento e pq para a comunicação com o banco de dados.

Estrutura Interna:

Database: Rodando em um contêiner Docker.

main.go: Ponto de entrada que inicializa o servidor e registra as rotas.

repository_clima.go: Camada de acesso a dados para a tabela clima_atual (ex: buscar por ID).

repository_previsao.go: Camada de acesso a dados para a tabela previsao_dias (ex: buscar por id_clima).

types.go: Define as estruturas de dados (structs) utilizadas no projeto.

clima_app/: Aplicação cliente multiplataforma desenvolvida com Flutter.

Estrutura Interna (lib/):

api/: Camada de serviço responsável pela comunicação com a API Go.

models/: Classes que modelam os dados da aplicação (ClimaAtual, PrevisaoDia).

screens/: Widgets que representam as telas do app (HomeScreen, DetailScreen).

utils/: Funções auxiliares, como o conversor de condição climática para ícones.

docs/: Documentos e apresentações do projeto.

docker-compose.yml: Arquivo de configuração para iniciar o ambiente do banco de dados PostgreSQL.

🚀 Como Executar
Para executar o projeto completo, siga os passos na ordem:

1. Ambiente e Banco de Dados
Certifique-se de ter o Docker e o Docker Compose instalados. Na pasta raiz do projeto, inicie o contêiner do banco de dados com o comando:

docker-compose up -d
2. Backend (API Go)
Abra um novo terminal, navegue até a pasta api_go/ e inicie o servidor. O servidor estará rodando na porta 8080.

go run .
3. Coletor (Rust)
Abra outro terminal, navegue até a pasta coleta_rust/ e execute o script de coleta.

cargo run
4. Frontend (Flutter)
Abra a pasta clima_app/ no seu editor de código (VS Code).

Importante: No arquivo lib/api/api_service.dart, verifique e atualize a variável baseUrl com o endereço IP correto do seu servidor Go.

Abra um terminal integrado e execute a aplicação para a plataforma desejada:

flutter run
⚖ Uso de IA
Este projeto não utiliza modelos de Inteligência Artificial.

🧪 Testes Automatizados e CI/CD
O ambiente de desenvolvimento foi padronizado utilizando Docker e Docker Compose para o banco de dados PostgreSQL, garantindo consistência na execução do backend em diferentes máquinas. O projeto Flutter inclui os testes de widget padrão (test/widget_test.dart) gerados pelo framework.



Módulo Rust - coleta_rust
Este módulo é responsável por consumir dados da API HG Brasil (clima), processar e enviar os dados atualizados para o backend em Go por meio de requisições HTTP (POST e PUT).
Funcionalidades:
1-Realiza chamadas GET para obter o clima atual e previsões dos 6 dias.
2-Envia os dados do clima atual via POST para a API Go e armazena o ID retornado.
3-Envia cada previsão via POST, vinculando-a ao clima atual pelo ID.
4-Armazena os IDs retornados (clima e previsões) para atualizações futuras.
5-A cada X segundos(configurado pelo código), realiza novas requisições para atualizar os dados (PUT), utilizando os IDs já armazenados.
6-Os dados enviados seguem o formato esperado pela API Go (em JSON).

Estrutura:
main.rs: Código principal que executa o loop de coleta e atualização.
cargo.toml: configuração das dependencies

Utiliza:
-reqwest para requisições HTTP.
-tokio para programação assíncrona.
-serde para serialização e deserialização de JSON.
-std::collections::HashMap para controle local de IDs.

Observação:
Para ambientes onde os dados já existem no banco de dados, o código permite comentar o trecho de POST e definir manualmente os IDs para continuar apenas com as atualizações (PUT). Isso evita duplicação de dados em bancos pequenos ou de teste.
