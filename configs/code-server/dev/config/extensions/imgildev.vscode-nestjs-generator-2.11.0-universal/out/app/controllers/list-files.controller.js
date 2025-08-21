"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ListFilesController = void 0;
const vscode_1 = require("vscode");
const configs_1 = require("../configs");
const helpers_1 = require("../helpers");
const models_1 = require("../models");
/**
 * The ListFilesController class.
 *
 * @class
 * @classdesc The class that represents the list files controller.
 * @export
 * @public
 * @example
 * const controller = new ListFilesController();
 */
class ListFilesController {
    // -----------------------------------------------------------------
    // Constructor
    // -----------------------------------------------------------------
    /**
     * Constructor for the ListFilesController class
     *
     * @constructor
     * @param {Config} config - The configuration object
     * @public
     * @memberof ListFilesController
     */
    constructor(config) {
        ListFilesController.config = config;
    }
    // -----------------------------------------------------------------
    // Methods
    // -----------------------------------------------------------------
    // Public methods
    /**
     * The getFiles method.
     *
     * @function getFiles
     * @param {number} maxResults - The maximum number of results
     * @public
     * @async
     * @memberof ListFilesController
     * @example
     * controller.getFiles();
     *
     * @returns {Promise<NodeModel[] | void>} - The list of files
     */
    static async getFiles(maxResults = Number.MAX_SAFE_INTEGER) {
        // Get the files in the folder
        const files = await (0, helpers_1.directoryMap)('/', {
            extensions: this.config.include,
            ignore: this.config.exclude,
            maxResults,
        });
        if (files.length !== 0) {
            let nodes = [];
            files.sort((a, b) => a.path.localeCompare(b.path));
            for (const file of files) {
                const document = await vscode_1.workspace.openTextDocument(file);
                const path = await (0, helpers_1.getRelativePath)(document.fileName);
                let filename = path.split('/').pop();
                if (filename && this.config.showPath) {
                    const folder = path.split('/').slice(0, -1).join('/');
                    filename += folder ? ` (${folder})` : ' (root)';
                }
                nodes.push(new models_1.NodeModel(filename ?? 'Untitled', new vscode_1.ThemeIcon('file'), {
                    command: `${configs_1.EXTENSION_ID}.list.openFile`,
                    title: 'Open File',
                    arguments: [document.uri],
                }, document.uri, document.fileName));
            }
            return nodes;
        }
        return;
    }
    /**
     * The openFile method.
     *
     * @function openFile
     * @param {Uri} uri - The file URI
     * @public
     * @memberof ListFilesController
     * @example
     * controller.openFile('file:///path/to/file');
     *
     * @returns {Promise<void>} - The promise
     */
    openFile(uri) {
        vscode_1.workspace.openTextDocument(uri).then((filename) => {
            vscode_1.window.showTextDocument(filename);
        });
    }
    /**
     * The gotoLine method.
     *
     * @function gotoLine
     * @param {Uri} uri - The file URI
     * @param {number} line - The line number
     * @public
     * @memberof ListFilesController
     * @example
     * controller.gotoLine('file:///path/to/file', 1);
     *
     * @returns {void} - The promise
     */
    gotoLine(uri, line) {
        vscode_1.workspace.openTextDocument(uri).then((document) => {
            vscode_1.window.showTextDocument(document).then((editor) => {
                const position = new vscode_1.Position(line, 0);
                editor.revealRange(new vscode_1.Range(position, position), vscode_1.TextEditorRevealType.InCenterIfOutsideViewport);
                editor.selection = new vscode_1.Selection(position, position);
            });
        });
    }
}
exports.ListFilesController = ListFilesController;
//# sourceMappingURL=list-files.controller.js.map