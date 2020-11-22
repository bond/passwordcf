#!/bin/sh

exec gcloud functions deploy passwordcf --region europe-west1 --entry-point GeneratePassword --runtime go113 --trigger-http --allow-unauthenticated --memory 128