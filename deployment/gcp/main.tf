terraform {
    backend "gcs" {
        bucket  = "mandelbrot-tinker-terraform"
        prefix  = "terraform/state"
    }
}