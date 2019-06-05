#!/bin/bash

LAST_NUM_STR=`find ./support/migrate/ -maxdepth 1 -name '*.go' | sort -t_ -nk2,2 | tail -n1 | awk -F/ '{print $NF}' | awk -F. '{print $NR}'`

NEXT_NUM_INT=$(($LAST_NUM_STR + 1))
NEXT_NUM_STR=`printf %09d $NEXT_NUM_INT`
TARGET_FILE="./support/migrate/${NEXT_NUM_STR}.go"

# Create new migration file
echo "package migrate" > $TARGET_FILE
echo "" >> $TARGET_FILE
echo "import (" >> $TARGET_FILE
echo "	\"golang-fave/engine/sqlw\"" >> $TARGET_FILE
echo ")" >> $TARGET_FILE
echo "" >> $TARGET_FILE
echo "func Migrate_${NEXT_NUM_STR}(db *sqlw.DB) error {" >> $TARGET_FILE
echo "	return nil" >> $TARGET_FILE
echo "}" >> $TARGET_FILE

# Update list
LIST_FILE="./support/migrate/000000001.go"
echo "package migrate" > $LIST_FILE
echo "" >> $LIST_FILE
echo "import (" >> $LIST_FILE
echo "	\"golang-fave/engine/sqlw\"" >> $LIST_FILE
echo ")" >> $LIST_FILE
echo "" >> $LIST_FILE
echo "var Migrations = map[string]func(*sqlw.DB) error{" >> $LIST_FILE
echo "	\"000000000\": nil," >> $LIST_FILE
echo "	\"000000001\": nil," >> $LIST_FILE
for i in `ls ./support/migrate/*.go | sort -V`; do
	IFILE=`echo "$i" | awk -F/ '{print $NF}' | awk -F. '{print $NR}'`
	if [ $IFILE == "000000000" ] || [ $IFILE == "000000001" ]; then
		continue
	fi
	echo "	\"${IFILE}\": Migrate_${IFILE}," >> $LIST_FILE
done
echo "}" >> $LIST_FILE

echo "New migration file created: ${TARGET_FILE}"
