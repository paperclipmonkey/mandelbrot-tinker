// Create a Cloud Run service
resource "google_cloud_run_v2_service" "mandelbrot" {
  name     = "mandelbrot"
  location = "europe-west2" // London
  deletion_protection = true
  ingress = "INGRESS_TRAFFIC_ALL"
  project = "mandelbrot-tinker"

  template {
    containers {
      image = "59vkckvlkjdfglkjdfv/mandelbrot-tinker:${var.image_tag}"
      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
      startup_probe {
        initial_delay_seconds = 0
        timeout_seconds = 1
        period_seconds = 3
        failure_threshold = 1
        tcp_socket {
          port = 8080
        }
      }
      liveness_probe {
        http_get {
          path = "/livez"
        }
      }
    }
  }
}

// Allow all users to invoke the service
resource "google_cloud_run_service_iam_binding" "mandelbrot" {
  project = "mandelbrot-tinker"
  location = google_cloud_run_v2_service.mandelbrot.location
  service  = google_cloud_run_v2_service.mandelbrot.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}