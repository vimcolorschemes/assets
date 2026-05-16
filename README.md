# vimcolorschemes assets

Generated visual assets for vimcolorschemes.

## Preview

![vimcolorschemes ANSI art](out/vimcolorschemes.svg)

![v ANSI art](out/v.svg)

## Generate

Requires Go 1.25 or newer.

```sh
make generate
```

Generate a single asset:

```sh
make generate-v
make generate-vimcolorschemes
```

## Test

```sh
make test
```

Generated SVGs are written to `out/`.
