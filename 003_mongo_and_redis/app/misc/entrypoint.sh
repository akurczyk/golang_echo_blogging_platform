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
    ./wait-for-it.sh $MONGO_HOST:$MONGO_PORT -t 600
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
