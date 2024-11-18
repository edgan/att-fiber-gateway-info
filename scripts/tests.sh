#!/bin/bash

if [[ "$TERM" == "xterm"* ]] || [[ "$TERM" == "screen"* ]] || [[ "$TERM" == "linux" ]] || [[ "$TERM" == "tmux" ]]; then
	SUPPORTS_COLOR=true
else
	SUPPORTS_COLOR=false
fi

# Text colors
BLUE='\e[34m'
GREEN='\e[32m'
RESET='\e[0m'

# Function to print colored text conditionally
print_color_conditional() {
  local color=$1
  local text=$2

  if [ "${SUPPORTS_COLOR}" == "true" ]; then
    # Directly print the color-corrected text using indirect expansion
    echo -e "${!color}${text}${RESET}"
  else
    echo "${text}"
  fi
}

print_separator() {
  local full_command="$1"
  local separator="##################################################"
  local text="${separator} ${full_command} ${separator}"
  print_color_conditional "BLUE" "${text}"
}

run_commands() {
  local flag="$1"
  shift
  local actions=("$@") # Store remaining arguments as actions

  if [ "${#actions[@]}" -eq 0 ]; then
    local full_command="${COMMAND} ${flag}"
    print_separator "${full_command}"
    ${full_command}
  else
    for action in "${actions[@]}"; do
      local full_command="${ACTION_COMMAND} ${action} ${flag}"
      print_separator "${full_command}"
      ${full_command}
    done
  fi
}

# Determine OS and architecture
detect_os_arch() {
  os=$(uname -s)
  case "$os" in
    Linux)
      if grep -q Microsoft /proc/version 2>/dev/null || grep -q WSL /proc/sys/kernel/osrelease 2>/dev/null; then
        GOOS="windows"
      else
        GOOS="linux"
      fi
      ;;
    Darwin)
      GOOS="darwin"
      ;;
    CYGWIN* | MINGW* | MSYS*)
      GOOS="windows"
      ;;
    *)
      echo "Unknown OS: $os"
      GOOS="unknown"
      ;;
  esac

  arch=$(uname -m)
  case "$arch" in
    arm64)
      GOARCH="arm64"
      ;;
    x86_64)
      GOARCH="amd64"
      ;;
  esac
}

# Build the binary if it doesn't exist
build_binary() {
  if [ ! -f "${COMMAND}" ]; then
    local build_script='scripts/build.sh'
    local text="Running ${build_script}"
    print_color_conditional "GREEN" "${text}"
    echo
    ${build_script}
  fi
}

TESTS="$1"
[ -z "${TESTS}" ] && TESTS="all"

detect_os_arch

if [ "${GOOS}" == "windows" ]; then
  COMMAND="bin/att-fiber-gateway-info_${GOOS}_${GOARCH}.exe"
else
  COMMAND="bin/att-fiber-gateway-info_${GOOS}_${GOARCH}"
fi

build_binary

ACTION_COMMAND="${COMMAND} -action"

if [[ "${TESTS}" == "nologin" || "${TESTS}" == "all" ]]; then
  NO_LOGIN_ACTIONS=(broadband-status device-list fiber-status home-network-status system-information)
  run_commands "" "${NO_LOGIN_ACTIONS[@]}"
fi

if [[ "${TESTS}" == "metrics" || "${TESTS}" == "all" ]]; then
  METRICS_ACTIONS=(broadband-status fiber-status home-network-status)
  run_commands "-metrics" "${METRICS_ACTIONS[@]}"

  run_commands "-allmetrics"
fi

if [[ "${TESTS}" == "login" || "${TESTS}" == "all" ]]; then
  LOGIN_ACTIONS=(ip-allocation nat-check nat-connections nat-destinations nat-sources nat-totals)
  run_commands "" "${LOGIN_ACTIONS[@]}"

  run_commands "-pretty" "nat-connections"
fi

if [[ "${TESTS}" == "reset" || "${TESTS}" == "all" ]]; then
  RESET_ACTIONS=(reset-connection reset-device reset-firewall reset-ip reset-wifi restart-gateway)
  run_commands "-no" "${RESET_ACTIONS[@]}"
fi
