FROM ubuntu:latest

WORKDIR /app

COPY  todo_final .
COPY  web ./web

RUN chmod +x finalTodoApp


ARG TODO_PORT=7540
ENV TODO_PORT=${TODO_PORT}

CMD [ "./todo_final" ]