VERSION ?= 0.1.0
BINARY = aryflow

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/aryflow

test:
	go test ./... -v

build-all:
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)-darwin-arm64 ./cmd/aryflow
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)-darwin-amd64 ./cmd/aryflow
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)-linux-amd64 ./cmd/aryflow

sync-embedded:
	@echo "Syncing .claude/ → embedded/ (source of truth: .claude/)"
	cp .claude/skills/spec-it/SKILL.md embedded/skills/spec-it/SKILL.md
	cp .claude/skills/execute-spec/SKILL.md embedded/skills/execute-spec/SKILL.md
	cp .claude/skills/commit/SKILL.md embedded/skills/commit/SKILL.md
	cp .claude/skills/pr/SKILL.md embedded/skills/pr/SKILL.md
	cp .claude/agents/merge-wave.md embedded/agents/merge-wave.md
	cp .claude/agents/post-spec-docs.md embedded/agents/post-spec-docs.md
	cp .claude/agents/knowledge-gc.md embedded/agents/knowledge-gc.md
	cp .claude/rules/aryflow.md embedded/rules/aryflow.md
	cp .claude/hooks/aryflow-session-start.sh embedded/hooks/aryflow-session-start.sh
	cp .claude/hooks/aryflow-stop.sh embedded/hooks/aryflow-stop.sh
	cp .claude/hooks/aryflow-subagent-stop.sh embedded/hooks/aryflow-subagent-stop.sh
	cp .claude/hooks/aryflow-statusline.js embedded/hooks/aryflow-statusline.js
	cp .claude/hooks/aryflow-context-monitor.js embedded/hooks/aryflow-context-monitor.js
	@echo "Done."

sync-check:
	@echo "Checking .claude/ ↔ embedded/ sync..."
	@FAIL=0; \
	for f in skills/spec-it/SKILL.md skills/execute-spec/SKILL.md skills/commit/SKILL.md skills/pr/SKILL.md \
	         agents/merge-wave.md agents/post-spec-docs.md agents/knowledge-gc.md \
	         rules/aryflow.md \
	         hooks/aryflow-session-start.sh hooks/aryflow-stop.sh hooks/aryflow-subagent-stop.sh \
	         hooks/aryflow-statusline.js hooks/aryflow-context-monitor.js; do \
	  if ! diff -q ".claude/$$f" "embedded/$$f" > /dev/null 2>&1; then \
	    echo "  MISMATCH: $$f"; FAIL=1; \
	  fi; \
	done; \
	if [ "$$FAIL" = "0" ]; then echo "  All files synced."; else echo "  Run 'make sync-embedded' to fix."; exit 1; fi

clean:
	rm -f $(BINARY) $(BINARY)-*

.PHONY: build test build-all sync-embedded sync-check clean
