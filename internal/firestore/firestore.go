package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

func InitFirestore(ctx context.Context, projectId string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	return client, nil
}
