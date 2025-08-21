# TypeScript Essentials - Extension Pack for Visual Studio Code

[![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/Gydunhn.typescript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.typescript-essentials) [![Installs](https://flat.badgen.net/vs-marketplace/i/Gydunhn.typescript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.typescript-essentials) [![Downloads](https://flat.badgen.net/vs-marketplace/d/Gydunhn.typescript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.typescript-essentials) [![Rating](https://flat.badgen.net/vs-marketplace/rating/Gydunhn.typescript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.typescript-essentials)

This extension pack for Visual Studio Code adds extensions that are useful for TypeScript projects. I reserve the right to update the extensions pack contents up to my own discretion. This extension is for my personal use, I think it's great if it works for other people too.

## Reasons

The [TypeScript Essentials extension pack] was made to automate and standardize the installation phase of the essential TypeScript extensions for Visual Studio Code every time a new member joins the team, or one of them restores a laptop, or exchanges it for a new one.

See the [CHANGELOG](https://github.com/Gydunhn/Typescript-Essentials/blob/HEAD/CHANGELOG.md) for the latest changes

## **settings.json**

It is strongly recommended that these settings be used in your workspace. You must copy and paste them, and if you need to adjust something you will already know where to do it.

``` json
{
    /**
     * JavaScript Essentials Config
     */
	"[javascript]": {
		"editor.defaultFormatter": "vscode.typescript-language-features"
	},
	"javascript.format.enable": true,
	"javascript.format.semicolons": "insert",
	"javascript.preferences.quoteStyle": "single",
	"eslint.enable": true,
	"eslint.codeAction.showDocumentation": {
		"enable": true
	},
	"eslint.codeAction.disableRuleComment": {
		"enable": true,
		"location": "sameLine"
	},
	"npm-intellisense.importES6": true,
	"npm-intellisense.importQuotes": "'",
	"npm-intellisense.importLinebreak": ";\r\n",
	"npm-intellisense.importDeclarationType": "const",
	/**
     * The following line is for the specific configuration of the 
     * Path-Intellisense extension over Javascript
     */
	"javascript.suggest.paths": false,
	"path-intellisense.showHiddenFiles": true,
	"formattingToggle.affects": [
		"editor.formatOnPaste",
		"editor.formatOnType"
	],
	/**
     * TypeScript Essentials Config
     */
	"[typescript]": {
		"editor.defaultFormatter": "vscode.typescript-language-features"
	},
	"typescript.format.enable": true,
	"typescript.format.semicolons": "insert",
	"typescript.preferences.quoteStyle": "single",
	/**
     * The following line is for the specific configuration of the 
     * Path-Intellisense extension over Typescript
     */
	"typescript.suggest.paths": false,
}
```

If you are using the [VSC-Essentials] extension pack additionally, you can see the complete settings file [here] ([settings.json])

Consider that if you would rather use ESLint as the default code formatter, rather than the one that comes with VSCode, you will need to change this settings in your settings.json file:

``` json
"[javascript]": {
    "editor.defaultFormatter": "vscode.typescript-language-features"
},
"[typescript]": {
    "editor.defaultFormatter": "vscode.typescript-language-features"
},
```

For this others:

``` json
"[javascript]": {
    "editor.defaultFormatter": "dbaeumer.vscode-eslint"
},
"[typescript]": {
    "editor.defaultFormatter": "dbaeumer.vscode-eslint"
},
```

As far as using the [Debugger for Firefox] entails, I highly recommend [reading its documentation] in order to get a full debugging experience.

## TypeScript TSLint Language Service Plugin

the TSLint extension used to be part of this extension pack until [it was deprecated]. If you search the marketplace, it is [still there] if you need it, But in most cases what is recommended is to migrate to ESlint, which supports the same functionality. There is a diverse range of documentation of interest about the process, so I will leave below what I think is the most convenient:

* [tslint-to-eslint-config], Converts your TSLint configuration to the closest reasonable ESLint equivalent.
* [Migrate from TSLint to ESLint], TSLint has been the recommended linter in the past but now TSLint is deprecated and ESLint is taking over its duties. This article will help you migrate from TSLint to ESLint.
* [typescript-eslint], The tooling that enables ESLint and Prettier to support TypeScript.
* [TSLint GitHub], the original TSLint repository, from [Palantir Technologies].

## Note

The [VSC-Essentials] project was used as a template for this one.

## Included

This extension pack includes the following extensions:

| Extension                | Stats                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Javascript-Essentials    | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/Gydunhn.javascript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.javascript-essentials) [![Installs](https://flat.badgen.net/vs-marketplace/i/Gydunhn.javascript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.javascript-essentials) [![Rating](https://flat.badgen.net/vs-marketplace/rating/Gydunhn.javascript-essentials?color=blue)](https://marketplace.visualstudio.com/items?itemName=Gydunhn.javascript-essentials)                                           |
| TypeScript Importer      | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/pmneo.tsimporter?color=blue)](https://marketplace.visualstudio.com/items?itemName=pmneo.tsimporter) [![Installs](https://flat.badgen.net/vs-marketplace/i/pmneo.tsimporter?color=blue)](https://marketplace.visualstudio.com/items?itemName=pmneo.tsimporter) [![Rating](https://flat.badgen.net/vs-marketplace/rating/pmneo.tsimporter?color=blue)](https://marketplace.visualstudio.com/items?itemName=pmneo.tsimporter)                                                                                                                         |
| Total TypeScript         | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/mattpocock.ts-error-translator?color=blue)](https://marketplace.visualstudio.com/items?itemName=mattpocock.ts-error-translator) [![Installs](https://flat.badgen.net/vs-marketplace/i/mattpocock.ts-error-translator?color=blue)](https://marketplace.visualstudio.com/items?itemName=mattpocock.ts-error-translator) [![Rating](https://flat.badgen.net/vs-marketplace/rating/mattpocock.ts-error-translator?color=blue)](https://marketplace.visualstudio.com/items?itemName=mattpocock.ts-error-translator)                                     |
| Pretty TypeScript Errors | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/yoavbls.pretty-ts-errors?color=blue)](https://marketplace.visualstudio.com/items?itemName=yoavbls.pretty-ts-errors) [![Installs](https://flat.badgen.net/vs-marketplace/i/yoavbls.pretty-ts-errors?color=blue)](https://marketplace.visualstudio.com/items?itemName=yoavbls.pretty-ts-errors) [![Rating](https://flat.badgen.net/vs-marketplace/rating/yoavbls.pretty-ts-errors?color=blue)](https://marketplace.visualstudio.com/items?itemName=yoavbls.pretty-ts-errors)                                                                         |
| TypeScript Toolbox       | [![Badge for version for Visual Studio Code extension](https://flat.badgen.net/vs-marketplace/v/DSKWRK.vscode-generate-getter-setter?color=blue)](https://marketplace.visualstudio.com/items?itemName=DSKWRK.vscode-generate-getter-setter) [![Installs](https://flat.badgen.net/vs-marketplace/i/DSKWRK.vscode-generate-getter-setter?color=blue)](https://marketplace.visualstudio.com/items?itemName=DSKWRK.vscode-generate-getter-setter) [![Rating](https://flat.badgen.net/vs-marketplace/rating/DSKWRK.vscode-generate-getter-setter?color=blue)](https://marketplace.visualstudio.com/items?itemName=DSKWRK.vscode-generate-getter-setter) |

[VSC-Essentials]: https://github.com/Gydunhn/VSC-Essentials
[TypeScript Essentials extension pack]: https://marketplace.visualstudio.com/items?itemName=Gydunhn.typescript-essentials
[it was deprecated]: https://blog.palantir.com/tslint-in-2019-1a144c2317a9
[still there]: https://marketplace.visualstudio.com/items?itemName=ms-vscode.vscode-typescript-tslint-plugin
[tslint-to-eslint-config]: https://github.com/typescript-eslint/tslint-to-eslint-config
[Migrate from TSLint to ESLint]: https://code.visualstudio.com/api/advanced-topics/tslint-eslint-migration
[typescript-eslint]: https://typescript-eslint.io/
[TSLint GitHub]: https://github.com/palantir/tslint
[Palantir Technologies]: https://github.com/palantir
[here]: /.vscode/settings.json
[settings.json]: /.vscode/settings.json
[Debugger for Firefox]: https://marketplace.visualstudio.com/items?itemName=firefox-devtools.vscode-firefox-debug
[reading its documentation]: https://github.com/firefox-devtools/vscode-firefox-debug
