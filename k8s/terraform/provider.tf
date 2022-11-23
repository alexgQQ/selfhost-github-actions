terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.38.0"
    }
  }

  required_version = ">= 0.14"
}

provider "google" {
  project = var.gke_project
  region  = var.gke_region
}
