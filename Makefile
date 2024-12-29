

.PHONY: update-testdata
update-testdata:
	-RUN_ACTUAL_TEST=1 go test -v -run TestCmp_ -json > testdata/cmp_test_results.json
	-perl -i -pe 's/"Time":".*?"/"Time":"2022-06-11T00:00:00.0Z"/' testdata/cmp_test_results.json
	-perl -i -pe 's/[0-9]+\.[0-9]+s/0.01s/g' testdata/cmp_test_results.json
	-perl -i -pe 's/"Elapsed":[0-9]+\.[0-9]+/"Elapsed":1.00/' testdata/cmp_test_results.json
