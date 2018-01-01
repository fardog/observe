package observe

import (
	"context"
	"errors"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
)

type BigQueryOptions struct{}

func NewBigQuery(projectID, tableName string, opts *BigQueryOptions) (*BigQuery, error) {
	if opts == nil {
		opts = &BigQueryOptions{}
	}

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	tn := strings.Split(tableName, ".")
	if len(tn) != 2 {
		return nil, errors.New("invalid table name")
	}
	ds, tb := tn[0], tn[1]

	dataset := client.Dataset(ds)
	table := dataset.Table(tb)

	return &BigQuery{
		options: opts,
		client:  client,
		table:   table,
	}, nil
}

type BigQuery struct {
	options *BigQueryOptions
	client  *bigquery.Client
	table   *bigquery.Table
}

func (b *BigQuery) Store(o *Observation) error {
	u := b.table.Uploader()
	err := u.Put(context.Background(), []*value{valueFromObservation(o)})

	return err
}

type header struct {
	Key   string
	Value []string
}

func valueFromObservation(o *Observation) *value {
	var h []header

	for k, v := range o.Header {
		h = append(h, header{Key: k, Value: v})
	}

	return &value{
		URL:        o.URL,
		RemoteAddr: o.RemoteAddr,
		Observed:   o.Observed,
		Header:     h,
	}
}

type value struct {
	URL        string
	RemoteAddr string
	Observed   time.Time
	Header     []header
}

func (v value) Save() (map[string]bigquery.Value, string, error) {

	return map[string]bigquery.Value{
		"URL":        v.URL,
		"RemoteAddr": v.RemoteAddr,
		"Observed":   v.Observed.Format("2006-01-02T15:04:05-07:00"),
		"Header":     v.Header,
	}, "", nil
}
