function handle_sigterm {
  echo "Localnet is being stopped..."
  ./pause.sh
}

trap handle_sigterm SIGTERM

while true; do
  sleep 1
done
