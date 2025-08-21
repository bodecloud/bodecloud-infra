"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Config = void 0;
const vscode_1 = require("vscode");
const constants_1 = require("./constants");
/**
 * The Config class.
 *
 * @class
 * @classdesc The class that represents the configuration of the extension.
 * @export
 * @public
 * @property {WorkspaceConfiguration} config - The workspace configuration
 * @property {string[]} include - The files to include
 * @property {string[]} exclude - The files to exclude
 * @property {string[]} watch - The files to watch
 * @property {boolean} showPath - Whether to show the path or not
 * @property {object[]} customCommands - The custom commands
 * @property {object[]} customTemplates - The custom templates
 * @property {object} activateItem - Whether to show the menu or not
 * @property {boolean} autoImport - The auto import setting
 * @property {boolean} skipFolderConfirmation - Whether to skip the folder confirmation or not
 * @property {string} orm - The orm
 * @example
 * const config = new Config(workspace.getConfiguration());
 * console.log(config.include);
 * console.log(config.exclude);
 * console.log(config.watch);
 * console.log(config.autoImport);
 */
class Config {
    // -----------------------------------------------------------------
    // Constructor
    // -----------------------------------------------------------------
    /**
     * Constructor for the Config class.
     *
     * @constructor
     * @param {WorkspaceConfiguration} config - The workspace configuration
     * @public
     * @memberof Config
     */
    constructor(config) {
        this.config = config;
        this.enable = config.get('enable', true);
        this.include = config.get('files.include', constants_1.INCLUDE);
        this.exclude = config.get('files.exclude', constants_1.EXCLUDE);
        this.watch = config.get('files.watch', constants_1.WATCH);
        this.showPath = config.get('files.showPath', constants_1.SHOW_PATH);
        this.cwd = config.get('terminal.cwd', vscode_1.workspace.workspaceFolders?.[0].uri.fsPath);
        this.customCommands = config.get('submenu.customCommands', constants_1.CUSTOM_COMMANDS);
        this.templates = config.get('submenu.templates', constants_1.CUSTOM_TEMPLATES);
        this.activateItem = config.get('submenu.activateItem', constants_1.ACTIVATE_MENU);
        this.autoImport = config.get('files.autoImport', constants_1.AUTO_IMPORT);
        this.skipFolderConfirmation = config.get('files.skipFolderConfirmation', constants_1.SKIP_FOLDER_CONFIRMATION);
        this.orm = config.get('files.orm', constants_1.ORM);
    }
    // -----------------------------------------------------------------
    // Methods
    // -----------------------------------------------------------------
    // Public methods
    /**
     * The update method.
     *
     * @function update
     * @param {WorkspaceConfiguration} config - The workspace configuration
     * @public
     * @memberof Config
     * @example
     * const config = new Config(workspace.getConfiguration());
     * config.update(workspace.getConfiguration());
     */
    update(config) {
        this.enable = config.get('enable', true);
        this.include = config.get('files.include', constants_1.INCLUDE);
        this.exclude = config.get('files.exclude', constants_1.EXCLUDE);
        this.watch = config.get('files.watch', constants_1.WATCH);
        this.showPath = config.get('files.showPath', constants_1.SHOW_PATH);
        this.cwd = config.get('terminal.cwd', vscode_1.workspace.workspaceFolders?.[0].uri.fsPath);
        this.customCommands = config.get('submenu.customCommands', constants_1.CUSTOM_COMMANDS);
        this.templates = config.get('submenu.templates', constants_1.CUSTOM_TEMPLATES);
        this.activateItem = config.get('submenu.activateItem', constants_1.ACTIVATE_MENU);
        this.autoImport = config.get('files.autoImport', constants_1.AUTO_IMPORT);
        this.skipFolderConfirmation = config.get('files.skipFolderConfirmation', constants_1.SKIP_FOLDER_CONFIRMATION);
        this.orm = config.get('files.orm', constants_1.ORM);
    }
}
exports.Config = Config;
//# sourceMappingURL=config.js.map