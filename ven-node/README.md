# VEC-A3C

This repository contains the source code for the **VEC-A3C** project.

## Description

The project utilizes Docker containers for development and production environments. There are two Dockerfiles included in this repository:

1. `Dockerfile`: Used for building the "production" level image. After making changes to the image, ensure that the `/cue/ven_template.cue` file is updated with the link to the latest image.

2. `Dockerfile.dev`: Used for development purposes and should be referenced in the `docker-compose.yaml` file.

## Getting Started

To run the project on your local machine, follow these steps:

1. Clone this repository to your local machine.

   ```bash
   git clone https://github.com/Durbek-Gafur/vec-a3c.git

2. Change to the project directory.

    ```bash
    cd repo

3. Start the project using Docker Compose.

    ```bash
    docker-compose up

This command will build and start the necessary containers based on the configuration specified in the docker-compose.yaml file.

## Updating the Database
If you need to update the database, you can do so by writing alter table SQL commands in the /migrations folder. After making the necessary changes, follow these additional steps:

Build and push a new image of the database to Dockerhub.

Update the /cue/ven_template.cue file with the link to the latest database image.


## Building images
docker build -t 39dj29dl2d9l2/vec-ven:13 .
docker push 39dj29dl2d9l2/vec-ven:13

## Post workflow

```
curl -X POST http://localhost:8080/workflow -H "Content-Type: application/jso
n" -d '{
            "name": "workflow1",
            "type": "demo_50per.fastq",
            "received_at": "2023-07-23T12:34:56Z"
        }'s




curl -X POST --http1.1 https://dgvkh-ven7.nrp-nautilus.io/workflow -H "Content-Type: application/json" -d '{
    "name": "workflow15",
    "type": "demo_50per.fastq",
}'

curl -X POST https://dgvkh-scheduler.nrp-nautilus.io/workflow \
                                                       -H "Content-Type: application/json" \
                                                       -d '{
                                                       "Name": "workflow11",
                                                       "Duration": "5m30s",
                                                       "StartedExecutionAt": "2023-07-23T12:34:56Z",
                                                       "CompletedAt": "2023-07-23T12:40:26Z"
                                                   }'
```