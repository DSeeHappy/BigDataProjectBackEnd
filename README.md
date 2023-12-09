### Final Project: Big Data
### Juan Daniel Sanchez Chavez


# Overview
This project allows a user to create a job which includes a location.
The server would obtain the weather from OpenWeatherMap API and store it in the database.
You would then be able to update the job to schedule based on the 30-day forecast of that location.

The idea would be if you are planning on creating your work schedule for the next month,
you would be able to see the weather for the next 30 days and schedule your work accordingly.

If this project were to be extended further, the objective would be to track and be notified of weather changes.
Helping the user make sure they don't schedule work when there is a major weather event, useful for outdoor jobs.
Plans would be to include test message, email notification or device notification if extended further beyond project scope.

# Setup

Start the server with docker-compose
Initiates go server, postgres db, prometheus, grafana, rabbitmq
```bash
docker-compose up
```

# Setup DB
navigate to dbscripts and select public_schema.sql
run the script in the postgres db to create the tables

# Testing
Unit testing can be run without any docker containers running

Integration testing requires the docker containers to be running

To test the server, run the following command
```bash
go test ./...
```

# Health monitoring
Prometheus is used to monitor the health of the server
Grafana is used to visualize the metrics

The health monitoring tracks metrics for job and weather api calls

# Deployment
The server is deployed to AWS using ECS
The server is deployed using a docker image
The docker image is built using the Dockerfile

JetBrains Space Packages is used to build the docker image and push it to JetBrains Space Container Registry
JetBrains Space Container Registry is also used to deploy the docker image to AWS ECS
