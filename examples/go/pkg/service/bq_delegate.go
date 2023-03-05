package service

import (
	"cloud.google.com/go/bigquery"
	"context"
	"examples/go/pb"
	"examples/go/pkg/model"
	"flag"
	"fmt"
	"google.golang.org/api/iterator"
)

var (
	gcpProjectId   = flag.String("gcp-project-id", "polyglot-bazel", "The ID of your project.")
	eventDatasetId = flag.String("event-ds-id", "ds_events", "The ID of the dataset to write to.")
	eventTableId   = flag.String("event-table-id", "tbl_event", "The ID of the table to write to.")
)

func getClient(ctx context.Context) (*bigquery.Client, error) {
	client, err := bigquery.NewClient(ctx, *gcpProjectId)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %w", err)
	}
	return client, nil
}

func insertEventRecord(record *pb.EventRecord) error {

	// Create a version of the event record that can be persisted.
	// here the persistable event maps the struct property names to the BigQuery (BQ) table structure.
	persistableEvent := model.NewServiceEventRecord(record)

	ctx := context.Background()
	client, err := getClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	inserter := client.Dataset(*eventDatasetId).Table(*eventTableId).Inserter()

	if err := inserter.Put(ctx, persistableEvent); err != nil {
		return err
	}
	return nil
}

const fmtQueryRange = "SELECT * FROM `%s.%s` WHERE TIMESTAMP_TRUNC(observed, DAY) between TIMESTAMP(DATE(\"%s 00:00:00+00\", \"UTC\")) AND TIMESTAMP(DATE(\"%s 00:00:00\", \"UTC\"))"
const fmtDate = "2006-01-02"

func findByDateRange(request *pb.DateRangeRequest) ([]*pb.EventRecord, error) {
	out := make([]*pb.EventRecord, 0)
	if request != nil && request.End != nil &&
		request.Start != nil && request.End.AsTime().After(request.Start.AsTime()) {
		endDate := request.End.AsTime().Format(fmtDate)
		startDate := request.Start.AsTime().Format(fmtDate)
		query := fmt.Sprintf(fmtQueryRange, *eventDatasetId, *eventTableId, startDate, endDate)

		ctx := context.Background()
		client, err := getClient(ctx)
		if err != nil {
			return out, err
		}
		defer client.Close()

		q := client.Query(query)
		it, err := q.Read(ctx)
		if err != nil {
			return out, err
		}

		for {
			var rpw model.ServiceEventRecord
			err := it.Next(&rpw)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return out, err
			}
			out = append(out, model.EventRecordFromServiceEventRecord(&rpw))
		}

	}

	return out, nil
}

func verify(request *pb.EventResponse) bool {
	return false
}
