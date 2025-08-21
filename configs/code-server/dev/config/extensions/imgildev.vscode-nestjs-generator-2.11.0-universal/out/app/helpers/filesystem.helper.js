"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isDirectory = exports.exists = exports.getRealpath = exports.getRelativePath = exports.setRealpath = exports.sameFile = exports.octalPermissions = exports.symbolicPermissions = exports.getDirFileInfo = exports.getFileInfo = exports.getFilenames = exports.deleteFiles = exports.saveFile = exports.directoryMap = void 0;
const fs_1 = require("fs");
const path_1 = require("path");
const vscode_1 = require("vscode");
/**
 * Reads the contents of the file specified in the path.
 *
 * @param {string} path - Path to the source directory
 * @param {object} [options] - Options for the directoryMap function
 * @param {string[]} [options.extensions] - File extensions to include
 * @param {string[]} [options.ignore] - Directories to ignore
 * @param {number} [options.maxResults] - An upper-bound for the result
 * @example
 * const files = await directoryMap('src', {
 *   extensions: ['ts'],
 *   ignore: ['**​/node_modules/**'],
 *   maxResults: 100,
 * });
 *
 * @returns {Promise<Uri[]>} - Array of files
 */
const directoryMap = async (path, options) => {
    let includes = path === '/' ? '**/*' : `${path}/**/*`;
    let exclude = '';
    if (options && options.extensions && options.extensions.length) {
        includes += `.{${options.extensions.join(',')}}`;
    }
    if (options && options.ignore && options.ignore.length) {
        exclude = `{${options.ignore.join(',')}}`;
    }
    return vscode_1.workspace.findFiles(includes, exclude, options?.maxResults);
};
exports.directoryMap = directoryMap;
/**
 * Writes data to the file specified in the path. If the file does not exist then the function will create it.
 *
 * @param {string} path - Path to the file
 * @param {string} filename - Name of the file
 * @param {string} data - Data to write to the file
 * @example
 * await saveFile('src', 'file.ts', 'console.log("Hello World")');
 *
 * @returns {Promise<void>} - Confirmation of the write operation
 */
const saveFile = async (path, filename, data) => {
    let folder = '';
    if (vscode_1.workspace.workspaceFolders) {
        folder = vscode_1.workspace.workspaceFolders[0].uri.fsPath;
    }
    else {
        const message = vscode_1.l10n.t('The file has not been created!');
        vscode_1.window.showErrorMessage(message);
        return;
    }
    const file = (0, path_1.join)(folder, path, filename);
    if (!(0, fs_1.existsSync)((0, path_1.dirname)(file))) {
        (0, fs_1.mkdirSync)((0, path_1.dirname)(file), { recursive: true });
    }
    (0, fs_1.access)(file, (err) => {
        if (err) {
            (0, fs_1.open)(file, 'w+', (err, fd) => {
                if (err) {
                    throw err;
                }
                (0, fs_1.writeFile)(fd, data, 'utf8', (err) => {
                    if (err) {
                        throw err;
                    }
                    const openPath = vscode_1.Uri.file(file);
                    vscode_1.workspace.openTextDocument(openPath).then((filename) => {
                        vscode_1.window.showTextDocument(filename);
                    });
                });
            });
            const message = vscode_1.l10n.t('Successfully created the file!');
            vscode_1.window.showInformationMessage(message);
        }
        else {
            const message = vscode_1.l10n.t('The name already exists!');
            vscode_1.window.showWarningMessage(message);
        }
    });
};
exports.saveFile = saveFile;
/**
 * Deletes ALL files contained in the supplied path.
 *
 * @param {string} path - Path to the directory
 * @param {object} [options] - Options for the deleteFiles function
 * @param {boolean} [options.recursive] - Delete the content recursively if a folder is denoted.
 * @param {boolean} [options.useTrash] - Use the trash instead of permanently deleting the files.
 * @example
 * await deleteFiles('src');
 *
 * @returns {Promise<void>} - No return value
 */
const deleteFiles = async (path, options) => {
    const files = await vscode_1.workspace.findFiles(`${path}/**/*`);
    files.forEach((file) => {
        (0, fs_1.access)(file.path, (err) => {
            if (err) {
                throw err;
            }
            vscode_1.workspace.fs.delete(file, options);
        });
    });
};
exports.deleteFiles = deleteFiles;
/**
 * Returns an array of filenames in the supplied path.
 *
 * @param {string} path - Path to the directory
 * @param {object} [options] - Options for the directoryMap function
 * @param {string[]} [options.extensions] - File extensions to include
 * @param {string[]} [options.ignore] - Directories to ignore
 * @param {number} [options.maxResults] - An upper-bound for the result.
 * @example
 * const files = await getFilenames('src');
 *
 * @returns {Promise<string[]>} - Array of filenames
 */
