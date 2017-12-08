#!/bin/bash

upsearch () {
  slashes=${PWD//[^\/]/}
  directory="$PWD"
  for (( n=${#slashes}; n>0; --n ))
  do
    # test -e "$directory/$1" && echo "$directory/$1" && return
    echo "hello ${directory}, ${#slashes}"
    directory="$directory/.."
  done
}

# upsearch
# exit 0

build () {
  DIRNAME=`dirname $1`
  if [[ "$DIRNAME" != "$DIRNAME_OLD" ]]; then
    if [[ $BUILT != *"${DIRNAME}"* ]]; then
      SLASHES=${PWD//[^\/]/}
      MKFILE=`echo "${DIRNAME}/Makefile"`
      if [ -f $MKFILE ]; then
        echo "build ${DIRNAME}"
        INCLUDE_MAKEFILE=$MKFILE make release
        BUILT=`echo "${BUILT};${DIRNAME}"`
      fi
    fi
  fi
  DIRNAME_OLD=$DIRNAME
}

echo "range ${TRAVIS_COMMIT_RANGE}"

# walk through each changed file
git diff --name-only $TRAVIS_COMMIT_RANGE | while read line; do
    build $line;
done
