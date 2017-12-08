#!/bin/bash

upsearch () {
  slashes=${PWD//[^\/]/}
  directory="$PWD"
  for (( n=${#slashes}; n>0; --n ))
  do
    test -e "$directory/$1" && echo "$directory/$1" && return
    directory="$directory/.."
  done
}

echo "range ${TRAVIS_COMMIT_RANGE}"
git diff --name-only $TRAVIS_COMMIT_RANGE | while read line; do
  DIRNAME=`dirname ${line}`;
  if [[ "$DIRNAME" != "$DIRNAME_OLD" ]]; then
    if [[ $BUILT != *"${DIRNAME}"* ]]; then
      MKFILE=`echo "${DIRNAME}/Makefile"`;
      if [ -f $MKFILE ]; then
        echo "build ${DIRNAME}";
        INCLUDE_MAKEFILE=$MKFILE make release;
        BUILT=`echo "${BUILT};${DIRNAME}"`;
      fi;
    fi;
  fi;
  DIRNAME_OLD=$DIRNAME;
done
