# ![RealWorld Example App](logo.png) - Users Service

> ### [Go](https://go.dev/) codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://demo.realworld.io/)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully fledged fullstack application built with **[Go](https://go.dev/)** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **[Go](https://go.dev/)** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# How it works

> Describe the general architecture of your app here

# Getting started

## Running locally

1. Install [Docker](https://docs.docker.com/get-docker/).
1. Run `docker compose up`.

## Testing

1. Run `./test.sh`.

# Deploy

## Deploy to [Google Cloud Platform](https://cloud.google.com/) (GCP)

1. The [`deploy-dev` workflow](./.github/workflows/deploy-dev.yaml) will build and deploy the Docker container image to GCP's [Artifact Registry](https://cloud.google.com/artifact-registry). Add the secrets it needs to run, which were created during the project's bootstrap.
