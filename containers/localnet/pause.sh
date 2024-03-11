sendSignalToProcessByPort() {
  local port=$1
  local pid=$(lsof -t -i:$port)
  local signal=$2

  if [ -n "$pid" ]
  then
    echo "Sending signal $signal to process $pid..."
    kill -n $signal $pid
  fi
}

sudo python3 /home/ubuntu/add_result.py "state" "paused"

pause_signal=2

sendSignalToProcessByPort 21500 $pause_signal
sendSignalToProcessByPort 21501 $pause_signal
sendSignalToProcessByPort 21502 $pause_signal
sendSignalToProcessByPort 21503 $pause_signal

sendSignalToProcessByPort 9999 $pause_signal

sendSignalToProcessByPort 7950 $pause_signal
