#!/bin/bash

# arg1 = relative directory (or go package) to test.
cd "$1"

RUN_ACTUAL_TEST=1 go test -v -trimpath -run TestCmp_ -json |
	# Replace time & timestamps
	perl -pe 's/"Time":".*?"/"Time":"2022-06-11T00:00:00.0Z"/' |
	perl -pe 's/[0-9]+\.[0-9]+s/0.01s/g' |
	perl -pe 's/"Elapsed":[0-9]+\.[0-9]+/"Elapsed":1.00/' |

	# dummy cat for consistently terminating above commands with |
	cat