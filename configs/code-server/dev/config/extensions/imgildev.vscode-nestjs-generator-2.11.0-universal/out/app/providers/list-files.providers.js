"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ListFilesProvider = void 0;
const vscode_1 = require("vscode");
const controllers_1 = require("../controllers");
const helpers_1 = require("../helpers");
const models_1 = require("../models");
/**
 * The ListFilesProvider class
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
 * const provider = new ListFilesProvider();
 *
 * @see https://code.visualstudio.com/api/references/vscode-api#TreeDataProvider
 */
class ListFilesProvider {
    // -----------------------------------------------------------------
    // Constructor
    // -----------------------------------------------------------------
    /**
     * Constructor for the ListFilesProvider class
     *
     * @constructor
     * @public
     * @memberof ListFilesProvider
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
     * @memberof ListFilesProvider
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
     * @memberof ListFilesProvider
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
        return this.getListFiles();
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
     * Gets the list of files.
     *
     * @function getListFiles
     * @private
     * @memberof ListFilesProvider
     * @example
     * const files = provider.getListFiles();
     *
     * @returns {Promise<NodeModel[] | undefined>} - The list of files
     */
    async getListFiles() {
        const files = await controllers_1.ListFilesController.getFiles();
        if (!files) {
            return;
        }
        const nodes = [];
        const fileTypes = controllers_1.ListFilesController.config.watch;
        for (const fileType of fileTypes) {
            const children = files.filter((file) => file.label.toString().includes(`${(0, helpers_1.singularize)(fileType)}.ts`));
            if (children.length !== 0) {
                const node = new models_1.NodeModel(`${(0, helpers_1.titleize)(fileType)}: ${children.length}`, new vscode_1.ThemeIcon('folder-opened'), undefined, undefined, fileType, children);
                nodes.push(node);
            }
        }
        if (nodes.length === 0) {
            return;
        }
        return nodes;
    }
}
exports.ListFilesProvider = ListFilesProvider;
//# sourceMappingURL=list-files.providers.js.map