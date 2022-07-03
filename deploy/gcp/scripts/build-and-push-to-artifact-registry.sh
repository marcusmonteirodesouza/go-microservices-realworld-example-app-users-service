#!/bin/bash

PROJECT=$1
REGION=$2
IMAGE=$3

DOCKERFILE_PATH='../../..'

pushd "$DOCKERFILE_PATH" || exit 1
gcloud builds submit --project "$PROJECT" --region "$REGION" --tag "$IMAGE"
popd || exit
