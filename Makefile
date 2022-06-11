

.PHONY: update-testdata
update-testdata:
	-RUN_ACTUAL_TEST=1 go test -v -run TestCmp_ -json > testdata/cmp_test_results.json
