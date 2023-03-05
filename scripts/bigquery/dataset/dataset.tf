#Deploy audit shop dataset
resource "google_bigquery_dataset" "event" {
  dataset_id                  = var.event
  friendly_name               = "events"
  description                 = "Event Dataset"
  location                    = var.event_DS_location #check the location
  #default_table_expiration_ms = 3600000
}