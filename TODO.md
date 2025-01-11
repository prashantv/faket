## TODO

## Core functionality

* Support for Nested `t.Cleanup` function calls
  * With support for `t.Skip` inside the nested cleanup calls.


## Cosmetic

 * Naming and API surface for getting logs, `testingLogOutput()`, `Logs()` and `LogsList()`
   and ideally some way to get detailed information for the log (structured message, function, caller, etc)

## Feature Requests

* Return caller information for logs.