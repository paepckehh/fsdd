#!/bin/sh
HARDLINK_LIMIT=32000 && TARGET=/tmp/inode-test && LOOP=true
rm -rf $TARGET && mkdir $TARGET && cd $TARGET && touch file.node.zero
echo "node done:"
while [ $LOOP ]; do
	LOOPS=$((LOOPS + 1))
	(
		touch file.$LOOPS && ln -f file.node.zero file.$LOOPS
		CONTENT=$(cat file.$LOOPS)
		if [ "$CONTENT" != "" ]; then
			echo "\n fail node $LOOPS: $CONTENT"
		fi
	) &
	if [ $LOOPS -gt $HARDLINK_LIMIT ]; then
		echo "\n test passed on [$TARGET]"
		break
	fi
	echo -n " $LOOPS"
done
