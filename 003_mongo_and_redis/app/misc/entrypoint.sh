#!/bin/sh

if [ "$#" -eq 0 ]; then
  echo "Error! No arguments."
  echo ""
  echo "Usage:"
  echo "  run - to start the application"
  echo "  <custom command> - to run custom command"
  echo ""
  exit 2
fi

case $@ in
  "run")
    env
    ./wait-for-it.sh $BP_MONGO_SERVICE_TO_CHECK -t 600
    ./wait-for-it.sh $BP_REDIS_SERVICE_TO_CHECK -t 600
    ./main
    ;;
  *)
    echo "Argument unknown. Falling back to execute all the supplied arguments as command..."
    echo "$@"
    echo ""

    set -e
    exec "$@"
    ;;
esac
