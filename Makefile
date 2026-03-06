.PHONY: all check-runtime check-runtime-probes check-runtime-behavior check-app-template check-emulator-smoke example hello strings slices interfaces emptyiface assertions runtimecheck timeprobe smokeapp sysinfo message ipc clean clean-example clean-hello clean-strings clean-slices clean-interfaces clean-emptyiface clean-assertions clean-runtimecheck clean-timeprobe clean-smokeapp clean-sysinfo clean-message clean-ipc rebuild-example rebuild-hello rebuild-strings rebuild-slices rebuild-interfaces rebuild-emptyiface rebuild-assertions rebuild-runtimecheck rebuild-timeprobe rebuild-smokeapp rebuild-sysinfo rebuild-message rebuild-ipc rebuild-all

all: example hello strings slices interfaces emptyiface assertions runtimecheck timeprobe sysinfo message ipc

check-runtime: check-runtime-probes check-runtime-behavior

check-runtime-probes:
	bash ./scripts/check-runtime-probes.sh

check-runtime-behavior:
	bash ./scripts/check-runtime-behavior.sh

check-app-template:
	bash ./scripts/check-app-template.sh

check-emulator-smoke:
	bash ./scripts/check-emulator-smoke.sh

example:
	$(MAKE) -C cmd/example all

hello:
	$(MAKE) -C cmd/hello all

strings:
	$(MAKE) -C cmd/strings all

slices:
	$(MAKE) -C cmd/slices all

interfaces:
	$(MAKE) -C cmd/interfaces all

emptyiface:
	$(MAKE) -C cmd/emptyiface all

assertions:
	$(MAKE) -C cmd/assertions all

runtimecheck:
	$(MAKE) -C cmd/runtimecheck all

timeprobe:
	$(MAKE) -C cmd/timeprobe all

smokeapp:
	$(MAKE) -C cmd/smokeapp all

sysinfo:
	$(MAKE) -C cmd/sysinfo all

message:
	$(MAKE) -C cmd/message all

ipc:
	$(MAKE) -C cmd/ipc all

clean: clean-example clean-hello clean-strings clean-slices clean-interfaces clean-emptyiface clean-assertions clean-runtimecheck clean-timeprobe clean-smokeapp clean-sysinfo clean-message clean-ipc

clean-example:
	$(MAKE) -C cmd/example clean

clean-hello:
	$(MAKE) -C cmd/hello clean

clean-strings:
	$(MAKE) -C cmd/strings clean

clean-slices:
	$(MAKE) -C cmd/slices clean

clean-interfaces:
	$(MAKE) -C cmd/interfaces clean

clean-emptyiface:
	$(MAKE) -C cmd/emptyiface clean

clean-assertions:
	$(MAKE) -C cmd/assertions clean

clean-runtimecheck:
	$(MAKE) -C cmd/runtimecheck clean

clean-timeprobe:
	$(MAKE) -C cmd/timeprobe clean

clean-smokeapp:
	$(MAKE) -C cmd/smokeapp clean

clean-sysinfo:
	$(MAKE) -C cmd/sysinfo clean

clean-message:
	$(MAKE) -C cmd/message clean

clean-ipc:
	$(MAKE) -C cmd/ipc clean

rebuild-example:
	$(MAKE) -C cmd/example clean all

rebuild-hello:
	$(MAKE) -C cmd/hello clean all

rebuild-strings:
	$(MAKE) -C cmd/strings clean all

rebuild-slices:
	$(MAKE) -C cmd/slices clean all

rebuild-interfaces:
	$(MAKE) -C cmd/interfaces clean all

rebuild-emptyiface:
	$(MAKE) -C cmd/emptyiface clean all

rebuild-assertions:
	$(MAKE) -C cmd/assertions clean all

rebuild-runtimecheck:
	$(MAKE) -C cmd/runtimecheck clean all

rebuild-timeprobe:
	$(MAKE) -C cmd/timeprobe clean all

rebuild-smokeapp:
	$(MAKE) -C cmd/smokeapp clean all

rebuild-sysinfo:
	$(MAKE) -C cmd/sysinfo clean all

rebuild-message:
	$(MAKE) -C cmd/message clean all

rebuild-ipc:
	$(MAKE) -C cmd/ipc clean all

rebuild-all:
	$(MAKE) clean all
