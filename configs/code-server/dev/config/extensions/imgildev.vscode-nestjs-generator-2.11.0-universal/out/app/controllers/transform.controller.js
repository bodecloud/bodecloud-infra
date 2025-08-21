"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TransformController = void 0;
const json_to_ts_1 = require("json-to-ts");
const json5 = require("json5");
const vscode_1 = require("vscode");
// Import the helper functions
const helpers_1 = require("../helpers");
/**
 * The TransformController class.
 *
 * @class
 * @classdesc The class that represents the example controller.
 * @export
 * @public
 * @example
 * const controller = new TransformController();
 */
class TransformController {
    // -----------------------------------------------------------------
    // Methods
    // -----------------------------------------------------------------
    // Public methods
    /**
     * The json2ts method.
     *
     * @function json2ts
     * @public
     * @async
     * @memberof TransformController
     * @example
     * await controller.json2ts();
     *
     * @returns {Promise<TextEditor | void>} The result
     */
    async json2ts() {
        let editor;
        if (vscode_1.workspace.workspaceFolders) {
            editor = vscode_1.window.activeTextEditor;
        }
        else {
            const message = vscode_1.l10n.t('No text editor is active!');
            (0, helpers_1.showError)(message);
            return;
        }
        const selection = editor?.selection;
        if (selection && !selection.isEmpty) {
            const selectionRange = new vscode_1.Range(selection.start.line, selection.start.character, selection.end.line, selection.end.character);
            let text = editor?.document.getText(selectionRange) || '';
            const languageId = editor?.document.languageId || '';
            if ([
                'javascript',
                'javascriptreact',
                'typescript',
                'typescriptreact',
            ].includes(languageId)) {
                text = text
                    .replace(/'([^']+)'/g, '"$1"')
                    .replace(/(['"])?([a-zA-Z0-9_]+)(['"])?:/g, '"$2":')
                    .replace(/,*\s*\n*\];*/g, ']')
                    .replace(/{\s*\n*/g, '{')
                    .replace(/,*\s*\n*};*/g, '}');
            }
            const jsonSchema = this.tryParseJSONObject(text);
            if (!jsonSchema) {
                const message = vscode_1.l10n.t('Invalid JSON Schema!');
                (0, helpers_1.showError)(message);
                return;
            }
            const tsSchema = (0, json_to_ts_1.default)(jsonSchema)
                .map((itf) => `export ${itf}\n`)
                .join('\n');
            const document = await vscode_1.workspace.openTextDocument({
                language: 'typescript',
                content: tsSchema,
            });
            return await vscode_1.window.showTextDocument(document);
        }
        const message = vscode_1.l10n.t('No text is selected!');
        (0, helpers_1.showError)(message);
        return;
    }
    // Private methods
    /**
     * The tryParseJSONObject method.
     *
     * @private
     * @memberof TransformController
     * @param {string} str - The string to parse
     * @returns {boolean | object} The result
     * @example
     * const object = controller.tryParseJSONObject(str);
     *
     * @returns {boolean | object} The result
     */
    tryParseJSONObject(str) {
        try {
            const object = json5.parse(str);
            if (object && typeof object === 'object') {
                return object;
            }
        }
        catch (e) {
            return false;
        }
        return false;
    }
}
exports.TransformController = TransformController;
//# sourceMappingURL=transform.controller.js.map