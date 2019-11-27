CLI=./cetcli # path/to/cetcli

if ! type $CLI > /dev/null; then
  echo 'cetcli' not found
  exit 0
fi
if [ -z "$1" ]; then
  echo 'Usage: sh mykey.sh <suffix>'
  exit 0
fi

printf "finding address with suffix '%s' ...\n" $1
for i in {1..50000}
do
  OUTPUT=`$CLI keys add xxx --dry-run 2>&1`
  ADDR=`echo $OUTPUT | grep -o 'address: [a-z0-9]*' | grep -o 'coinex[0-9a-z]*'`
  printf "\t%d\t%s\r" $i $ADDR
  if [[ $ADDR == *$1 ]]; then
    echo ok $ADDR
    echo $OUTPUT
    break
  fi
done
