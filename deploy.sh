#!/bin/bash 

go build && \
  go vet && \
  go test && \
  gcloud functions deploy DetectAndCrop \
    --runtime go111 \
    --trigger-resource incoming-images \
    --trigger-event google.storage.object.finalize
