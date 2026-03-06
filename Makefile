.PHONY: all example clean clean-example rebuild-example

all: example

example:
	$(MAKE) -C cmd/example all

clean: clean-example

clean-example:
	$(MAKE) -C cmd/example clean

rebuild-example:
	$(MAKE) -C cmd/example clean all
