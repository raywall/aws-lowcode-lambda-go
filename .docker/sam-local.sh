#!/bin/bash

# Define a porta que o SAM Local vai usar
SAM_PORT=3000

echo "Iniciando o SAM Local API na porta ${SAM_PORT}..."
echo "Usando o template: $SAM_TEMPLATE"

# Navega até o diretório do seu projeto mapeado (opcional, dependendo de onde o template está)
cd /local/workspace

# Executa o comando do SAM Local para iniciar a API, especificando o template
sam local start-api --host 0.0.0.0 --port ${SAM_PORT} --template "$SAM_TEMPLATE"

echo "SAM Local API iniciada. Aguardando requisições..."

# Mantenha o script rodando em foreground
tail -f /dev/null