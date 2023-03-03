provider "google" {
  region      = var.region
  zone        = var.zone
  project     = var.project_id
}

# modules
module "dataset" {
  source = "./dataset"
  audit_DS_location = var.audit_DS_location
  audit = var.audit

}

module "table" {
  source = "./table"
  audit = module.dataset.audit_dataset_id
  tbl_audit = var.tbl_audit
}