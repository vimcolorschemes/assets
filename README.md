# vimcolorschemes assets

Generated visual assets for vimcolorschemes.

## Preview

![vimcolorschemes ANSI art](out/vimcolorschemes/vimcolorschemes.svg)

![v ANSI art](out/v/v.svg)

## Generate

Requires Go 1.25 or newer.

```sh
make generate
```

Files are generated in `./out`.

## Test

```sh
make test
```

Generated SVG and PNG files are written to per-asset directories under `out/`, including dark/default, `*-light`, and `*-transparent` variants. Non-transparent variants also include WebP files. Each background variant has a bordered and `*-borderless` form.

## Theme

Edit `theme.toml` to change the generated images without changing Go code. It controls the font, dark and light backgrounds, gradient colors, text spacing, shadow, and per-asset padding/offsets.
