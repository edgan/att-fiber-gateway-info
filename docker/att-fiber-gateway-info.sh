#!/bin/bash

ARGS='-allmetrics -continuous -datadog'

if [ ! -z "${PASSWORD}" ]; then
  ARGS="${ARGS} -password ${PASSWORD}"
fi

if [ ! -z "${URL}" ]; then
  ARGS="${ARGS} -url ${URL} "
fi

if [ ! -z "${PASSWORD}" ]; then
  ARGS="${ARGS} -password ${PASSWORD}"
fi

if [ ! -z "${STATSDIPPORT}" ]; then
 ARGS="${ARGS} -statsdipport ${STATSDIPPORT}"
fi

att-fiber-gateway-info ${ARGS}
