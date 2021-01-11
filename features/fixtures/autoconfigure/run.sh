#!/usr/bin/env bash

# SIGTERM or SIGINT trapped (likely SIGTERM from docker), pass it onto app
# process
function _term_or_init {
  kill -TERM "$APP_PID" 2>/dev/null
  wait $APP_PID
}

# The bugsnag notifier monitor process needs at least 300ms, in order to ensure
# that it can send its notify
function _exit {
  sleep 1
}

trap _term_or_init SIGTERM SIGINT
trap _exit EXIT

PROC="${@:1}"
$PROC &

# Wait on the app process to ensure that this script is able to trap the SIGTERM
# signal
APP_PID=$!
wait $APP_PID
