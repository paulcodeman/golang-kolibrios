.PHONY: all examples apps check-runtime check-runtime-probes check-runtime-behavior check-app-template check-emulator-smoke check-diagnostics window runtime time system input ipc files path strings bytes io os fmt console diag smokeapp clean clean-examples clean-apps clean-window clean-runtime clean-time clean-system clean-input clean-ipc clean-files clean-path clean-strings clean-bytes clean-io clean-os clean-fmt clean-console clean-diag clean-smokeapp rebuild-window rebuild-runtime rebuild-time rebuild-system rebuild-input rebuild-ipc rebuild-files rebuild-path rebuild-strings rebuild-bytes rebuild-io rebuild-os rebuild-fmt rebuild-console rebuild-diag rebuild-smokeapp rebuild-all

all: examples apps smokeapp

examples: window runtime time system input ipc files path strings bytes io os fmt console

apps: diag

check-runtime: check-runtime-probes check-runtime-behavior

check-runtime-probes:
	bash ./scripts/check-runtime-probes.sh

check-runtime-behavior:
	bash ./scripts/check-runtime-behavior.sh

check-app-template:
	bash ./scripts/check-app-template.sh

check-emulator-smoke:
	bash ./scripts/check-emulator-smoke.sh

check-diagnostics:
	bash ./scripts/check-diagnostics.sh

window:
	$(MAKE) -C examples/window all

runtime:
	$(MAKE) -C examples/runtime all

time:
	$(MAKE) -C examples/time all

system:
	$(MAKE) -C examples/system all

input:
	$(MAKE) -C examples/input all

ipc:
	$(MAKE) -C examples/ipc all

files:
	$(MAKE) -C examples/files all

path:
	$(MAKE) -C examples/path all

strings:
	$(MAKE) -C examples/strings all

bytes:
	$(MAKE) -C examples/bytes all

io:
	$(MAKE) -C examples/io all

os:
	$(MAKE) -C examples/os all

fmt:
	$(MAKE) -C examples/fmt all

console:
	$(MAKE) -C examples/console all

diag:
	$(MAKE) -C apps/diag all

smokeapp:
	$(MAKE) -C tests/smokeapp all

clean: clean-examples clean-apps clean-smokeapp

clean-examples: clean-window clean-runtime clean-time clean-system clean-input clean-ipc clean-files clean-path clean-strings clean-bytes clean-io clean-os clean-fmt clean-console

clean-apps: clean-diag

clean-window:
	$(MAKE) -C examples/window clean

clean-runtime:
	$(MAKE) -C examples/runtime clean

clean-time:
	$(MAKE) -C examples/time clean

clean-system:
	$(MAKE) -C examples/system clean

clean-input:
	$(MAKE) -C examples/input clean

clean-ipc:
	$(MAKE) -C examples/ipc clean

clean-files:
	$(MAKE) -C examples/files clean

clean-path:
	$(MAKE) -C examples/path clean

clean-strings:
	$(MAKE) -C examples/strings clean

clean-bytes:
	$(MAKE) -C examples/bytes clean

clean-io:
	$(MAKE) -C examples/io clean

clean-os:
	$(MAKE) -C examples/os clean

clean-fmt:
	$(MAKE) -C examples/fmt clean

clean-console:
	$(MAKE) -C examples/console clean

clean-diag:
	$(MAKE) -C apps/diag clean

clean-smokeapp:
	$(MAKE) -C tests/smokeapp clean

rebuild-window:
	$(MAKE) -C examples/window clean all

rebuild-runtime:
	$(MAKE) -C examples/runtime clean all

rebuild-time:
	$(MAKE) -C examples/time clean all

rebuild-system:
	$(MAKE) -C examples/system clean all

rebuild-input:
	$(MAKE) -C examples/input clean all

rebuild-ipc:
	$(MAKE) -C examples/ipc clean all

rebuild-files:
	$(MAKE) -C examples/files clean all

rebuild-path:
	$(MAKE) -C examples/path clean all

rebuild-strings:
	$(MAKE) -C examples/strings clean all

rebuild-bytes:
	$(MAKE) -C examples/bytes clean all

rebuild-io:
	$(MAKE) -C examples/io clean all

rebuild-os:
	$(MAKE) -C examples/os clean all

rebuild-fmt:
	$(MAKE) -C examples/fmt clean all

rebuild-console:
	$(MAKE) -C examples/console clean all

rebuild-diag:
	$(MAKE) -C apps/diag clean all

rebuild-smokeapp:
	$(MAKE) -C tests/smokeapp clean all

rebuild-all:
	$(MAKE) clean all
