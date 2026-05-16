.PHONY: generate generate-v generate-vimcolorschemes test

generate:
	go run ./cmd/generate-assets

generate-v:
	go run ./cmd/generate-assets v

generate-vimcolorschemes:
	go run ./cmd/generate-assets vimcolorschemes

test:
	go test ./...
