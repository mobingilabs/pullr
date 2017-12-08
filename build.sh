#!/bin/bash

# we don't include root by default
BUILT=`readlink -e ${PWD}/Makefile`

build () {
    DIRNAME=`dirname $1`
    SLASHES=${PWD//[^\/]/}
    MKFILE=`echo "${DIRNAME}/Makefile"`

    # try walking up the path until we find a makefile
    for (( n=${#SLASHES}; n>0; --n )); do
        if [ -f $MKFILE ]; then
            echo "Found Makefile in ${DIRNAME}"
            break
        else
            DIRNAME="${DIRNAME}/.."
            MKFILE=`echo "${DIRNAME}/Makefile"`
        fi
    done
        
    MKFILE_FULL=`readlink -e ${MKFILE}`

    if [[ $BUILT != *"${MKFILE_FULL}"* ]]; then
        echo "Build ${DIRNAME} (${MKFILE_FULL})"
        INCLUDE_MAKEFILE=$MKFILE make release
        BUILT=`echo "${BUILT};${MKFILE_FULL}"`
    else
        echo "Skip ${MKFILE_FULL} (already built, or root)"
    fi
}

echo "Range ${TRAVIS_COMMIT_RANGE}"

# walk through each changed file within the range
git diff --name-only $TRAVIS_COMMIT_RANGE | while read line; do
    if [[ $line != vendor* ]]; then
        build $line
        echo "-"
    fi
done
