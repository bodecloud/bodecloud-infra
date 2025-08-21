# JavaScript Essentials - Extension Pack for Visual Studio Code

[![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/Gydunhn.javascript-essentials?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.javascript-essentials) [![Installs](https://flat.badgen.net/vs-marketplace/i/Gydunhn.javascript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.javascript-essentials) [![Downloads](https://flat.badgen.net/vs-marketplace/d/Gydunhn.javascript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.javascript-essentials) [![Rating](https://flat.badgen.net/vs-marketplace/rating/Gydunhn.javascript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.javascript-essentials)

This extension pack for Visual Studio Code adds extensions that are useful for JavaScript projects. I reserve the right to update the extensions pack contents up to my own discretion. This extension is for my personal use, I think it's great if it works for other people too.

## Reasons

The [JavaScript Essentials extension pack] was made to automate and standardize the installation phase of the essential JavaScript extensions for Visual Studio Code every time a new member joins the team, or one of them restores a laptop, or exchanges it for a new one.

See the [CHANGELOG](https://github.com/Gydunhn/Javascript-Essentials/blob/HEAD/CHANGELOG.md) for the latest changes

## **settings.json**

It is strongly recommended that these settings be used in your workspace. You must copy and paste them, and if you need to adjust something you will already know where to do it.

``` json
{
    /**
     * JavaScript Essentials Config
     */
    // General JavaScript Settings
    "[javascript]": {
        "editor.defaultFormatter": "vscode.typescript-language-features"
    },
    "javascript.format.enable": true,
    "javascript.format.semicolons": "insert",
    "javascript.preferences.quoteStyle": "single",
    // ESLint Settings
    "eslint.enable": true,
    "eslint.codeAction.showDocumentation": {
        "enable": true
    },
    "eslint.codeAction.disableRuleComment": {
        "enable": true,
        "location": "sameLine"
    },
    // npm Intellisense Settings
    "npm-intellisense.importES6": true,
    "npm-intellisense.importQuotes": "'",
    "npm-intellisense.importLinebreak": ";\r\n",
    "npm-intellisense.importDeclarationType": "const",
    // Visual Studio IntelliCode Settings
    "editor.suggestSelection": "first",
    "vsintellicode.modify.editor.suggestSelection": "automaticallyOverrodeDefaultValue",
    // Debugger for Firefox Settings
    "firefox.keepProfileChanges": false,
    "firefox.port": 6000,
    "firefox.reloadOnChange": true
}
```

If you are using the [VSC-Essentials] extension pack additionally, you can see the complete settings file [here] ([settings.json])

Consider that if you would rather use ESLint as the default code formatter, rather than the one that comes with VSCode, you will need to change this setting in your settings.json file:

``` json
"[javascript]": {
    "editor.defaultFormatter": "vscode.typescript-language-features"
},
```

For this other one:

``` json
"[javascript]": {
    "editor.defaultFormatter": "dbaeumer.vscode-eslint"
},
```

As far as using the [Debugger for Firefox] entails, I highly recommend [reading its documentation] in order to get a full debugging experience.

## Note

The [VSC-Essentials] project was used as a template for this one. And it is highly recommended that it be installed in conjunction with this pack.

## Included

This extension pack includes the following extensions:

| Extension                      | Stats                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| ------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| VSC-Essentials-Core            | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/Gydunhn.vsc-essentials-core?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.vsc-essentials-core) [![Installs](https://flat.badgen.net/vs-marketplace/i/Gydunhn.vsc-essentials-core?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.vsc-essentials-core) [![Rating](https://flat.badgen.net/vs-marketplace/rating/Gydunhn.vsc-essentials-core?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.vsc-essentials-core)                                                                   |
| ESLint                         | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/dbaeumer.vscode-eslint?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) [![Installs](https://flat.badgen.net/vs-marketplace/i/dbaeumer.vscode-eslint?color=blue)](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) [![Rating](https://flat.badgen.net/vs-marketplace/rating/dbaeumer.vscode-eslint?color=blue)](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint)                                                                                                 |
| npm Intellisense               | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/christian-kohler.npm-intellisense?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=christian-kohler.npm-intellisense) [![Installs](https://flat.badgen.net/vs-marketplace/i/christian-kohler.npm-intellisense?color=blue)](https://marketplace.visualstudio.com/items?itemName=christian-kohler.npm-intellisense) [![Rating](https://flat.badgen.net/vs-marketplace/rating/christian-kohler.npm-intellisense?color=blue)](https://marketplace.visualstudio.com/items?itemName=christian-kohler.npm-intellisense)                               |
| Visual Studio IntelliCode      | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/VisualStudioExptTeam.vscodeintellicode?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=VisualStudioExptTeam.vscodeintellicode) [![Installs](https://flat.badgen.net/vs-marketplace/i/VisualStudioExptTeam.vscodeintellicode?color=blue)](https://marketplace.visualstudio.com/items?itemName=VisualStudioExptTeam.vscodeintellicode) [![Rating](https://flat.badgen.net/vs-marketplace/rating/VisualStudioExptTeam.vscodeintellicode?color=blue)](https://marketplace.visualstudio.com/items?itemName=VisualStudioExptTeam.vscodeintellicode) |
| JavaScript (ES6) code snippets | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/xabikos.JavaScriptSnippets?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=xabikos.JavaScriptSnippets) [![Installs](https://flat.badgen.net/vs-marketplace/i/xabikos.JavaScriptSnippets?color=blue)](https://marketplace.visualstudio.com/items?itemName=xabikos.JavaScriptSnippets) [![Rating](https://flat.badgen.net/vs-marketplace/rating/xabikos.JavaScriptSnippets?color=blue)](https://marketplace.visualstudio.com/items?itemName=xabikos.JavaScriptSnippets)                                                                         |
| Backticks                      | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/fractalbrew.backticks?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=fractalbrew.backticks) [![Installs](https://flat.badgen.net/vs-marketplace/i/fractalbrew.backticks?color=blue)](https://marketplace.visualstudio.com/items?itemName=fractalbrew.backticks) [![Rating](https://flat.badgen.net/vs-marketplace/rating/fractalbrew.backticks?color=blue)](https://marketplace.visualstudio.com/items?itemName=fractalbrew.backticks)                                                                                                       |
| Debugger for Firefox           | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/firefox-devtools.vscode-firefox-debug?icon=visualstudio&color=blue)](https://marketplace.visualstudio.com/items?itemName=firefox-devtools.vscode-firefox-debug) [![Installs](https://flat.badgen.net/vs-marketplace/i/firefox-devtools.vscode-firefox-debug?color=blue)](https://marketplace.visualstudio.com/items?itemName=firefox-devtools.vscode-firefox-debug) [![Rating](https://flat.badgen.net/vs-marketplace/rating/firefox-devtools.vscode-firefox-debug?color=blue)](https://marketplace.visualstudio.com/items?itemName=firefox-devtools.vscode-firefox-debug)       |

[VSC-Essentials]: https://github.com/Gydunhn/VSC-Essentials
[here]: /.vscode/settings.json
[settings.json]: /.vscode/settings.json
[Debugger for Firefox]: https://marketplace.visualstudio.com/items?itemName=firefox-devtools.vscode-firefox-debug
[reading its documentation]: https://github.com/firefox-devtools/vscode-firefox-debug
