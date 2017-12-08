# This included makefile should define the 'custom' target rule which is called here.
include $(INCLUDE_MAKEFILE)

.PHONY: release prune
release: custom prune

prune:
	@docker system prune -f
