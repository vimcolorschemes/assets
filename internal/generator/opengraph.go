package generator

type openGraphCard struct {
	X       int
	Y       int
	Width   int
	Height  int
	Name    string
	Bg      string
	Border  string
	Palette []string
}

func openGraphCards() []openGraphCard {
	palettes := []openGraphCard{
		{Name: "gruvbox", Bg: "#1d2021", Border: "#a89984", Palette: []string{"#fb4934", "#b8bb26", "#fabd2f", "#83a598", "#d3869b"}},
		{Name: "tokyonight", Bg: "#16161e", Border: "#7aa2f7", Palette: []string{"#bb9af7", "#7dcfff", "#9ece6a", "#e0af68", "#f7768e"}},
		{Name: "kanagawa", Bg: "#1f1f28", Border: "#c8c093", Palette: []string{"#c34043", "#76946a", "#c0a36e", "#7e9cd8", "#957fb8"}},
		{Name: "catppuccin", Bg: "#24273a", Border: "#8aadf4", Palette: []string{"#ed8796", "#a6da95", "#eed49f", "#8bd5ca", "#c6a0f6"}},
		{Name: "everforest", Bg: "#0f1c1e", Border: "#81b29a", Palette: []string{"#e07a5f", "#81b29a", "#f2cc8f", "#83c5be", "#cdb4db"}},
		{Name: "rose-pine", Bg: "#191724", Border: "#ebbcba", Palette: []string{"#eb6f92", "#f6c177", "#9ccfd8", "#c4a7e7", "#31748f"}},
	}

	positions := []struct {
		X int
		Y int
	}{
		{X: 10, Y: -126}, {X: 310, Y: -126}, {X: 610, Y: -126}, {X: 910, Y: -126},
		{X: 10, Y: 60}, {X: 310, Y: 60}, {X: 610, Y: 60}, {X: 910, Y: 60},
		{X: 10, Y: 246}, {X: 310, Y: 246}, {X: 610, Y: 246}, {X: 910, Y: 246},
		{X: 10, Y: 432}, {X: 310, Y: 432}, {X: 610, Y: 432}, {X: 910, Y: 432},
		{X: 10, Y: 618}, {X: 310, Y: 618}, {X: 610, Y: 618}, {X: 910, Y: 618},
	}

	cards := make([]openGraphCard, 0, len(positions))
	for i, position := range positions {
		card := palettes[i%len(palettes)]
		card.X = position.X
		card.Y = position.Y
		card.Width = 280
		card.Height = 166
		cards = append(cards, card)
	}
	return cards
}

func truncateOpenGraphCodeLine(line string, cardWidth int) string {
	maxChars := (cardWidth - 36) / 7
	if len(line) <= maxChars {
		return line
	}
	if maxChars <= 1 {
		return ""
	}
	return line[:maxChars-1] + "."
}

func openGraphCodeLines() []string {
	return []string{
		`$ nvim '+set termguicolors' '+Telescope colorscheme'`,
		`:colorscheme gruvbox-material       " warm contrast`,
		`local api = vim.api                 -- highlight preview`,
		`for _, scheme in ipairs(colorschemes) do`,
		`  preview:render(scheme.name, scheme.repository)`,
		`  api.nvim_set_hl(0, 'Normal', scheme.palette.bg)`,
		`end`,
		`:set background=dark cursorline number`,
		`require('vimcolorschemes').browse({ sort = 'updated' })`,
		`-- compare palettes, screenshots, terminal colors`,
		`highlight String  guifg=#a7c080 | highlight Comment guifg=#7a8478`,
		`$ vimcolorschemes search tokyonight catppuccin kanagawa`,
	}
}
