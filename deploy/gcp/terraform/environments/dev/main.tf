provider "google" {
  project = local.project_id
  region  = local.region
}

provider "google-beta" {
  project = local.project_id
  region  = local.region
}

resource "google_project_service" "artifactregistry" {
  service            = "artifactregistry.googleapis.com"
  disable_on_destroy = false
}

resource "google_artifact_registry_repository" "users_service" {
  provider = google-beta

  location      = local.region
  repository_id = "users-service-repository"
  description   = "Users Service Repository"
  format        = "DOCKER"

  depends_on = [
    google_project_service.artifactregistry
  ]
}

resource "google_project_service" "cloudbuild" {
  service            = "cloudbuild.googleapis.com"
  disable_on_destroy = false
}

resource "null_resource" "build_and_push_container_image" {
  provisioner "local-exec" {
    command     = "./build-and-push-container-image.sh '${local.project_id}' '${local.region}' '${local.image}'"
    working_dir = "${path.module}/../../../scripts"
  }

  depends_on = [
    google_project_service.cloudbuild
  ]
}
