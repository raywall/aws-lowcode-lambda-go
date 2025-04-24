# Defina as variáveis importantes
IMAGE_NAME = raysouz/aws_lambda_provided.al2023
POSTMAN_NAME = raysouz/aws_lambda_newman
TAG = go1.24

# Comando para construir a imagem Docker
build-image:
	@echo "Construindo a imagem Docker..."
	docker build -f .docker/Repository -t $(IMAGE_NAME):$(TAG) .
	@echo "Imagem Docker construída com sucesso: $(IMAGE_NAME):$(TAG)"

build-postman:
	@echo "Construindo a imagem Docker..."
	docker build -f .docker/Postman -t $(POSTMAN_NAME):latest .
	@echo "Imagem Docker construída com sucesso: $(POSTMAN_NAME):latest"

# Comando para enviar a imagem para o Docker Hub
push:
	@echo "Enviando a imagem Docker para o Docker Hub..."
	docker push $(IMAGE_NAME):$(TAG)
	@echo "Imagem Docker enviada com sucesso para o Docker Hub: $(IMAGE_NAME):$(TAG)"

# Comando para enviar a imagem para o Docker Hub
push-postman:
	@echo "Enviando a imagem Docker para o Docker Hub..."
	docker push $(POSTMAN_NAME):latest
	@echo "Imagem Docker enviada com sucesso para o Docker Hub: $(POSTMAN_NAME):latest"

# Comando para construir e enviar a imagem em uma única etapa
build-and-push: build-image push
	@echo "Processo de build e push concluído."

compose-up:
	@echo "Construindo o Docker compose..."
	docker compose -f .docker/docker-compose.yaml up -d

compose-down:
	@echo "Destruindo o Docker compose..."
	docker compose -f .docker/docker-compose.yaml down

# Limpar quaisquer imagens locais (opcional)
clean:
	@echo "Removendo imagens Docker locais (opcional)..."
	docker rmi -f $(IMAGE_NAME):$(TAG)
	docker rmi -f $(POSTMAN_NAME):latest
	@echo "Imagens removidas."

.PHONY: build-image build-postman push push-postman build-and-push compose-up compose-down clean