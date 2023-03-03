output "audit_dataset_id" {
  description = "The ID of the dataset "
  value       = google_bigquery_dataset.audit.dataset_id
}