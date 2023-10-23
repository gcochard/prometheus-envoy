BIN=prometheus-envoy
INS_TARGET=/usr/local/bin/$(BIN)
SVC_TARGET=/etc/systemd/system/prometheus-envoy.service
TARGET=cmd/prometheus-envoy/$(BIN)

$(TARGET): cmd/**/*.go pkg/*.go
	cd cmd/prometheus-envoy && go build

.PHONY: clean
clean:
	rm $(TARGET)

.PHONY: install
install: $(INS_TARGET)

$(INS_TARGET): $(TARGET)
	cp $(TARGET) /usr/local/bin/$(BIN)

.PHONY: install-svc
install-svc: $(SVC_TARGET) $(INS_TARGET)

$(SVC_TARGET): init/prometheus-envoy.service
	cp init/prometheus-envoy.service $(SVC_TARGET)
	systemctl daemon-reload
