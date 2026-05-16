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

Generated SVG, PNG, and WebP files are written to `out/`, including `*-borderless` variants without the TUI border.

## Theme

Edit `theme.toml` to change the generated images without changing Go code. It controls the font, background, gradient colors, text spacing, shadow, and per-asset padding/offsets.
