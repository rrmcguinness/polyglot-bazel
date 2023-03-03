#Deploy audit shop dataset
resource "google_bigquery_dataset" "audit" {
  dataset_id                  = var.audit
  friendly_name               = "audit"
  description                 = "Audit Dataset"
  location                    = var.audit_DS_location #check the location
  #default_table_expiration_ms = 3600000
}