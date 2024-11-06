# EDSB
O EDSB (É DIFÍCIL SER BRASILEIRO) é uma plataforma/aplicação web que permite que brasileiros se juntem para discutir sobre suas existências em sua terra tupiniquin. A plataforma visa fomentar debates sobre cultura, política, e cotidiano, promovendo um espaço de apoio e troca de experiências entre os usuários.

DOMÍNIO: edificilserbrasileiro.com.br

# Pré-requisistos

- Certifique-se de que o [Docker](https://www.docker.com/get-started) e o [Docker Compose](https://docs.docker.com/compose/) estão instalados no servidor.

## Executando o Projeto

Para iniciar o projeto, execute o seguinte comando no terminal:

```bash
docker-compose up --build
```
Isso irá construir a imagem do Docker e iniciar todos os serviços definidos no docker-compose.yml.

## Acessando a Aplicação

Após iniciar o projeto, você pode acessar a aplicação em seu navegador através de `http://localhost:8000` (ou a porta especificada no seu `docker-compose.yml`).

# Estrutura EDSB
```
edsb
.
├── Dockerfile
├── README.md
├── api
│   ├── comment
│   │   └── comment.go
│   ├── like
│   │   └── like.go
│   ├── post
│   │   └── post.go
│   └── user
│       └── user.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── models
├── static
└── templates

9 directories, 10 files
```

## Diretórios e subdiretórios
- `api`: Este diretório contém a lógica de negócios e as interações com o banco de dados.

    - `comment`: Gerencia os comentários (como a tabela de comentários).
    - `like`: Lida com a lógica de likes (como a tabela de likes).
    - `post`: Trata a lógica dos posts (como a tabela de posts).
    - `user`: Contém a lógica relacionada aos usuários (como a tabela de usuários).

- `models`: Este diretório pode ser usado para definir as estruturas de dados (structs) que correspondem às suas tabelas do banco de dados. Isso ajuda a mapear os dados que você recebe e envia.

- `static`: Você pode colocar arquivos estáticos aqui, como CSS, JavaScript e imagens que sua aplicação web pode servir.

- `templates`: Coloque seus arquivos de template HTML aqui, caso você esteja usando um framework como o html/template do Go para renderizar páginas web.