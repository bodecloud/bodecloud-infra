"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ListModulesProvider = void 0;
const vscode_1 = require("vscode");
const configs_1 = require("../configs");
const controllers_1 = require("../controllers");
const models_1 = require("../models");
/**
 * The ListModulesProvider class
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
 * const provider = new ListModulesProvider();
 *
 * @see https://code.visualstudio.com/api/references/vscode-api#TreeDataProvider
 */
class ListModulesProvider {
    // -----------------------------------------------------------------
    // Constructor
    // -----------------------------------------------------------------
    /**
     * Constructor for the ListModulesProvider class
     *
     * @constructor
     * @public
     * @memberof ListModulesProvider
     */
    constructor() {
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
     * @memberof ListModulesProvider
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
     * @memberof ListModulesProvider
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
        return this.getListModules();
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
     * @function getListModules
     * @private
     * @memberof ListModulesProvider
     * @example
     * const files = provider.getListModules();
     *
     * @returns {Promise<NodeModel[] | undefined>} - The list of files
     */
    async getListModules() {
        const files = await controllers_1.ListFilesController.getFiles();
        if (!files) {
            return;
        }
        // List of Modules
        const nodes = files.filter((file) => file.label.toString().includes('module.ts'));
        for (const file of nodes) {
            const document = await vscode_1.workspace.openTextDocument(file.resourceUri?.path ?? '');
            const children = Array.from({ length: document.lineCount }, (_, index) => {
                const line = document.lineAt(index);
                let node;
                if (line.text.match(/(providers|controllers|imports|exports): \[/g)) {
                    node = new models_1.NodeModel(line.text.trim(), new vscode_1.ThemeIcon('symbol-module'), {
                        command: `${configs_1.EXTENSION_ID}.list.gotoLine`,
                        title: line.text,
                        arguments: [file.resourceUri, index],
                    });
                }
                return node;
            });
            file.setChildren(children.filter((child) => child !== undefined));
        }
        return nodes.filter((file) => file.children && file.children.length !== 0);
    }
}
exports.ListModulesProvider = ListModulesProvider;
//# sourceMappingURL=list-modules.providers.js.map