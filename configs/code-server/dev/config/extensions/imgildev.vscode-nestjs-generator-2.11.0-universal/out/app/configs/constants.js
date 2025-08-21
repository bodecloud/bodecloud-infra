"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ORM = exports.SKIP_FOLDER_CONFIRMATION = exports.AUTO_IMPORT = exports.ACTIVATE_MENU = exports.CUSTOM_TEMPLATES = exports.CUSTOM_COMMANDS = exports.SHOW_PATH = exports.WATCH = exports.EXCLUDE = exports.INCLUDE = exports.EXTENSION_PAYPAL_URL = exports.EXTENSION_SPONSOR_URL = exports.EXTENSION_MARKETPLACE_URL = exports.EXTENSION_REPOSITORY_URL = exports.USER_PUBLISHER = exports.USER_NAME = exports.EXTENSION_DISPLAY_NAME = exports.EXTENSION_NAME = exports.EXTENSION_ID = void 0;
/**
 * EXTENSION_ID: The unique identifier of the extension.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(EXTENSION_ID);
 *
 * @returns {string} - The unique identifier of the extension
 */
exports.EXTENSION_ID = 'nestjs';
/**
 * EXTENSION_NAME: The repository ID of the extension.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(EXTENSION_NAME);
 *
 * @returns {string} - The repository ID of the extension
 */
exports.EXTENSION_NAME = 'vscode-nestjs-generator';
/**
 * EXTENSION_DISPLAY_NAME: The name of the extension.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(EXTENSION_DISPLAY_NAME);
 *
 * @returns {string} - The name of the extension
 */
exports.EXTENSION_DISPLAY_NAME = 'NestJS File Generator';
/**
 * USER_NAME: The githubUsername of the extension.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(USER_NAME);
 *
 * @returns {string} - The githubUsername of the extension
 */
exports.USER_NAME = 'ManuelGil';
/**
 * USER_PUBLISHER: The publisher of the extension.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(USER_PUBLISHER);
 *
 * @returns {string} - The publisher of the extension
 */
exports.USER_PUBLISHER = 'imgildev';
/**
 * EXTENSION_REPOSITORY_URL: The repository URL of the extension.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(EXTENSION_REPOSITORY_URL);
 *
 * @returns {string} - The repository URL of the extension
 */
exports.EXTENSION_REPOSITORY_URL = `https://github.com/${exports.USER_NAME}/${exports.EXTENSION_NAME}`;
/**
 * MARKETPLACE_URL: The marketplace URL of the extension.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(MARKETPLACE_URL);
 *
 * @returns {string} - The marketplace URL of the extension
 */
exports.EXTENSION_MARKETPLACE_URL = `https://marketplace.visualstudio.com/items?itemName=${exports.USER_PUBLISHER}.${exports.EXTENSION_NAME}`;
/**
 * EXTENSION_SPONSOR_URL: The sponsor URL of the extension.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(EXTENSION_SPONSOR_URL);
 *
 * @returns {string} - The sponsor URL of the extension
 */
exports.EXTENSION_SPONSOR_URL = 'https://github.com/sponsors/ManuelGil';
/**
 * EXTENSION_PAYPAL_URL: The PayPal URL of the extension.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(EXTENSION_PAYPAL_URL);
 *
 * @returns {string} - The PayPal URL of the extension
 */
exports.EXTENSION_PAYPAL_URL = 'https://www.paypal.com/paypalme/ManuelFGil';
/**
 * INCLUDE: The files to include.
 * @type {string[]}
 * @public
 * @memberof Constants
 * @example
 * console.log(INCLUDE);
 *
 * @returns {string[]} - The files to include
 */
exports.INCLUDE = ['ts'];
/**
 * EXCLUDE: The files to exclude.
 * @type {string[]}
 * @public
 * @memberof Constants
 * @example
 * console.log(EXCLUDE);
 *
 * @returns {string[]} - The files to exclude
 */
exports.EXCLUDE = [
    '**/node_modules/**',
    '**/dist/**',
    '**/out/**',
    '**/build/**',
    '**/.*/**',
];
/**
 * WATCH: The files to watch.
 * @type {string[]}
 * @public
 * @memberof Constants
 * @example
 * console.log(WATCH);
 *
 * @returns {string[]} - The files to watch
 */
exports.WATCH = ['controllers', 'dtos', 'services'];
/**
 * SHOW_PATH: Whether to show the path or not.
 * @type {boolean}
 * @public
 * @memberof Constants
 * @example
 * console.log(SHOW_PATH);
 *
 * @returns {boolean} - Whether to show the path or not
 */
exports.SHOW_PATH = true;
/**
 * CUSTOM_COMMANDS: The custom commands.
 * @type {object[]}
 * @public
 * @memberof Constants
 * @example
 * console.log(CUSTOM_COMMANDS);
 *
 * @returns {object[]} - The custom commands
 */
exports.CUSTOM_COMMANDS = [
    {
        'name': 'Template 1',
        'command': 'nest g co',
        'args': '--flat',
    },
    {
        'name': 'Template 2',
        'command': 'nest g co',
        'args': '--no-flat',
    },
];
/**
 * CUSTOM_TEMPLATES: The custom templates.
 * @type {object[]}
 * @public
 * @memberof Constants
 * @example
 * console.log(CUSTOM_TEMPLATES);
 *
 * @returns {object[]} - The custom templates
 */
exports.CUSTOM_TEMPLATES = [
    {
        'name': 'Custom Service',
        'description': 'Generate a custom service',
        'type': 'service',
        'template': [
            "import { Injectable } from '@nestjs/common';",
            '',
            '@Injectable()',
            'export class CustomService {',
            '}',
        ],
    },
];
/**
 * ACTIVATE_MENU: Whether to show the menu or not.
 * @type {object}
 * @public
 * @memberof Constants
 * @example
 * console.log(ACTIVATE_MENU);
 *
 * @returns {object} - Whether to show the menu or not
 */
exports.ACTIVATE_MENU = {
    file: {
        class: true,
        controller: true,
        decorator: true,
        dto: true,
        exception: true,
        exceptionFilter: true,
        filter: true,
        gateway: true,
        guard: true,
        interceptor: true,
        interface: true,
        jwtGuard: true,
        jwtStrategy: true,
        middleware: true,
        logger: true,
        module: true,
        pipe: true,
        provider: true,
        resolver: true,
        service: true,
        test: true,
        template: true,
    },
    terminal: {
        controller: true,
        gateway: true,
        library: true,
        module: true,
        provider: true,
        resolver: true,
        resource: true,
        service: true,
        custom: true,
        start: true,
        startDev: true,
        startDebug: true,
        startProd: true,
    },
};
/**
 * AUTO_IMPORT: The auto import setting.
 * @type {boolean}
 * @public
 * @memberof Constants
 * @example
 * console.log(AUTO_IMPORT);
 *
 * @returns {boolean} - The auto import setting
 */
exports.AUTO_IMPORT = true;
/**
 * SKIP_FOLDER_CONFIRMATION: Whether to skip the folder confirmation or not.
 * @type {boolean}
 * @public
 * @memberof Constants
 * @example
 * console.log(SKIP_FOLDER_CONFIRMATION);
 *
 * @returns {boolean} - Whether to skip the folder confirmation or not
 */
exports.SKIP_FOLDER_CONFIRMATION = false;
/**
 * ORM: The orm.
 * @type {string}
 * @public
 * @memberof Constants
 * @example
 * console.log(ORM);
 *
 * @returns {string} - The orm
 */
exports.ORM = 'typeorm';
//# sourceMappingURL=constants.js.map