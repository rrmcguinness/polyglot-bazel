provider "google" {
  region      = var.region
  zone        = var.zone
  project     = var.project_id
}

# modules
module "dataset" {
  source = "./dataset"
  event_DS_location = var.event_DS_location
  event = var.event

}

module "table" {
  source = "./table"
  event = module.dataset.event_dataset_id
  tbl_event = var.tbl_event
}