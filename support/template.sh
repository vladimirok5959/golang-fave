#!/bin/bash

find ./engine/assets/template/ ! -name 'template.go' -type f -exec rm -f {} +

echo "package template" > ./engine/assets/template/template.go
echo "" >> ./engine/assets/template/template.go
echo "var AllData = map[string][]byte{" >> ./engine/assets/template/template.go

for FILE_FULL in $(find ./hosts/localhost/template/ | grep -v 'hosts/localhost/template/$' | grep -v '.keep'); do
	FILE_BASE="$(basename -- $FILE_FULL)"
	FILE_GO_BASE="${FILE_BASE//[\-\.]/_}_file.go"
	FILE_GO_FULL="./engine/assets/template/${FILE_GO_BASE}"
	GO_VAR_NAME="${FILE_GO_BASE}"
	GO_VAR_NAME=$(echo "$GO_VAR_NAME" | sed -E 's/^([a-zA-Z]{1})/\U\1/g')
	GO_VAR_NAME=$(echo "$GO_VAR_NAME" | sed -E 's/(_)([a-zA-Z]{1})/\U\2/g')
	GO_VAR_NAME=$(echo "$GO_VAR_NAME" | sed -e 's/\.go$//g')
	GO_VAR_NAME="Var${GO_VAR_NAME}"
	FILE_CONTENT=$(cat ${FILE_FULL})
	FILE_CONTENT=$(echo "$FILE_CONTENT" | sed -E 's/([`]+)/` + "\1" + `/g')

	# Write target file
	echo "package template" > ${FILE_GO_FULL}
	echo "" >> ${FILE_GO_FULL}
	echo "var ${GO_VAR_NAME} = []byte(\`${FILE_CONTENT}\`)" >> ${FILE_GO_FULL}

	# Add files to hash
	echo "	\"${FILE_BASE}\": ${GO_VAR_NAME}," >> ./engine/assets/template/template.go
done

echo "}" >> ./engine/assets/template/template.go