const getFilenames = async (path, options) => {
    const files = await (0, exports.directoryMap)(path, options);
    return files.map((file) => file.path);
};
exports.getFilenames = getFilenames;
/**
 * Returns an object containing the file information for the supplied path.
 *
 * @param {string} path - Path to the file
 * @example
 * const fileInfo = await getFileInfo('src/file.ts');
 *
 * @returns {Promise<object>} - File information
 */
const getFileInfo = async (path) => {
    return await vscode_1.workspace.fs.stat(vscode_1.Uri.file(path));
};
exports.getFileInfo = getFileInfo;
/**
 * Returns an object containing the directory information for the supplied path.
 *
 * @param {string} path - Path to the directory
 * @example
 * const dirInfo = await getDirFileInfo('src');
 *
 * @returns {Promise<object>} - Directory information
 */
const getDirFileInfo = async (path) => {
    return await vscode_1.workspace.fs.stat(vscode_1.Uri.file(path));
};
exports.getDirFileInfo = getDirFileInfo;
/**
 * Returns the symbolic permissions for the supplied path.
 *
 * @param {string} path - Path to the file
 * @example
 * const permissions = await symbolicPermissions('src/file.ts');
 *
 * @returns {Promise<FilePermission | undefined>} - Symbolic permissions
 */
const symbolicPermissions = async (path) => {
    return await vscode_1.workspace.fs
        .stat(vscode_1.Uri.file(path))
        .then((file) => file.permissions);
};
exports.symbolicPermissions = symbolicPermissions;
/**
 * Returns the octal permissions for the supplied path.
 *
 * @param {string} path - Path to the file
 * @example
 * const permissions = await octalPermissions('src/file.ts');
 *
 * @returns {Promise<string | undefined>} - Octal permissions
 */
const octalPermissions = async (path) => {
    const file = await vscode_1.workspace.fs
        .stat(vscode_1.Uri.file(path))
        .then((file) => file.permissions);
    return file?.toString(8);
};
exports.octalPermissions = octalPermissions;
/**
 * Returns a boolean indicating whether the two supplied files are the same.
 *
 * @param {string} file1 - Path to the first file
 * @param {string} file2 - Path to the second file
 * @example
 * const isSame = await sameFile('src/file1.ts', 'src/file2.ts');
 *
 * @returns {Promise<boolean>} - Confirmation of the comparison
 */
const sameFile = async (file1, file2) => {
    const file1Info = await (0, exports.getFileInfo)(file1);
    const file2Info = await (0, exports.getFileInfo)(file2);
    return file1Info === file2Info;
};
exports.sameFile = sameFile;
/**
 * Sets the realpath for the supplied path.
 *
 * @param {string} path - Path to the file
 * @example
 * await setRealpath('src/file.ts');
 *
 * @returns {Promise<Uri | FileStat>} - Uri or FileStat object for the file
 */
const setRealpath = async (path) => {
    return (await vscode_1.workspace.fs.stat)
        ? await vscode_1.workspace.fs.stat(vscode_1.Uri.file(path))
        : await vscode_1.workspace
            .openTextDocument(vscode_1.Uri.file(path))
            .then((filename) => filename.uri);
};
exports.setRealpath = setRealpath;
/**
 * Returns the relative path from the workspace root to the supplied path.
 *
 * @param {string} path - Path to the file
 * @example
 * const relativePath = await getRelativePath('src/file.ts');
 *
 * @returns {Promise<string>} - Relative path
 */
const getRelativePath = async (path) => {
    return vscode_1.workspace.asRelativePath(path);
};
exports.getRelativePath = getRelativePath;
/**
 * Returns the realpath for the supplied path.
 *
 * @param {string} path - Path to the file
 * @example
 * const realpath = await getRealpath('src/file.ts');
 *
 * @returns {Promise<string>} - Realpath
 */
const getRealpath = async (path) => {
    return vscode_1.Uri.file(path).fsPath;
};
exports.getRealpath = getRealpath;
/**
 * Returns a boolean indicating whether the supplied path exists.
 *
 * @param {string} path - Path to the file or directory
 * @example
 * const fileExists = await exists('src/file.ts');
 *
 * @returns {Promise<boolean>} - Confirmation of the existence
 */
const exists = async (path) => {
    return (0, fs_1.existsSync)(path);
};
exports.exists = exists;
// isDirectory
/**
 * Returns a boolean indicating whether the supplied path is a directory.
 *
 * @param {string} path - Path to the file or directory
 * @example
 * const isDir = await isDirectory('src');
 *
 * @returns {Promise<boolean>} - Confirmation of the directory
 */
const isDirectory = async (path) => {
    return (await vscode_1.workspace.fs.stat(vscode_1.Uri.file(path))).type === 2;
};
exports.isDirectory = isDirectory;
//# sourceMappingURL=filesystem.helper.js.map