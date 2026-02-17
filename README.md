# backoffice

## Backoffice é um consumidor de filas rabbitMQ.

### Como executar:

1. Tenha o docker instalado;
2. Tenha o Golang instalado (De preferência 1.25);
3. Crie um arquivo .env com as mesmas envs do env.example;
4. Execute make up para subir as imagens do rabbitmq e redis;
5. Execute make backoffice para construir e executar a imagem da aplicação;
6. O envio de mensagens pode ser realizado tanto com rabbitmq ui ou via script na pasta test;

#### Teste via [RabbitMQ UI](http://localhost:15672/):
    1. Via rabbitmq ui basta logar com usuário: user e senha: password;
    2. Clicar em queues and streams, selecionar a fila task.in;
    3. Rolar ate publish message, inserir a mensagem que desejada e enviar;

#### Via script:
    1. Definir a mensagem que deseja enviar em integration_test.json;
    2. Executar make integration;


### Decisões de implementação.

* A arquitetura foi baseada no [Go scaffold](https://github.com/go-scaffold/go-scaffold)
* Foi utilizado abstração da implementação do RabbitMQ para evitar retrabalho caso o estilo de fila/stream mude para kafka ou sqs.
* Foi implementado tanto o consumer da fila task.in conforme requisito, quanto da task.out.
* No consumo de task.out foi adicionado a persistência da mensagem enviado via task.in.
* Além das fila de task.in e task.out, foi criada uma fila de dlq, no qual recebe toda mensagem recutada no processamento de task.in e task.out.

### Ferramentas

1. Golang 1.25;
2. Redis;
3. RabbitMQ;
4. Gracefull shutdown;
5. Docker;
6. Jaeger

#### Para acessar o jaeger e pesquisar os traces basta acessar [jaeger local](http://localhost:16686/search)