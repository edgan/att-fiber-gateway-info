#!/bin/bash

DEBUG="$1"

BUILDS_DIR='bin'
BUILDS_FILE='.go_builds'
DELIMITER='_'
NAME=`grep module go.mod | cut -f2 -d/`

readarray -t BUILDS < ${BUILDS_FILE}

for BUILD in "${BUILDS[@]}"; do
  BUILD_CMD='go build'
  BUILD_FLAGS_STRIP='-ldflags="-w -s"'
  SPLIT=(${BUILD//${DELIMITER}/ })
  GOOS=${SPLIT[0]} # Opearting system
  GOARCH=${SPLIT[1]} # Architecture

  if [ "${DEBUG}" != "debug" ]; then
    if [ "${GOOS}" == "linux" ]; then
      BUILD_CMD="${BUILD_CMD} ${BUILD_FLAGS_STRIP}"
    fi
  fi

  if [ "${GOOS}" == "windows" ]; then
    EXTENSION='.exe'
  fi

  BUILD_FILENAME="${NAME}${DELIMITER}${BUILD}${EXTENSION}"
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} ${BUILD_CMD} -o ${BUILDS_DIR}/${BUILD_FILENAME}"
  eval ${CMD}

done
