"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TerminalController = void 0;
const vscode_1 = require("vscode");
const helpers_1 = require("../helpers");
/**
 * The TerminalController class.
 *
 * @class
 * @classdesc The class that represents the example controller.
 * @export
 * @public
 * @property {Config} config - The configuration
 * @example
 * const controller = new TerminalController(config);
 */
class TerminalController {
    // -----------------------------------------------------------------
    // Constructor
    // -----------------------------------------------------------------
    /**
     * Constructor for the TerminalController class.
     *
     * @constructor
     * @param {Config} config - The configuration
     * @public
     * @memberof TerminalController
     */
    constructor(config) {
        this.config = config;
    }
    // -----------------------------------------------------------------
    // Methods
    // -----------------------------------------------------------------
    // Public methods
    /**
     * The generateController method.
     *
     * @function generateController
     * @param {Uri} [path] The path to the folder.
     * @public
     * @async
     * @memberof TerminalController
     * @example
     * await generateController();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    async generateController(path) {
        // Get the relative path
        let folderPath = path ? vscode_1.workspace.asRelativePath(path.path) : '';
        if (this.config.cwd) {
            const cwd = vscode_1.workspace.asRelativePath(vscode_1.Uri.file(this.config.cwd).path);
            folderPath = folderPath.replace(cwd, '');
            if (folderPath.startsWith('/')) {
                folderPath = folderPath.substring(1);
            }
        }
        // Get the path to the folder
        const folder = await (0, helpers_1.getPath)('Controller name', 'Controller name. E.g. modules/cats, modules/users, modules/projects...', `${folderPath}/`, (path) => {
            if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                return 'The folder name must be a valid name';
            }
            return;
        });
        if (!folder) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const items = [
            {
                label: '--dry-run',
                description: 'Report actions that would be taken without writing out results.',
            },
            {
                label: '--flat',
                description: 'Enforce flat structure of generated element.',
            },
            {
                label: '--skip-import',
                description: 'Skip importing (default: false)',
            },
            {
                label: '--no-spec',
                description: 'Disable spec files generation.',
            },
        ];
        const options = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the options for the controller generation (optional)'),
            canPickMany: true,
        });
        const filename = folder.replace('src/', '');
        const command = `nest g co ${filename}` +
            (options ? ' ' + options.map((item) => item.label).join(' ') : '');
        (0, helpers_1.runCommand)('generate controller', command, this.config.cwd);
    }
    /**
     * The generateGateway method.
     *
     * @function generateGateway
     * @param {Uri} [path] The path to the folder.
     * @public
     * @async
     * @memberof TerminalController
     * @example
     * await generateGateway();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    async generateGateway(path) {
        // Get the relative path
        let folderPath = path ? vscode_1.workspace.asRelativePath(path.path) : '';
        if (this.config.cwd) {
            const cwd = vscode_1.workspace.asRelativePath(vscode_1.Uri.file(this.config.cwd).path);
            folderPath = folderPath.replace(cwd, '');
            if (folderPath.startsWith('/')) {
                folderPath = folderPath.substring(1);
            }
        }
        // Get the path to the folder
        const folder = await (0, helpers_1.getPath)('Gateway name', 'Gateway name. E.g. modules/cats, modules/users, modules/projects...', `${folderPath}/`, (path) => {
            if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                return 'The folder name must be a valid name';
            }
            return;
        });
        if (!folder) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const items = [
            {
                label: '--dry-run',
                description: 'Report actions that would be taken without writing out results.',
            },
            {
                label: '--flat',
                description: 'Enforce flat structure of generated element.',
            },
            {
                label: '--skip-import',
                description: 'Skip importing (default: false)',
            },
            {
                label: '--no-spec',
                description: 'Disable spec files generation.',
            },
        ];
        const options = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the options for the gateway generation (optional)'),
            canPickMany: true,
        });
        const filename = folder.replace('src/', '');
        const command = `nest g ga ${filename}` +
            (options ? ' ' + options.map((item) => item.label).join(' ') : '');
        (0, helpers_1.runCommand)('generate gateway', command, this.config.cwd);
    }
    /**
     * The generateLibrary method.
     *
     * @function generateLibrary
     * @public
     * @async
     * @memberof TerminalController
     * @example
     * await generateLibrary();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    async generateLibrary() {
        const folder = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter library name'), 'Library name. E.g. cats, users, projects...', (path) => {
            if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                return 'The folder name must be a valid name';
            }
            return;
        });
        if (!folder) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const items = [
            {
                label: '--dry-run',
                description: 'Report actions that would be taken without writing out results.',
            },
            {
                label: '--flat',
                description: 'Enforce flat structure of generated element.',
            },
            {
                label: '--skip-import',
                description: 'Skip importing (default: false)',
            },
            {
                label: '--no-spec',
                description: 'Disable spec files generation.',
            },
        ];
        const options = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the options for the library generation (optional)'),
            canPickMany: true,
        });
        const filename = folder.replace('src/', '');
        const command = `nest g lib ${filename}` +
            (options ? ' ' + options.map((item) => item.label).join(' ') : '');
        (0, helpers_1.runCommand)('generate library', command, this.config.cwd);
    }
    /**
     * The generateModule method.
     *
     * @function generateModule
     * @param {Uri} [path] The path to the folder.
     * @public
     * @async
     * @memberof TerminalController
     * @example
     * await generateModule();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    async generateModule(path) {
        // Get the relative path
        let folderPath = path ? vscode_1.workspace.asRelativePath(path.path) : '';
        if (this.config.cwd) {
            const cwd = vscode_1.workspace.asRelativePath(vscode_1.Uri.file(this.config.cwd).path);
            folderPath = folderPath.replace(cwd, '');
            if (folderPath.startsWith('/')) {
                folderPath = folderPath.substring(1);
            }
        }
        // Get the path to the folder
        const folder = await (0, helpers_1.getPath)('Module name', 'Module name. E.g. modules/cats, modules/users, modules/projects...', `${folderPath}/`, (path) => {
            if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                return 'The folder name must be a valid name';
            }
            return;
        });
        if (!folder) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const items = [
            {
                label: '--dry-run',
                description: 'Report actions that would be taken without writing out results.',
            },
            {
                label: '--flat',
                description: 'Enforce flat structure of generated element.',
            },
            {
                label: '--skip-import',
                description: 'Skip importing (default: false)',
            },
            {
                label: '--no-spec',
                description: 'Disable spec files generation.',
            },
        ];
        const options = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the options for the module generation (optional)'),
            canPickMany: true,
        });
        const filename = folder.replace('src/', '');
        const command = `nest g mo ${filename}` +
            (options ? ' ' + options.map((item) => item.label).join(' ') : '');
        (0, helpers_1.runCommand)('generate module', command, this.config.cwd);
    }
    /**
     * The generateProvider method.
     *
     * @function generateProvider
     * @param {Uri} [path] The path to the folder.
     * @public
     * @async
     * @memberof TerminalController
     * @example
     * await generateProvider();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    async generateProvider(path) {
        // Get the relative path
        let folderPath = path ? vscode_1.workspace.asRelativePath(path.path) : '';
        if (this.config.cwd) {
            const cwd = vscode_1.workspace.asRelativePath(vscode_1.Uri.file(this.config.cwd).path);
            folderPath = folderPath.replace(cwd, '');
            if (folderPath.startsWith('/')) {
                folderPath = folderPath.substring(1);
            }
        }
        // Get the path to the folder
        const folder = await (0, helpers_1.getPath)('Provider name', 'Provider name. E.g. providers/cats, providers/users, providers/projects...', `${folderPath}/`, (path) => {
            if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                return 'The folder name must be a valid name';
            }
            return;
        });
        if (!folder) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const items = [
            {
                label: '--dry-run',
                description: 'Report actions that would be taken without writing out results.',
            },
            {
                label: '--flat',
                description: 'Enforce flat structure of generated element.',
            },
            {
                label: '--skip-import',
                description: 'Skip importing (default: false)',
            },
            {
                label: '--no-spec',
                description: 'Disable spec files generation.',
            },
        ];
        const options = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the options for the provider generation (optional)'),
            canPickMany: true,
        });
        const filename = folder.replace('src/', '');
        const command = `nest g pr ${filename}` +
            (options ? ' ' + options.map((item) => item.label).join(' ') : '');
        (0, helpers_1.runCommand)('generate provider', command, this.config.cwd);
    }
    /**
     * The generateResolver method.
     *
     * @function generateResolver
     * @param {Uri} [path] The path to the folder.
     * @public
     * @async
     * @memberof TerminalController
     * @example
     * await generateResolver();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    async generateResolver(path) {
        // Get the relative path
        let folderPath = path ? vscode_1.workspace.asRelativePath(path.path) : '';
        if (this.config.cwd) {
            const cwd = vscode_1.workspace.asRelativePath(vscode_1.Uri.file(this.config.cwd).path);
            folderPath = folderPath.replace(cwd, '');
            if (folderPath.startsWith('/')) {
                folderPath = folderPath.substring(1);
            }
        }
        // Get the path to the folder
        const folder = await (0, helpers_1.getPath)('Resolver name', 'Resolver name. E.g. resolvers/cats, resolvers/users, resolvers/projects...', `${folderPath}/`, (path) => {
            if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                return 'The folder name must be a valid name';
            }
            return;
        });
        if (!folder) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const items = [
            {
                label: '--dry-run',
                description: 'Report actions that would be taken without writing out results.',
            },
            {
                label: '--flat',
                description: 'Enforce flat structure of generated element.',
            },
            {
                label: '--skip-import',
                description: 'Skip importing (default: false)',
            },
            {
                label: '--no-spec',
                description: 'Disable spec files generation.',
            },
        ];
        const options = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the options for the resolver generation (optional)'),
            canPickMany: true,
        });
        const filename = folder.replace('src/', '');
        const command = `nest g r ${filename}` +
            (options ? ' ' + options.map((item) => item.label).join(' ') : '');
        (0, helpers_1.runCommand)('generate resolver', command, this.config.cwd);
    }
    /**
     * The generateResource method.
     *
     * @function generateResource
     * @param {Uri} [path] The path to the folder.
     * @public
     * @async
     * @memberof TerminalController
     * @example
     * await generateResource();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    async generateResource(path) {
        // Get the relative path
        let folderPath = path ? vscode_1.workspace.asRelativePath(path.path) : '';
        if (this.config.cwd) {
            const cwd = vscode_1.workspace.asRelativePath(vscode_1.Uri.file(this.config.cwd).path);
            folderPath = folderPath.replace(cwd, '');
            if (folderPath.startsWith('/')) {
                folderPath = folderPath.substring(1);
            }
        }
        // Get the path to the folder
        const folder = await (0, helpers_1.getPath)('Resource name', 'Resource name. E.g. modules/cats, modules/users, modules/projects...', `${folderPath}/`, (path) => {
            if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                return 'The folder name must be a valid name';
            }
            return;
        });
        if (!folder) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const items = [
            {
                label: '--dry-run',
                description: 'Report actions that would be taken without writing out results.',
            },
            {
                label: '--flat',
                description: 'Enforce flat structure of generated element.',
            },
            {
                label: '--skip-import',
                description: 'Skip importing (default: false)',
            },
            {
                label: '--no-spec',
                description: 'Disable spec files generation.',
            },
        ];
        const options = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the options for the resource generation (optional)'),
            canPickMany: true,
        });
        const filename = folder.replace('src/', '');
        const command = `nest g res ${filename}` +
            (options ? ' ' + options.map((item) => item.label).join(' ') : '');
        (0, helpers_1.runCommand)('generate resource', command, this.config.cwd);
    }
    /**
     * The start method.
     *
     * @function start
     * @public
     * @memberof TerminalController
     * @example
     * start();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    start() {
        (0, helpers_1.runCommand)('start', 'nest start', this.config.cwd);
    }
    /**
     * The startDev method.
     *
     * @function startDev
     * @public
     * @memberof TerminalController
     * @example
     * await startDev();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    startDev() {
        (0, helpers_1.runCommand)('dev', 'nest start --watch', this.config.cwd);
    }
    /**
     * The startDebug method.
     *
     * @function startDebug
     * @public
     * @memberof TerminalController
     * @example
     * await startDebug();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    startDebug() {
        (0, helpers_1.runCommand)('debug', 'nest start --debug --watch', this.config.cwd);
    }
    /**
     * The startProd method.
     *
     * @function startProd
     * @public
     * @memberof TerminalController
     * @example
     * await startProd();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    startProd() {
        (0, helpers_1.runCommand)('prod', 'node dist/main', this.config.cwd);
    }
    /**
     * The generateService method.
     *
     * @function generateService
     * @param {Uri} [path] The path to the folder.
     * @public
     * @async
     * @memberof TerminalController
     * @example
     * await generateService();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    async generateService(path) {
        // Get the relative path
        let folderPath = path ? vscode_1.workspace.asRelativePath(path.path) : '';
        if (this.config.cwd) {
            const cwd = vscode_1.workspace.asRelativePath(vscode_1.Uri.file(this.config.cwd).path);
            folderPath = folderPath.replace(cwd, '');
            if (folderPath.startsWith('/')) {
                folderPath = folderPath.substring(1);
            }
        }
        // Get the path to the folder
        const folder = await (0, helpers_1.getPath)('Service name', 'Service name. E.g. services/cats, services/users, services/projects...', `${folderPath}/`, (path) => {
            if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                return 'The folder name must be a valid name';
            }
            return;
        });
        if (!folder) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const items = [
            {
                label: '--dry-run',
                description: 'Report actions that would be taken without writing out results.',
            },
            {
                label: '--flat',
                description: 'Enforce flat structure of generated element.',
            },
            {
                label: '--skip-import',
                description: 'Skip importing (default: false)',
            },
            {
                label: '--no-spec',
                description: 'Disable spec files generation.',
            },
        ];
        const options = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the options for the service generation (optional)'),
            canPickMany: true,
        });
        const filename = folder.replace('src/', '');
        const command = `nest g s ${filename}` +
            (options ? ' ' + options.map((item) => item.label).join(' ') : '');
        (0, helpers_1.runCommand)('generate service', command, this.config.cwd);
    }
    /**
     * The generateSubApp method.
     *
     * @function generateSubApp
     * @public
     * @memberof TerminalController
     * @example
     * await generateSubApp();
     *
     * @returns {Promise<void>} The promise that resolves the method.
     */
    async generateSubApp() {
        const folder = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter sub-app name'), 'Sub-app name. E.g. cats, users, projects...', (path) => {
            if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                return 'The folder name must be a valid name';
            }
            return;
        });
        if (!folder) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const items = [
            {
                label: '--dry-run',
                description: 'Report actions that would be taken without writing out results.',
            },
            {
                label: '--flat',
                description: 'Enforce flat structure of generated element.',
            },
            {
                label: '--skip-import',
                description: 'Skip importing (default: false)',
            },
            {
                label: '--no-spec',
                description: 'Disable spec files generation.',
            },
        ];
        const options = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the options for the sub-app generation (optional)'),
            canPickMany: true,
        });
        const filename = folder.replace('src/', '');
        const command = `nest g app ${filename}` +
            (options ? ' ' + options.map((item) => item.label).join(' ') : '');
        (0, helpers_1.runCommand)('generate sub-app', command, this.config.cwd);
    }
    /**
     * Generates a custom element.
     *
     * @function generateCustomElement
     * @param {Uri} [path] - The path
     * @public
     * @async
     * @memberof TerminalController
     * @example
     * controller.generateCustomElement();
     *
     * @returns {Promise<void>} - No return value
     */
    async generateCustomElement(path) {
        // Get the relative path
        let folderPath = path ? vscode_1.workspace.asRelativePath(path.path) : '';
        if (this.config.cwd) {
            const cwd = vscode_1.workspace.asRelativePath(vscode_1.Uri.file(this.config.cwd).path);
            folderPath = folderPath.replace(cwd, '');
            if (folderPath.startsWith('/')) {
                folderPath = folderPath.substring(1);
            }
        }
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', `${folderPath}/`, (path) => {
                if (!/^(?!\/)[^\sÀ-ÿ]+?$/.test(path)) {
                    return 'The folder name must be a valid name';
                }
                return;
            });
            if (!folder) {
                const message = vscode_1.l10n.t('Operation cancelled!');
                (0, helpers_1.showError)(message);
                return;
            }
        }
        else {
            folder = folderPath;
        }
        if (this.config.customCommands.length === 0) {
            const message = vscode_1.l10n.t('The custom commands list is empty. Please add custom commands to the configuration');
            vscode_1.window.showErrorMessage(message);
            return;
        }
        const items = this.config.customCommands.map((item) => {
            return {
                label: item.name,
                description: item.command,
                detail: item.args,
            };
        });
        const option = await vscode_1.window.showQuickPick(items, {
            placeHolder: vscode_1.l10n.t('Select the template for the custom element generation'),
        });
        if (!option) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        folder = folder.replace('src/', '');
        const command = `${option.description} ${folder} ${option.detail}`;
        (0, helpers_1.runCommand)('generate custom element', command, this.config.cwd);
    }
}
exports.TerminalController = TerminalController;
//# sourceMappingURL=terminal.controller.js.map