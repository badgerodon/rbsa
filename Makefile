include $(GOROOT)/src/Make.inc

TARG=github.com/badgerodon/rbsa
GOFILES=\
  rbsa.go\
  analyze.go\
  cache.go\

include $(GOROOT)/src/Make.pkg
