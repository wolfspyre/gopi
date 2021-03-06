#!/bin/bash
##############################################################
# Build Darwin (MacOS) Flavours
##############################################################

CURRENT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO=`which go`
LDFLAGS="-w -s"
TAGS="darwin"
cd "${CURRENT_PATH}/.."

##############################################################
# Sanity checks

if [ ! -d ${CURRENT_PATH} ] ; then
  echo "Not found: ${CURRENT_PATH}" >&2
  exit -1
fi
if [ "${GO}" == "" ] || [ ! -x ${GO} ] ; then
  echo "go not installed or executable" >&2
  exit -1
fi

##############################################################
# Generate

GENERATE=(
    rpc/rpc_server.go
)

for COMMAND in ${GENERATE[@]}; do
    echo "go generate cmd/${COMMAND}"
done

##############################################################
# Install

COMMANDS=(
    helloworld/helloworld.go
    timer/timer_tester.go
    rpc/rpc_server.go
)

for COMMAND in ${COMMANDS[@]}; do
    echo "go install cmd/${COMMAND}"
    go install -ldflags "${LDFLAGS}" -tags "${TAGS}" "cmd/${COMMAND}" || exit -1
done


