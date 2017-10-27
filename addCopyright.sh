#!/usr/bin/env bash

# Copyright 2017, Dell EMC, Inc.

COPYRIGHT=" Copyright 2017, Dell EMC, Inc."

declare -A COMMENT=()

COMMENT["go_comment"]="//${COPYRIGHT}"
COMMENT["yaml_comment"]="#${COPYRIGHT}"
COMMENT["yml_comment"]="#${COPYRIGHT}"
COMMENT["js_comment"]="//${COPYRIGHT}"
COMMENT["css_comment"]="/*${COPYRIGHT}*/"
COMMENT["html_comment"]="<!--${COPYRIGHT}-->"
COMMENT["lock_comment"]="#${COPYRIGHT}"
COMMENT["md_comment"]="[//]: # (${COPYRIGHT})"
COMMENT["sh_comment"]="#${COPYRIGHT}"
COMMENT["swagger-codegen-ignore_comment"]="#${COPYRIGHT}"
COMMENT["gitignore_comment"]="#${COPYRIGHT}"


GENERATED_FOLDERS="client cmd models restapi"

FILE_TYPES=`find . -type f -name '*.*' | sed 's|.*\.||' | sort -u`

IGNORED_TYPES="xml sh sample png idx iml pack css js map"

function part_of_list(){
    echo $1 | grep -w $2
}

for folder in $GENERATED_FOLDERS; do
    pushd $folder

    for fileType in $FILE_TYPES; do
        if [[ ! $(part_of_list "$IGNORED_TYPES" $fileType) ]]; then
            echo
            echo "Processing file type: $fileType"
            for file in $(find . -type f -name "*.$fileType"); do
            if grep -Fq "$COPYRIGHT" $file
                then
                    echo "  File [$file] Already has copyrights"
                else
                    comment=${COMMENT[${fileType}_comment]}
                    echo "  Adding comment [$comment] to file [$file]"
                    sed -i "1s;^;$comment\n\n;" $file
                fi
            done;
        fi

    done;

    popd
done;



