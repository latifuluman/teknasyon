#!/bin/sh
# Ensure the correct permissions on the .erlang.cookie file
chown rabbitmq:rabbitmq /var/lib/rabbitmq/.erlang.cookie
chmod 400 /var/lib/rabbitmq/.erlang.cookie

# Start RabbitMQ
exec docker-entrypoint.sh rabbitmq-server
