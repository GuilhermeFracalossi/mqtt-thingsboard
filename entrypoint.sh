#!/bin/bash

# Variáveis de ambiente
THINGSBOARD_HOST=${THINGSBOARD_HOST:-"mytb"}
THINGSBOARD_PORT=${THINGSBOARD_PORT:-1883}

echo "Aguardando o ThingsBoard iniciar em $THINGSBOARD_HOST:$THINGSBOARD_PORT..."

# Loop até que a porta esteja disponível
while ! nc -z $THINGSBOARD_HOST $THINGSBOARD_PORT; do
    sleep 2
done

echo "ThingsBoard está pronto. Iniciando o sensor..."

# Executa o comando original
exec "$@"
