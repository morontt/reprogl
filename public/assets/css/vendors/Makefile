# Generation all Pygments styles
# Thanks https://github.com/richleland/pygments-css

STYLES = default
STYLES += emacs
STYLES += friendly
STYLES += friendly_grayscale
STYLES += colorful
STYLES += autumn
STYLES += murphy
STYLES += manni
STYLES += material
STYLES += monokai
STYLES += perldoc
STYLES += pastie
STYLES += borland
STYLES += trac
STYLES += native
STYLES += fruity
STYLES += bw
STYLES += vim
STYLES += vs
STYLES += tango
STYLES += rrt
STYLES += xcode
STYLES += igor
STYLES += paraiso-light
STYLES += paraiso-dark
STYLES += lovelace
STYLES += algol
STYLES += algol_nu
STYLES += arduino
STYLES += rainbow_dash
STYLES += abap
STYLES += solarized-dark
STYLES += solarized-light
STYLES += sas
STYLES += stata
STYLES += stata-light
STYLES += stata-dark
STYLES += inkpot
STYLES += zenburn
STYLES += gruvbox-dark
STYLES += gruvbox-light
STYLES += dracula
STYLES += one-dark
STYLES += lilypond

gen_html = pygmentize -S $(style) -f html -a .highlight > pygments/$(style).css;

cssfiles:
	@$(foreach style, $(STYLES), $(gen_html))
