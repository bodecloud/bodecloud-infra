# VSC Essentials Core - Pack for Developers

[![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/Gydunhn.vsc-essentials-core?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.vsc-essentials-core) [![Installs](https://flat.badgen.net/vs-marketplace/i/Gydunhn.vsc-essentials-core?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.vsc-essentials-core) [![Downloads](https://flat.badgen.net/vs-marketplace/d/Gydunhn.vsc-essentials-core?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.vsc-essentials-core) [![Rating](https://flat.badgen.net/vs-marketplace/rating/Gydunhn.vsc-essentials-core?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.vsc-essentials-core)

This extension pack for Visual Studio Code adds extensions that are convenient and useful for any development (regardless of language). I reserve the right to update the content of the extension pack at my own discretion.

This **Detailed** version of the extension pack is for a series of very specific projects in which I am currently involved; projects with multiple repositories that share the same stack of technologies transversally.

![Preview](https://firebasestorage.googleapis.com/v0/b/vsc-essentials.appspot.com/o/VSC-Essentials-Core%2FVSC-Essentials-Core-Preview.png?alt=media&token=2a85a17f-bf06-4d68-8b6a-cee289451515)

## Reasons

New Purpose of "VSC Essentials Core"

We have given a new purpose to the "VSC Essentials Core" extension to better align with the hierarchy of extension packages used by our team in Visual Studio Code (VSCode) development.

**VSC Essentials Core** will now serve as the core of our extensions. This package will contain the essential extensions that all projects within our team must have, regardless of the framework or programming language used. This set of extensions is fundamental to ensuring a common and consistent base across all our developments.

**[VSC Essentials]**, the original extension, will now feature a broader range of extensions that, while also considered essential, are not mandatory for all developers on the team. These additional extensions may be useful depending on the project or individual preferences but are not part of the required basic set.

This change will allow us to maintain a more organized and efficient structure, facilitating collaboration and maintenance of our projects in VSCode.

See the [CHANGELOG](https://github.com/Gydunhn/VSC-Essentials-Core/blob/HEAD/CHANGELOG.md) for the latest changes

## **settings.json**

It is imperative that the settings be added to settings.json, inside the ".vscode" folder, and that this file be inside Git version control for this extension pack to work correctly.

``` json
{
    /**
     * Core Basic VSC Essentials Config
     */
    // Editor Native Settings
    "editor.fontLigatures": true,
    "editor.stickyScroll.enabled": true,
    "editor.cursorBlinking": "expand",
    "editor.cursorSmoothCaretAnimation": "on",
    "editor.guides.highlightActiveBracketPair": true,
    "editor.bracketPairColorization.enabled": true,
    "editor.bracketPairColorization.independentColorPoolPerBracketType": true,
    "editor.guides.bracketPairs": "active",
    "editor.guides.indentation": true,
    "editor.indentSize": 2,
    "editor.tabCompletion": "on",
    "editor.showDeprecated": true,
    "editor.autoIndent": "full",
    "editor.rulers": [
        140
    ],
    "editor.wordWrap": "bounded",
    "editor.wordWrapColumn": 145,
    // Markdown Related Settings
    "[markdown]": {
        "editor.defaultFormatter": "yzhang.markdown-all-in-one"
    },
    "markdownlint.config": {
        "default": true,
        "MD001": false,
        "MD010": false,
        "MD022": false,
        "MD024": false,
        "MD025": false
    },
    "emojisense.languages": {
        "plaintext": false,
        "markdown": true,
        "json": true,
        "scminput": true
    },
    // Todo Tree Settings
    "todo-tree.tree.showCountsInTree": true,
    "todo-tree.general.statusBar": "top three",
    "todo-tree.general.showIconsInsteadOfTagsInStatusBar": true,
    "todo-tree.general.tags": [
        "TODO",
        "FIXME",
        "FIXIT",
        "FIX",
        "BUG"
    ],
    "todo-tree.general.tagGroups": {
        "FIXME": [
            "FIXME",
            "FIXIT",
            "FIX",
            "BUG",
        ]
    },
    "todo-tree.highlights.customHighlight": {
        "TODO": {
            "gutterIcon": true,
            "icon": "tasklist",
            "iconColour": "#FF8C00",
            "type": "tag",
            "background": "#CF7200",
            "foreground": "#FFFFFF",
            "fontWeight": "bold"
        },
        "FIXME": {
            "gutterIcon": true,
            "icon": "tools",
            "iconColour": "#00FF00",
            "type": "tag",
            "background": "#008000",
            "foreground": "#FFFF00",
            "fontWeight": "bold"
        }
    },
    // Better Comments Settings
    "better-comments.multilineComments": true,
    "better-comments.tags": [
        {
            "tag": "!",
            "color": "#FF2D00",
            "strikethrough": false,
            "underline": false,
            "backgroundColor": "transparent",
            "bold": true,
            "italic": false
        },
        {
            "tag": "?",
            "color": "#3498DB",
            "strikethrough": false,
            "underline": false,
            "backgroundColor": "transparent",
            "bold": false,
            "italic": false
        },
        {
            "tag": "//",
            "color": "#474747",
            "strikethrough": true,
            "underline": false,
            "backgroundColor": "transparent",
            "bold": false,
            "italic": false
        },
        {
            "tag": "todo",
            "color": "#FF8C00",
            "strikethrough": false,
            "underline": false,
            "backgroundColor": "transparent",
            "bold": false,
            "italic": false
        },
        {
            "tag": "fixme",
            "color": "#008000",
            "strikethrough": false,
            "underline": false,
            "backgroundColor": "transparent",
            "bold": false,
            "italic": false
        },
        {
            "tag": "fixit",
            "color": "#008000",
            "strikethrough": false,
            "underline": false,
            "backgroundColor": "transparent",
            "bold": false,
            "italic": false
        },
        {
            "tag": "fix",
            "color": "#008000",
            "strikethrough": false,
            "underline": false,
            "backgroundColor": "transparent",
            "bold": false,
            "italic": false
        },
        {
            "tag": "bug",
            "color": "#008000",
            "strikethrough": false,
            "underline": false,
            "backgroundColor": "transparent",
            "bold": false,
            "italic": false
        },
        {
            "tag": "*",
            "color": "#98C379",
            "strikethrough": false,
            "underline": false,
            "backgroundColor": "transparent",
            "bold": true,
            "italic": false
        }
    ],
    // Terminal in Status Bar Settings
    "terminal-in-status-bar.statusBarAlignment": "right",
    "terminal-in-status-bar.statusBarPriority": 10000,
    // indent-rainbow Settings
    "indentRainbow.ignoreErrorLanguages": [
        "haskell"
    ],
    // Bookmarks Settings
    "bookmarks.saveBookmarksInProject": false,
    "bookmarks.showCommandsInContextMenu": true,
    // Native JSON Settings
    "[json]": {
        "editor.defaultFormatter": "vscode.json-language-features"
    },
    "[jsonc]": {
        "editor.defaultFormatter": "vscode.json-language-features"
    }
}

```

## Note

This extension pack was made from their original [VSC Essentials], which you can find [here].
[This extension] can be found at [open-vsx.org] as well.

## Included

This **Core** extension pack includes the following extensions:

### Productivity and Code Management

| Extension                  | Stats                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Better Comments            | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/aaron-bond.better-comments?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=aaron-bond.better-comments) [![Installs](https://flat.badgen.net/vs-marketplace/i/aaron-bond.better-comments?color=blue)](https://marketplace.visualstudio.com/items?itemName=aaron-bond.better-comments) [![Rating](https://flat.badgen.net/vs-marketplace/rating/aaron-bond.better-comments?color=blue)](https://marketplace.visualstudio.com/items?itemName=aaron-bond.better-comments)                         |
| Todo Tree                  | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/Gruntfuggly.todo-tree?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=Gruntfuggly.todo-tree) [![Installs](https://flat.badgen.net/vs-marketplace/i/Gruntfuggly.todo-tree?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gruntfuggly.todo-tree) [![Rating](https://flat.badgen.net/vs-marketplace/rating/Gruntfuggly.todo-tree?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gruntfuggly.todo-tree)                                                       |
| Terminal in Status Bar     | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/flyfly6.terminal-in-status-bar?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=flyfly6.terminal-in-status-bar) [![Installs](https://flat.badgen.net/vs-marketplace/i/flyfly6.terminal-in-status-bar?color=blue)](https://marketplace.visualstudio.com/items?itemName=flyfly6.terminal-in-status-bar) [![Rating](https://flat.badgen.net/vs-marketplace/rating/flyfly6.terminal-in-status-bar?color=blue)](https://marketplace.visualstudio.com/items?itemName=flyfly6.terminal-in-status-bar) |
| Bookmarks                  | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/alefragnani.Bookmarks?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=alefragnani.Bookmarks) [![Installs](https://flat.badgen.net/vs-marketplace/i/alefragnani.Bookmarks?color=blue)](https://marketplace.visualstudio.com/items?itemName=alefragnani.Bookmarks) [![Rating](https://flat.badgen.net/vs-marketplace/rating/alefragnani.Bookmarks?color=blue)](https://marketplace.visualstudio.com/items?itemName=alefragnani.Bookmarks)                                                       |
| Path Intellisense          | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/christian-kohler.path-intellisense?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=christian-kohler.path-intellisense) [![Installs](https://flat.badgen.net/vs-marketplace/i/christian-kohler.path-intellisense?color=blue)](https://marketplace.visualstudio.com/items?itemName=christian-kohler.path-intellisense) [![Rating](https://flat.badgen.net/vs-marketplace/rating/christian-kohler.path-intellisense?color=blue)](https://marketplace.visualstudio.com/items?itemName=christian-kohler.path-intellisense) |
| Formatting Toggle           | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/tombonnike.vscode-status-bar-format-toggle?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=tombonnike.vscode-status-bar-format-toggle) [![Installs](https://flat.badgen.net/vs-marketplace/i/tombonnike.vscode-status-bar-format-toggle?color=blue)](https://marketplace.visualstudio.com/items?itemName=tombonnike.vscode-status-bar-format-toggle) [![Rating](https://flat.badgen.net/vs-marketplace/rating/tombonnike.vscode-status-bar-format-toggle?color=blue)](https://marketplace.visualstudio.com/items?itemName=tombonnike.vscode-status-bar-format-toggle) |


### Visual Enhancements

| Extension                  | Stats                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| indent-rainbow             | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/oderwat.indent-rainbow?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=oderwat.indent-rainbow) [![Installs](https://flat.badgen.net/vs-marketplace/i/oderwat.indent-rainbow?color=blue)](https://marketplace.visualstudio.com/items?itemName=oderwat.indent-rainbow) [![Rating](https://flat.badgen.net/vs-marketplace/rating/oderwat.indent-rainbow?color=blue)](https://marketplace.visualstudio.com/items?itemName=oderwat.indent-rainbow)                                                 |
| Error Lens                 | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/usernamehw.errorlens?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=usernamehw.errorlens) [![Installs](https://flat.badgen.net/vs-marketplace/i/usernamehw.errorlens?color=blue)](https://marketplace.visualstudio.com/items?itemName=usernamehw.errorlens) [![Rating](https://flat.badgen.net/vs-marketplace/rating/usernamehw.errorlens?color=blue)](https://marketplace.visualstudio.com/items?itemName=usernamehw.errorlens)                                                       |
| Output Colorizer           | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/IBM.output-colorizer?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=IBM.output-colorizer) [![Installs](https://flat.badgen.net/vs-marketplace/i/IBM.output-colorizer?color=blue)](https://marketplace.visualstudio.com/items?itemName=IBM.output-colorizer) [![Rating](https://flat.badgen.net/vs-marketplace/rating/IBM.output-colorizer?color=blue)](https://marketplace.visualstudio.com/items?itemName=IBM.output-colorizer)                                                     |

### Markdown Support

| Extension                  | Stats                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Markdown All in One        | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/yzhang.markdown-all-in-one?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=yzhang.markdown-all-in-one) [![Installs](https://flat.badgen.net/vs-marketplace/i/yzhang.markdown-all-in-one?color=blue)](https://marketplace.visualstudio.com/items?itemName=yzhang.markdown-all-in-one) [![Rating](https://flat.badgen.net/vs-marketplace/rating/yzhang.markdown-all-in-one?color=blue)](https://marketplace.visualstudio.com/items?itemName=yzhang.markdown-all-in-one)                         |
| markdownlint               | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/DavidAnson.vscode-markdownlint?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=DavidAnson.vscode-markdownlint) [![Installs](https://flat.badgen.net/vs-marketplace/i/DavidAnson.vscode-markdownlint?color=blue)](https://marketplace.visualstudio.com/items?itemName=DavidAnson.vscode-markdownlint) [![Rating](https://flat.badgen.net/vs-marketplace/rating/DavidAnson.vscode-markdownlint?color=blue)](https://marketplace.visualstudio.com/items?itemName=DavidAnson.vscode-markdownlint) |
| :emojisense:               | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/bierner.emojisense?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=bierner.emojisense) [![Installs](https://flat.badgen.net/vs-marketplace/i/bierner.emojisense?color=blue)](https://marketplace.visualstudio.com/items?itemName=bierner.emojisense) [![Rating](https://flat.badgen.net/vs-marketplace/rating/bierner.emojisense?color=blue)](https://marketplace.visualstudio.com/items?itemName=bierner.emojisense)                                                                         |
| Markdown Emoji             | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/bierner.markdown-emoji?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-emoji) [![Installs](https://flat.badgen.net/vs-marketplace/i/bierner.markdown-emoji?color=blue)](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-emoji) [![Rating](https://flat.badgen.net/vs-marketplace/rating/bierner.markdown-emoji?color=blue)](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-emoji)                                                 |
| Markdown Checkboxes        | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/bierner.markdown-checkbox?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-checkbox) [![Installs](https://flat.badgen.net/vs-marketplace/i/bierner.markdown-checkbox?color=blue)](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-checkbox) [![Rating](https://flat.badgen.net/vs-marketplace/rating/bierner.markdown-checkbox?color=blue)](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-checkbox)                               |

[vsc essentials]: https://marketplace.visualstudio.com/items?itemName=Gydunhn.vsc-essentials
[here]: https://marketplace.visualstudio.com/items?itemName=Gydunhn.vsc-essentials
[This extension]: https://open-vsx.org/extension/Gydunhn/vsc-essentials-core
[open-vsx.org]: https://open-vsx.org/
