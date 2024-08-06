FROM alpine:latest

RUN mkdir /app

COPY accountApp /app

CMD [ "/app/accountApp"]