export $(grep -v '^#' .env | xargs)
goose up
unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)