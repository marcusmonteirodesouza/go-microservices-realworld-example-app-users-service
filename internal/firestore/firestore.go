package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

func InitFirestore(projectId string) (*firestore.Client, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	return client, nil
}
