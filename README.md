# bubbles

[![Go Reference](https://pkg.go.dev/badge/github.com/shawalli/bubbles.svg)](https://pkg.go.dev/github.com/shawalli/bubbles)

bubbles is a collection of TUI elements for [Bubble Tea](https://github.com/charmbracelet/bubbletea) applications.

## Radio

![Grouped radio button demo](assets/radio-grouped.gif)

`radio` simplifies the management of radio buttons, which may be presented vertically or horizontally.

* [Example code, basic radio buttons](examples/radio/simple/main.go)
* [Example code, pillbox-style buttons](examples/tabs/pillbox/main.go)
* [Example code, grouped buttons](examples/tabs/resizeable/main.go)

## Tabs

![Wraparound tab demo](assets/tabs-wraparound.gif)

`tabs` is a remix on the [Bubble Tea tabs example](https://github.com/charmbracelet/bubbletea/tree/main/examples/tabs).
It abstracts away the tab logic into a separate model and provides management of tab-content.

* [Example code, basic tabs](examples/tabs/simple/main.go)
* [Example code, wraparound tabs](examples/tabs/wraparound/main.go)
* [Example code, resizeable tabs](examples/tabs/resizeable/main.go)

## License

This project is licensed under the terms of the MIT license.
