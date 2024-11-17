#!/bin/bash

TESTS="$1"

if [ "${TESTS}" == "" ]; then
  TESTS=all
fi

COMMAND='./att-fiber-gateway-info -action'
SEPARATOR='##################################################'

if [[ ${TESTS} == "nologin" || ${TESTS} == "all" ]]; then
  NO_LOGIN_ACTIONS=(broadband-status device-list fiber-status home-network-status system-information)

  for ACTION in ${NO_LOGIN_ACTIONS[@]}; do
    echo
    echo "${SEPARATOR} ${ACTION} ${SEPARATOR}"
    ${COMMAND} ${ACTION}
  done
fi

if [[ ${TESTS} == "login" || ${TESTS} == "all" ]]; then
  LOGIN_ACTIONS=(ip-allocation nat-check nat-connections nat-destinations nat-sources nat-totals)

  for ACTION in ${LOGIN_ACTIONS[@]}; do
    echo
    echo "${SEPARATOR} ${ACTION} ${SEPARATOR}"
    ${COMMAND} ${ACTION}
  done
fi

if [[ ${TESTS} == "reset" || ${TESTS} == "all" ]]; then
  RESET_ACTIONS=(reset-connection reset-device reset-firewall reset-ip reset-wifi restart-gateway)

  for ACTION in ${RESET_ACTIONS[@]}; do
    echo
    echo "${SEPARATOR} ${ACTION} ${SEPARATOR}"
    ${COMMAND} ${ACTION} -no
  done
fi
