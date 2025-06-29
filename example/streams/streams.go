package main

import (
	"context"
	"encoding/base64"
	"ksql/config"
	"ksql/kinds"
	"ksql/ksql"
	"ksql/shared"
	"ksql/streams"
	"log/slog"
)

const (
	ksqlURL = "http://localhost:8088"
)

func main() {
	//ctx := context.Background()
	//List(ctx)
	//Create(ctx)
	//Describe(ctx)
	//Drop(ctx)
	//t := time.NewTicker(5 * time.Second)
	//go func() {
	//	for {
	//		select {
	//		case <-t.C:
	//			Insert(ctx)
	//		}
	//	}
	//}()
	//Insert(ctx)
	//Select(ctx)
	//SelectWithEmit(ctx)
	//CreateAsSelect(ctx)
	//InsertAsSelect(ctx)
}

func init() {
	cfg := config.New(ksqlURL, 15, false)
	if err := cfg.Configure(context.Background()); err != nil {
		slog.Error("cannot configure ksql", "error", err.Error())
	}
}

func List(ctx context.Context) {
	streamsList, err := streams.ListStreams(ctx)
	if err != nil {
		slog.Error("cannot list streams", "error", err.Error())
		return
	}

	slog.Info("successfully executed!", "streams", streamsList)
}

type ExampleStream struct {
	ID    int    `ksql:"ID"`
	Token []byte `ksql:"TOKEN"`
}

const (
	streamName = "MY_STREAM"
)

func Create(ctx context.Context) {
	sourceTopic := "examples-topics"
	//partitions := 1 // if topic doesnt exists, partitions are required

	exampleTable, err := streams.CreateStream[ExampleStream](
		ctx, streamName, shared.StreamSettings{
			SourceTopic: &sourceTopic,
			Format:      kinds.JSON,
		})

	if err != nil {
		slog.Error("cannot create stream", "error", err.Error())
		return
	}

	slog.Info("stream created!", "name", exampleTable.Name)
}

func Describe(ctx context.Context) {
	description, err := streams.Describe(ctx, streamName)
	if err != nil {
		slog.Error("cannot describe stream", "error", err.Error())
		return
	}

	slog.Info("successfully executed", "description", description)
}

func Drop(ctx context.Context) {
	if err := streams.Drop(ctx, streamName); err != nil {
		slog.Error("cannot drop stream", "error", err.Error())
		return
	}

	slog.Info("stream dropped!", "name", streamName)
}

func Insert(ctx context.Context) {
	exampleStream, err := streams.GetStream[ExampleStream](ctx, streamName)
	if err != nil {
		slog.Error("cannot get stream", "error", err.Error())
		return
	}

	data := []byte("SECRET_BASE64_DATA")
	token := []byte(base64.StdEncoding.EncodeToString(data))

	if err = exampleStream.Insert(ctx, ksql.Row{
		"ID":    1,
		"TOKEN": token,
	}); err != nil {
		slog.Error("cannot insert data to stream", "error", err.Error())
		return
	}

	slog.Info("successfully inserted")
}

func Select(ctx context.Context) {
	exampleStream, err := streams.GetStream[ExampleStream](ctx, streamName)
	if err != nil {
		slog.Error("cannot get stream", "error", err.Error())
		return
	}

	rows, err := exampleStream.SelectOnce(ctx)
	if err != nil {
		slog.Error("cannot select from stream", "error", err.Error())
		return
	}

	slog.Info("successfully selected rows", "rows", rows)
}

func SelectWithEmit(ctx context.Context) {
	// Fix tommorow: invalid select builder: GROUP BY requires WINDOW clause on streams"
	exampleStream, err := streams.GetStream[ExampleStream](ctx, streamName)
	if err != nil {
		slog.Error("cannot get stream", "error", err.Error())
		return
	}

	notesStream, err := exampleStream.SelectWithEmit(ctx)
	if err != nil {
		slog.Error("error during emit", "error", err.Error())
		return
	}

	for note := range notesStream {
		slog.Info("received note", "note", note)
	}
}

func CreateAsSelect(ctx context.Context) {
	sql := ksql.Select(ksql.F("ID"), ksql.F("TOKEN")).From("EXAMPLESTREAM", ksql.STREAM)
	sourceTopic := "examples-topics"
	dublicateStream, err := streams.CreateStreamAsSelect[ExampleStream](ctx, "dublicateStream", shared.StreamSettings{
		SourceTopic: &sourceTopic,
		Format:      kinds.JSON,
	}, sql)
	if err != nil {
		slog.Error("cannot create stream as select", "error", err.Error())
		return
	}

	slog.Info("stream created!", dublicateStream.Name)
}

func InsertAsSelect(ctx context.Context) {
	sql := ksql.Select(ksql.F("ID"), ksql.F("TOKEN")).From("DUBLICATESTREAM", ksql.STREAM)

	stream, err := streams.GetStream[ExampleStream](ctx, "EXAMPLESTREAM")
	if err != nil {
		slog.Error("cannot get stream", "error", err.Error())
		return
	}

	err = stream.InsertAsSelect(ctx, sql)
	if err != nil {
		slog.Error("cannot insert as select to stream", "error", err.Error())
		return
	}

	slog.Info("inserted as select")
}
