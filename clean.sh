#!/usr/bin/env bash
find . | grep -E  "(log$|data$|/.db$|./log$)"  |  xargs rm -rf