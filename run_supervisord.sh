#!/bin/bash

function run_offline_mode() {
  export HEDERA_MIRROR_ROSETTA_NETWORK=$NETWORK
  supervisord --configuration supervisord-offline.conf
}

function run_online_mode() {
  if [[ -n "$HEDERA_MIRROR_IMPORTER_DOWNLOADER_ACCESS_KEY" ]]; then
    yq write -i application.yml hedera.mirror.importer.downloader.accessKey $HEDERA_MIRROR_IMPORTER_DOWNLOADER_ACCESS_KEY
  fi

  if [[ -n "$HEDERA_MIRROR_IMPORTER_DOWNLOADER_BUCKET_NAME" ]]; then
    yq write -i application.yml hedera.mirror.importer.downloader.bucketName $HEDERA_MIRROR_IMPORTER_DOWNLOADER_BUCKET_NAME
  fi

  if [[ -n "$HEDERA_MIRROR_IMPORTER_DOWNLOADER_CLOUD_PROVIDER" ]]; then
    yq write -i application.yml hedera.mirror.importer.downloader.cloudProvider $HEDERA_MIRROR_IMPORTER_DOWNLOADER_CLOUD_PROVIDER
  fi

  if [[ -n "$HEDERA_MIRROR_IMPORTER_DOWNLOADER_GCP_PROJECT_ID" ]]; then
    yq write -i application.yml hedera.mirror.importer.downloader.gcpProjectId $HEDERA_MIRROR_IMPORTER_DOWNLOADER_GCP_PROJECT_ID
  fi

  if [[ -n "$HEDERA_MIRROR_IMPORTER_DOWNLOADER_SECRET_KEY" ]]; then
    yq write -i application.yml hedera.mirror.importer.downloader.secretKey $HEDERA_MIRROR_IMPORTER_DOWNLOADER_SECRET_KEY
  fi

  if [[ -n "$HEDERA_MIRROR_IMPORTER_NETWORK" ]]; then
    yq write -i application.yml hedera.mirror.importer.network $HEDERA_MIRROR_IMPORTER_NETWORK
  fi

  supervisord
}

case $MODE in
    "offline")
      run_offline_mode
    ;;
    *)
      run_online_mode
    ;;
esac