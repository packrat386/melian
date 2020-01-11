package melian

import (
	"context"
	"net/http"
)

type extraKeyType struct{}

var extraKey = extraKeyType{}

func AddExtra(ctx context.Context, key string, value interface{}) context.Context {
	extra, ok := getExtra(ctx)
	if !ok || extra == nil {
		extra = map[string]interface{}{}
	}

	extra[key] = value
	return context.WithValue(ctx, extraKey, extra)
}

func getExtra(ctx context.Context) (map[string]interface{}, bool) {
	extra, ok := ctx.Value(extraKey).(map[string]interface{})
	return extra, ok
}

type tagKeyType struct{}

var tagKey = tagKeyType{}

func AddTag(ctx context.Context, key string, value string) context.Context {
	tags, ok := getTags(ctx)
	if !ok || tags == nil {
		tags = map[string]string{}
	}

	tags[key] = value
	return context.WithValue(ctx, tagKey, tags)
}

func getTags(ctx context.Context) (map[string]string, bool) {
	tags, ok := ctx.Value(tagKey).(map[string]string)
	return tags, ok
}

type requestKeyType struct{}

var requestKey = requestKeyType{}

func AddRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, fromHTTPRequest(req))
}

func getRequest(ctx context.Context) (request, bool) {
	r, ok := ctx.Value(requestKey).(request)
	return r, ok
}

type transactionKeyType struct{}

var transactionKey = transactionKeyType{}

func AddTransaction(ctx context.Context, transaction string) context.Context {
	return context.WithValue(ctx, transactionKey, transaction)
}

func getTransaction(ctx context.Context) (string, bool) {
	transaction, ok := ctx.Value(transactionKey).(string)
	return transaction, ok
}
