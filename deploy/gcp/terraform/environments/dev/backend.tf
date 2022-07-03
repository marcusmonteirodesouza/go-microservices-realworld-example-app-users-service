terraform {
  backend "gcs" {
    bucket = "tfstate-marcus-go-ms-rw-example-1-dev"
    prefix = "users-service"
  }
}
