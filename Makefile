.PHONY: all example hello strings clean clean-example clean-hello clean-strings rebuild-example rebuild-hello rebuild-strings rebuild-all

all: example hello strings

example:
	$(MAKE) -C cmd/example all

hello:
	$(MAKE) -C cmd/hello all

strings:
	$(MAKE) -C cmd/strings all

clean: clean-example clean-hello clean-strings

clean-example:
	$(MAKE) -C cmd/example clean

clean-hello:
	$(MAKE) -C cmd/hello clean

clean-strings:
	$(MAKE) -C cmd/strings clean

rebuild-example:
	$(MAKE) -C cmd/example clean all

rebuild-hello:
	$(MAKE) -C cmd/hello clean all

rebuild-strings:
	$(MAKE) -C cmd/strings clean all

rebuild-all:
	$(MAKE) clean all
