#!/bin/sh
APP=fsdd
TESTCASE=/tmp/test.$APP.$(date +%s)
APPOPT="$TESTCASE --verbose --hard-link --clean-symlinks"
GO=$GOROOT/bin/go
if [ $BSD_DEV ]; then . /etc/.bsdconf; fi
prep_payload() {
	if [ -e "$TESTCASE" ]; then rm -rf $TESTCASE; fi
	mkdir -p $TESTCASE && cd $TESTCASE && (
		PAYLOAD="IM 23 BYTES OF CONTENT"
		echo $PAYLOAD > original
		echo $PAYLOAD > same-content-individual-file
		ln -f original hardlink-ok
		ln -fs original symlink-ok
		ln -fs not-available symlink-fail
	)
}
echo
if [ -x $GO ]; then
	echo
	echo "#### go run from source"
	prep_payload
	(cd $BSD_DEV/$APP/APP/$APP && $GO run main.go $APPOPT)
	echo "# verify"
	(cd $BSD_DEV/$APP/APP/$APP && $GO run main.go $TESTCASE --verbose)
fi
if [ -x /usr/bin/$APP ]; then
	echo
	echo "#### compared to: /usr/bin/$APP output"
	prep_payload
	/usr/bin/$APP $APPOPT
	echo "# verify"
	/usr/bin/$APP $TESTCASE --verbose
fi
echo
if [ -e "$TESTCASE" ]; then rm -rf $TESTCASE; fi
