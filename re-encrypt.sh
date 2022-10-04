#!/bin/bash

# This is a little helper script, which enables you to re-encrypt Sealedsecrets after the 30 day expiration of the
# certificate and key. Make sure you use the file used for the creation of the initial inventory and are connected to the cluster

usage(){
  echo """Usage: ./re-encrypt.sh cluster
  Make sure you use a current configuration for the cluster named 'config.yaml' and have it on the same level as this script.
  Don't forget to commit it to the private inventory repository on github if there have been any changes."""
}

main() {
set -e

CLUSTER="${1}"
if test -z "${1}"; then
  usage
  exit 1
fi

mkdir -p src/generated

CURRSECRET=$(kubectl -n gp-infrastructure get secret --sort-by metadata.creationTimestamp | grep sealed-secrets-key | head -n 1 | cut -d " " -f 1)
kubectl -n gp-infrastructure get secret ${CURRSECRET} -o jsonpath='{.data.tls\.key}' | base64 -d > src/generated/"${CLUSTER}".key
kubectl -n gp-infrastructure get secret ${CURRSECRET} -o jsonpath='{.data.tls\.crt}' | base64 -d > src/generated"/${CLUSTER}".crt

cd src || exit 1
go build -o day-x-generator gepaplexx/day-x-generator
./day-x-generator ../config.yaml

rm day-x-generator && cd ..

}

main "${1}"

