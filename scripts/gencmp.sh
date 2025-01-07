#!/bin/bash

# arg1 = relative directory (or go package) to test.
cd "$1"

RUN_ACTUAL_TEST=1 go test -v -trimpath -run TestCmp_ -json |
	# Replace time & timestamps
	perl -pe 's/"Time":".*?"/"Time":"2022-06-11T00:00:00.0Z"/' |
	perl -pe 's/[0-9]+\.[0-9]+s/0.01s/g' |
	perl -pe 's/"Elapsed":[0-9]+\.[0-9]+/"Elapsed":1.00/' |

	# Replace pointers & goroutine IDs 
	perl -pe 's/0x[0-9a-fA-F]+/0x0000/g' |
	perl -pe 's/goroutine [0-9]+/goroutine 1/g' |

	# Replace line numbers in stdlib files (matches */*.go)
	perl -pe 's/(\\t[a-z]+\/[a-z]+.go):[0-9]+/\1:1/g' |

	# dummy cat for consistently terminating above commands with |
	cat
