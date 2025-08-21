"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.showWarning = exports.showError = exports.showMessage = exports.pickItem = exports.getName = exports.getPath = void 0;
const vscode_1 = require("vscode");
/**
 * Displays a message box with the provided message
 *
 * @param {string} prompt - The prompt to display
 * @param {string} placeHolder - The input placeholder
 * @param {string} currentPath - The current path
 * @param {string} validate - The validation function
 * @example
 * const path = await getPath('Enter a path', 'src/app', 'src/app', (path) => {
 *   if (path.length === 0) {
 *     return 'Path cannot be empty';
 * });
 *
 * @returns {Promise<string | undefined>} - The selected path
 */
const getPath = async (prompt, placeHolder, currentPath, validate) => {
    return await vscode_1.window.showInputBox({
        prompt,
        placeHolder,
        value: currentPath,
        validateInput: validate,
    });
};
exports.getPath = getPath;
/**
 * Displays a message box with the provided message
 *
 * @param {string} prompt - The prompt to display
 * @param {string} placeHolder - The input placeholder
 * @param {string} validate - The validation function
 * @example
 * const name = await getName('Enter a name', 'foo', (name) => {
 *   if (name.length === 0) {
 *     return 'Name cannot be empty';
 * });
 *
 * @returns {Promise<string | undefined>} - The selected name
 */
const getName = async (prompt, placeHolder, validate) => {
    return await vscode_1.window.showInputBox({
        prompt,
        placeHolder,
        validateInput: validate,
    });
};
exports.getName = getName;
/**
 * Displays a message box with the provided message
 *
 * @param {string[]} items - The list of items to select from
 * @param {string} placeHolder - The input placeholder
 * @example
 * const item = await pickItem(['foo', 'bar'], 'Select an item');
 *
 * @returns {Promise<string | undefined>} - The selected item
 */
const pickItem = async (items, placeHolder) => {
    return await vscode_1.window.showQuickPick(items, {
        placeHolder,
    });
};
exports.pickItem = pickItem;
/**
 * Displays a message box with the provided message
 *
 * @param {string} message - The message to display
 * @example
 * await showMessage('Hello, world!');
 *
 * @returns {Promise<void>} - No return value
 */
const showMessage = async (message) => {
    vscode_1.window.showInformationMessage(message);
};
exports.showMessage = showMessage;
/**
 * Displays a message box with the provided message
 *
 * @param {string} message - The message to display
 * @example
 * await showError('An error occurred');
 *
 * @returns {Promise<void>} - No return value
 */
const showError = async (message) => {
    vscode_1.window.showErrorMessage(message);
};
exports.showError = showError;
/**
 * Displays a message box with the provided message
 *
 * @param {string} message - The message to display
 * @example
 * await showWarning('This is a warning');
 *
 * @returns {Promise<void>} - No return value
 */
const showWarning = async (message) => {
    vscode_1.window.showWarningMessage(message);
};
exports.showWarning = showWarning;
//# sourceMappingURL=dialog.helper.js.map