FROM raysouz/aws_lambda_provided.al2023:go1.24

WORKDIR /app

# ENV HTTP_PROXY=
# ENV HTTPS_PROXY=
# ENV NO_PROXY=localhost,127.0.0.1
# ENV AWS_CA_BUNDLE=

# ENV GOINSECURE=
# ENV GOPROXY=
# ENV GONOPROXY=127.0.0.1,localhost,kubernetes.docker.internal
# ENV GOPRIVATE=

# Define as variáveis de ambiente GOROOT e PATH
ENV GOROOT=/usr/local/go
ENV PATH="$GOROOT/bin:${PATH}"

# Define a variável de ambiente GO111MODULE para habilitar módulos Go (opcional, mas recomendado)
ENV GO111MODULE=on

# Copia e compila o projeto para gerar o bootstrap executável
COPY app/. .
RUN go build -o bin/bootstrap cmd/lambda/main.go

# Define o ponto de entrada para o seu handler Go (ajuste conforme necessário)
ENTRYPOINT [ "/app/bin/bootstrap" ]