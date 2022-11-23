data "google_app_engine_default_service_account" "default" {
  project = var.gke_project
}

resource "google_container_cluster" "actions_cluster" {
  name     = "primary-cluster"
  location = var.gke_zone

  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "primary_node" {
  name       = "primary-node"
  location   = var.gke_zone
  cluster    = google_container_cluster.actions_cluster.name
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "e2-small"
    disk_size_gb = 20
    disk_type    = "pd-standard"
    image_type   = "COS_CONTAINERD"

    # The GAE app deploy uses its default service account
    # lets use that as the default sa for our cluster so the runners have the same perms
    service_account = data.google_app_engine_default_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/service.management",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}
