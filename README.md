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

1. Install the [gcloud CLI](https://cloud.google.com/sdk/docs/install).
1. Install [terraform](https://www.terraform.io).
1. `cd` into your [environment's directory](./deploy/gcp/terraform/environments/dev). Update the `backend.tf` to use the Storage Bucket created during the project's bootstraping as your [backend](https://www.terraform.io/language/settings/backends/gcs). Update the `locals.tf` with your project id and other values as you seem fit.
1. Run `terraform init` and `terraform apply`. It will setup the [Workload Identity Provider](https://github.com/terraform-google-modules/terraform-google-github-actions-runners/tree/v3.0.0/modules/gh-oidc) to be used in the CI/CD pipeline.
1. The [`deploy` workflow](./.github/workflows/deploy.yaml) will build and deploy the Docker container image to GCP's [Artifact Registry](https://cloud.google.com/artifact-registry). Create the secrets it needs according to the [example](https://github.com/terraform-google-modules/terraform-google-github-actions-runners/tree/v3.0.0/modules/gh-oidc#github-workflow).
