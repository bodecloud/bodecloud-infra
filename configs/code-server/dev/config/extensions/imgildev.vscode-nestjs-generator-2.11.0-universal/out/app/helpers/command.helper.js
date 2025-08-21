"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.runCommand = void 0;
const vscode_1 = require("vscode");
/**
 * Runs a command in the terminal
 *
 * @param {string} title - Title of the terminal
 * @param {string} command - Command to run
 * @example
 * runCommand('echo "Hello, World!"');
 *
 * @returns {Promise<void>} - No return value
 */
const runCommand = async (name, command, path) => {
    let cwd;
    if (path && vscode_1.workspace.getWorkspaceFolder(vscode_1.Uri.file(path))) {
        cwd = vscode_1.Uri.file(path);
    }
    const terminal = vscode_1.window.createTerminal({
        name,
        cwd,
    });
    terminal.show();
    terminal.sendText(command);
};
exports.runCommand = runCommand;
//# sourceMappingURL=command.helper.js.map