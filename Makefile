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
	@echo "Syncing embedded files from kiwi_nexus..."
	cp -r ../kiwi_nexus/.claude/skills/spec-it/SKILL.md embedded/skills/spec-it/SKILL.md
	cp -r ../kiwi_nexus/.claude/skills/execute-spec/SKILL.md embedded/skills/execute-spec/SKILL.md
	cp -r ../kiwi_nexus/.claude/skills/commit/SKILL.md embedded/skills/commit/SKILL.md
	cp -r ../kiwi_nexus/.claude/skills/pr/SKILL.md embedded/skills/pr/SKILL.md
	cp ../kiwi_nexus/.claude/agents/merge-wave.md embedded/agents/merge-wave.md
	cp ../kiwi_nexus/.claude/agents/post-spec-docs.md embedded/agents/post-spec-docs.md
	cp ../kiwi_nexus/.claude/agents/knowledge-gc.md embedded/agents/knowledge-gc.md
	cp ../kiwi_nexus/.claude/rules/aryflow.md embedded/rules/aryflow.md
	cp ../kiwi_nexus/.claude/hooks/aryflow-session-start.sh embedded/hooks/aryflow-session-start.sh
	cp ../kiwi_nexus/.claude/hooks/aryflow-stop.sh embedded/hooks/aryflow-stop.sh
	cp ../kiwi_nexus/.claude/hooks/aryflow-subagent-stop.sh embedded/hooks/aryflow-subagent-stop.sh
	@echo "Done."

clean:
	rm -f $(BINARY) $(BINARY)-*

.PHONY: build test build-all sync-embedded clean
