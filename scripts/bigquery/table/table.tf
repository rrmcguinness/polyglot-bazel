# notebook table
resource "google_bigquery_table" "tbl_audit" {
  dataset_id = var.audit
  table_id = var.tbl_audit

  time_partitioning {
    type = "DAY"
  }
  schema = <<EOF
[
  {
    "name": "id",
    "type": "STRING",
    "mode": "REQUIRED",
    "description": "System generated identifier"
  },
  {
    "name": "created",
    "type": "TIMESTAMP",
    "mode": "REQUIRED",
    "description": "The time the action occured"
  },
  {
    "name": "action",
    "type": "STRING",
    "mode": "REQUIRED",
    "description": "The action that was executed."
  },
  {
    "name": "context",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "Additional context values"
  },
  {
    "name": "principal",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "Additional context values"
  },
  {
    "name": "context_variables",
    "type": "RECORD",
    "mode": "REPEATED",
    "description": "Additional context variables",
    "fields": [
      {
        "name": "k",
        "type": "STRING",
        "mode": "REQUIRED",
        "description": "Key"
      },
      {
        "name": "v",
        "type": "STRING",
        "mode": "NULLABLE",
        "description": "Value"
      }
    ]
  }
]
EOF
}