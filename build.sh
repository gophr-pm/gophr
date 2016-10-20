#!/bin/bash
EXIT_STATUS=0

echo "mode: count" > coverage-all.out
for pkg in $( find . -name *.go -print0 | xargs -0 -n1 dirname | sort --unique ); do
    if [ $pkg == "./gophr-test" ]; then continue; fi; # integration test should be skipped 
    go test -coverprofile=coverage.out -covermode=count $pkg;
    if [ $? -ne 0 ]; then EXIT_STATUS=1; fi;
    tail -n +2 coverage.out >> coverage-all.out;
done
echo "exiting with status $EXIT_STATUS";
exit $EXIT_STATUS;