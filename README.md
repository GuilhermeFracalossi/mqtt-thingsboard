# Simulação de sensores + Broker MQTT ThingsBoard

## Como instalar
1. Dentro do WSL ou do terminal do Linux, instale o Docker e o Docker Compose.
2. Instale o Go na versão mais recente.
3. Clone o repositório.
4. Mova o dataset para a pasta `data` na raiz do projeto.
5. Copie o arquivo `.sample.env` para `.env`.
6. Atualize o arquivo `.env` se necessário, com o caminho do dataset e os tokens dos dispositivos do ThingsBoard.
7. Execute o comando `docker compose build` e depois `docker compose up`.
8. Acesse o endereço `http://localhost:8080` para acessar o ThingsBoard. (Pode levar alguns instantes para ficar disponível)
