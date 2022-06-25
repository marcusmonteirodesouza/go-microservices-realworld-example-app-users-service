locals {
  project_id = "marcus-go-ms-rw-example-1"
  region = "us-central1"
  image = "${google_artifact_registry_repository.users_service.location}-docker.pkg.dev/${local.project_id}/${google_artifact_registry_repository.users_service.repository_id}/users-service"
}
