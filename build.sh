#!/bin/bash

# we don't include root by default
BUILT=`readlink -e ${PWD}/Makefile`

build () {
    DIRNAME=`dirname $1`
    echo "built makefiles: ${BUILT}"

    SLASHES=${PWD//[^\/]/}
    MKFILE=`echo "${DIRNAME}/Makefile"`

    for (( n=${#SLASHES}; n>0; --n )); do
        if [ -f $MKFILE ]; then
            echo "found makefile in ${DIRNAME}"
            break
        else
            DIRNAME="${DIRNAME}/.."
            MKFILE=`echo "${DIRNAME}/Makefile"`
        fi
    done
        
    MKFILE_FULL=`readlink -e ${MKFILE}`
    echo "full makefile path: ${MKFILE_FULL}"

    if [[ $BUILT != *"${MKFILE_FULL}"* ]]; then
        echo "build ${DIRNAME} (${MKFILE_FULL})"
        INCLUDE_MAKEFILE=$MKFILE make release
        BUILT=`echo "${BUILT};${MKFILE_FULL}"`
    else
        echo "${MKFILE_FULL} already built, skipping"
    fi
}

# for test
# cat filelist | while read line; do
#     echo "line: ${line}"
#     build $line
#     echo "-"
# done
# exit

echo "range ${TRAVIS_COMMIT_RANGE}"

# walk through each changed file
git diff --name-only $TRAVIS_COMMIT_RANGE | while read line; do
    build $line
    echo "-"
done
