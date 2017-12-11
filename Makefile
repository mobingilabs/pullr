# This included makefile should define the 'custom' target rule which is called here.
include $(INCLUDE_MAKEFILE)

.PHONY: pre release

release: pre custom 

pre:
	@./pre.sh
