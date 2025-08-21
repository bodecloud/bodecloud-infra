"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ListDTOsProvider = void 0;
const vscode_1 = require("vscode");
const configs_1 = require("../configs");
const controllers_1 = require("../controllers");
const models_1 = require("../models");
/**
 * The ListDTOsProvider class
 *
 * @class
 * @classdesc The class that represents the list of files provider.
 * @export
 * @public
 * @implements {TreeDataProvider<NodeModel>}
 * @property {EventEmitter<NodeModel | undefined | null | void>} _onDidChangeTreeData - The onDidChangeTreeData event emitter
 * @property {Event<NodeModel | undefined | null | void>} onDidChangeTreeData - The onDidChangeTreeData event
 * @property {ListFilesController} controller - The list of files controller
 * @example
 * const provider = new ListDTOsProvider();
 *
 * @see https://code.visualstudio.com/api/references/vscode-api#TreeDataProvider
 */
class ListDTOsProvider {
    // -----------------------------------------------------------------
    // Constructor
    // -----------------------------------------------------------------
    /**
     * Constructor for the ListDTOsProvider class
     *
     * @constructor
     * @public
     * @memberof ListDTOsProvider
     */
    constructor(controller) {
        this.controller = controller;
        this._onDidChangeTreeData = new vscode_1.EventEmitter();
        this.onDidChangeTreeData = this._onDidChangeTreeData.event;
    }
    // -----------------------------------------------------------------
    // Methods
    // -----------------------------------------------------------------
    // Public methods
    /**
     * Returns the tree item for the supplied element.
     *
     * @function getTreeItem
     * @param {NodeModel} element - The element
     * @public
     * @memberof ListDTOsProvider
     * @example
     * const treeItem = provider.getTreeItem(element);
     *
     * @returns {TreeItem | Thenable<TreeItem>} - The tree item
     *
     * @see https://code.visualstudio.com/api/references/vscode-api#TreeDataProvider
     */
    getTreeItem(element) {
        return element;
    }
    /**
     * Returns the children for the supplied element.
     *
     * @function getChildren
     * @param {NodeModel} [element] - The element
     * @public
     * @memberof ListDTOsProvider
     * @example
     * const children = provider.getChildren(element);
     *
     * @returns {ProviderResult<NodeModel[]>} - The children
     *
     * @see https://code.visualstudio.com/api/references/vscode-api#TreeDataProvider
     */
    getChildren(element) {
        if (element) {
            return element.children;
        }
        return this.getListDtos();
    }
    /**
     * Refreshes the tree data.
     *
     * @function refresh
     * @public
     * @memberof FeedbackProvider
     * @example
     * provider.refresh();
     *
     * @returns {void} - No return value
     */
    refresh() {
        this._onDidChangeTreeData.fire();
    }
    // Private methods
    /**
     * Returns the list of files.
     *
     * @function getListDtos
     * @private
     * @memberof ListDTOsProvider
     * @example
     * const files = provider.getListDtos();
     *
     * @returns {Promise<NodeModel[] | undefined>} - The list of files
     */
    async getListDtos() {
        const files = await controllers_1.ListFilesController.getFiles();
        if (!files) {
            return;
        }
        // List of DTOs
        const nodes = files.filter((file) => file.label.toString().includes('dto.ts'));
        const importRegex = this.controller.getAnnotationsRegex();
        if (!importRegex) {
            return;
        }
        for (const file of nodes) {
            const document = await vscode_1.workspace.openTextDocument(file.resourceUri?.path ?? '');
            const children = Array.from({ length: document.lineCount }, (_, index) => {
                const line = document.lineAt(index);
                let node;
                if (importRegex.test(line.text)) {
                    node = new models_1.NodeModel(line.text.trim(), new vscode_1.ThemeIcon('symbol-method'), {
                        command: `${configs_1.EXTENSION_ID}.list.gotoLine`,
                        title: line.text,
                        arguments: [file.resourceUri, index],
                    });
                }
                return node;
            });
            file.setChildren(children.filter((child) => child !== undefined));
        }
        return files.filter((file) => file.children && file.children.length !== 0);
    }
}
exports.ListDTOsProvider = ListDTOsProvider;
//# sourceMappingURL=list-dtos.providers.js.map