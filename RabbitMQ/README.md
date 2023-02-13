# RabbitMQ

Simple RabbitMQ setup. Eventually we will replace the consumer with our transaction service and the sender with our 
frontend and our log file parser.

## Usage

Run Docker Compose: `docker-compose up --build`

Send a message via the sender: `http://localhost:3000/send?msg=testing`

View the message production and consumption via RabbitMQ UI using guest as the password and username: `http://localhost:15672`
