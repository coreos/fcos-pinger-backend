## Description

Backend for Fedora CoreOS Pinger (https://github.com/coreos/fedora-coreos-pinger).

Currently the backend server will only interpret the sent request body as a [`postData`](https://github.com/zonggen/fcos-pinger-backend/blob/67020ce36c8d1d2c38b0a23c1d2bba6d8189539f/main.go#L33) struct and store the `minimal` and `full` data to the corresponding `minimal` and `full` collections under `fcos_pinger` database.

## Pre-requisites:
 - `podman`: https://podman.io/getting-started/installation.html
 - `podman-compose`: https://github.com/containers/podman-compose

## Docker images used:
 - MongoDB official image from Docker Hub ([link_to_image](https://hub.docker.com/_/mongo))
 - Red Hat Universal Base Image 8 (UBI) ([link_to_image](https://access.redhat.com/containers/#/registry.access.redhat.com/ubi8/ubi))

## To run locally

```bash
$ podman-compose up
```

## To stop running

```bash
$ podman-compose down
```

## Networks
The default port for mongod and mongos instances is 27017, and the backend server is listening on port 5000. The networks between docker containers and host are shared, therefore to send POST request to server, all one needed was to send request to `127.0.0.1:5000`.

## Tools used during development
 - MongoDB Compass to visualize the DB ([link_to_official-page](https://docs.mongodb.com/compass/master/install/))
 - Postman for visualizing requests ([link_to_official_page](https://www.getpostman.com/))
