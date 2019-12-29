# go-gcp

# Pre-Requistes

```{bash}
sudo snap install google-cloud-sdk --classic
gcloud config set project <MyProject>
gcloud auth login
```

## Project Setup

This was a small project to play with some face recognition using GCP functions

- Go to https://console.cloud.google.com/home/dashboard?project=MyProject
- Go to Resource->Storage to see the buckets.
- Create a bucket to watch for incoming webcam images, examing `incoming-images`
- Create another bucket called `outgoing-images`

## Build

```{bash}
go build && go test
```

## Deploy

```{bash}
./deploy.sh
```