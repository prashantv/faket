SHELL := /bin/bash

# Find DIR/testdata/cmp_test_results.json files, and map to ./DIR
CMP_TEST_DIRS := $(shell \
	git ls-files '**/cmp_test_results.json' | \
	sed -e 's|^|./|' -e 's|testdata/cmp_test_results.json||' \
)

# Note(prashant): Helpful alias since I can never remember `-B`
.PHONY: force-update-testdata
force-update-testdata:
	make -B update-testdata

# cmp_rule will append actions to these rules.
.PHONY: update-testdata
update-testdata::

.PHONY: diff-testdata
diff-testdata::

define cmp_rule
update-testdata:: $1/testdata/cmp_test_results.json
$1/testdata/cmp_test_results.json: $$(wildcard $1/*_test.go) $$(wildcard ./internal/cmptest/*.go)
	./scripts/gencmp.sh "$1" > "$$@"

diff-testdata::
	diff $1/testdata/cmp_test_results.json <(./scripts/gencmp.sh "$1")
endef
$(foreach d,$(CMP_TEST_DIRS),$(eval $(call cmp_rule,$d)))

