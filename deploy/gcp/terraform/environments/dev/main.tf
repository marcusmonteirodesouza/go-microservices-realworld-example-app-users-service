provider "google" {
  project = local.project_id
  region  = local.region
}

resource "google_service_account" "github_deployer" {
  account_id = "users-service-github-deployer"
}

resource "google_project_iam_member" "github_deployer_workload_identity_pool_admin" {
  project = local.project_id
  role    = "roles/iam.workloadIdentityPoolAdmin"
  member  = "serviceAccount:${google_service_account.github_deployer.email}"
}

resource "google_project_iam_member" "github_deployer_service_account_admin" {
  project = local.project_id
  role    = "roles/iam.serviceAccountAdmin"
  member  = "serviceAccount:${google_service_account.github_deployer.email}"
}

module "github_oidc" {
  source      = "terraform-google-modules/github-actions-runners/google//modules/gh-oidc"
  project_id  = local.project_id
  pool_id     = "users-service-github-pool"
  provider_id = "users-service-github-provider"
  sa_mapping = {
    "users-service-github-deployer" = {
      sa_name   = google_service_account.github_deployer.id
      attribute = "*"
    }
  }

  depends_on = [
    google_project_iam_member.github_deployer_workload_identity_pool_admin,
    google_project_iam_member.github_deployer_service_account_admin
  ]
}
