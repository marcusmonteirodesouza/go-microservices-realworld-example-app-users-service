output "github_oidc_provider" {
  description = "Workload Identity Provider name."
  value = module.github_oidc.provider_name
}
