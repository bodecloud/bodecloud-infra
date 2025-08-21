"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.titleize = exports.singularize = exports.pluralize = exports.ordinalize = exports.ordinal = exports.dasherize = exports.isPluralizable = exports.humanize = exports.decamelize = exports.underscore = exports.pascalize = exports.camelize = void 0;
/**
 * Changes a string of words separated by spaces or underscores to camel case.
 *
 * @param {string} str - The string to camelize
 * @example
 * camelize('foo bar');
 *
 * @returns {string} - The camelized string
 */
const camelize = (str) => {
    return str
        .replace(/(?:^\w|[A-Z]|\b\w)/g, (word, index) => index === 0 ? word.toLowerCase() : word.toUpperCase())
        .replace(/\s+/g, '');
};
exports.camelize = camelize;
/**
 * Changes a string of words separated by spaces or underscores to pascal case.
 *
 * @param {string} str - The string to pascalize
 * @example
 * pascalize('foo bar');
 *
 * @returns {string} - The pascalized string
 */
const pascalize = (str) => {
    return str
        .replace(/(?:^\w|[A-Z]|\b\w)/g, (word) => word.toUpperCase())
        .replace(/\s+/g, '');
};
exports.pascalize = pascalize;
/**
 * Changes a string of words separated by spaces or camel or pascal case.
 *
 * @param {string} str - The string to underscore
 * @example
 * underscore('foo bar');
 *
 * @returns {string} - The underscored string
 */
const underscore = (str) => {
    return str
        .replace(/(?:^\w|[A-Z]|\b\w)/g, (word, index) => index === 0 ? word.toLowerCase() : `_${word.toLowerCase()}`)
        .replace(/\s+/g, '_');
};
exports.underscore = underscore;
/**
 * Changes a string of words separated by spaces or camel or pascal case to lowercase with underscores.
 *
 * @param {string} str - The string to decamelize
 * @example
 * decamelize('foo bar');
 *
 * @returns {string} - The decamelized string
 */
const decamelize = (str) => {
    return str
        .replace(/(?:^\w|[A-Z]|\b\w)/g, (word, index) => index === 0 ? word.toLowerCase() : `_${word.toLowerCase()}`)
        .replace(/\s+/g, '_');
};
exports.decamelize = decamelize;
/**
 * Changes a string of words separated by spaces or camel or pascal case to human readable form.
 *
 * @param {string} str - The string to humanize
 * @example
 * humanize('foo bar');
 *
 * @returns {string} - The humanized string
 */
const humanize = (str) => {
    return str
        .replace(/(?:^\w|[A-Z]|\b\w)/g, (word, index) => index === 0 ? word.toUpperCase() : ` ${word.toLowerCase()}`)
        .replace(/\s+/g, ' ');
};
exports.humanize = humanize;
/**
 * Checks if a string is pluralizable.
 *
 * @param {string} str - The string to check
 * @example
 * isPluralizable('foo');
 *
 * @returns {boolean} - Whether the string is pluralizable
 */
const isPluralizable = (str) => {
    return str.endsWith('s');
};
exports.isPluralizable = isPluralizable;
/**
 * Changes a string of words separated by spaces or camel or pascal case to lowercase with dashes.
 *
 * @param {string} str - The string to dasherize
 * @example
 * dasherize('foo bar');
 *
 * @returns {string} - The dasherized string
 */
const dasherize = (str) => {
    return str
        .replace(/(?:^\w|[A-Z]|\b\w)/g, (word, index) => index === 0 ? word.toLowerCase() : `-${word.toLowerCase()}`)
        .replace(/\s+/g, '-');
};
exports.dasherize = dasherize;
/**
 * Changes a number to its ordinal form.
 *
 * @param {number} num - The number to ordinalize
 * @example
 * ordinalize(1);
 *
 * @returns {string} - The ordinalized number
 */
const ordinal = (num) => {
    const j = num % 10;
    const k = num % 100;
    if (j === 1 && k !== 11) {
        return `${num}st`;
    }
    if (j === 2 && k !== 12) {
        return `${num}nd`;
    }
    if (j === 3 && k !== 13) {
        return `${num}rd`;
    }
    return `${num}th`;
};
exports.ordinal = ordinal;
/**
 * Changes a number to its ordinal form.
 *
 * @param {number} num - The number to ordinalize
 * @example
 * ordinalize(1);
 *
 * @returns {string} - The ordinalized number
 */
const ordinalize = (num) => {
    return `${num}${(0, exports.ordinal)(num)}`;
};
exports.ordinalize = ordinalize;
/**
 * Changes a string to its plural form.
 *
 * @param {string} str - The string to pluralize
 * @example
 * pluralize('foo');
 *
 * @returns {string} - The pluralized string
 */
const pluralize = (str) => {
    if (str.endsWith('y')) {
        return str.slice(0, -1) + 'ies';
    }
    if (str.endsWith('s')) {
        return str;
    }
    return str + 's';
};
exports.pluralize = pluralize;
/**
 * Changes a string to its singular form.
 *
 * @param {string} str - The string to singularize
 * @example
 * singularize('foos');
 *
 * @returns {string} - The singularized string
 */
const singularize = (str) => {
    if (str.endsWith('ies')) {
        return str.slice(0, -3) + 'y';
    }
    if (str.endsWith('s')) {
        return str.slice(0, -1);
    }
    return str;
};
exports.singularize = singularize;
/**
 * Changes a string to its title case form.
 *
 * @param {string} str - The string to titleize
 * @example
 * titleize('foo bar');
 *
 * @returns {string} - The titleized string
 */
const titleize = (str) => {
    return str
        .split(' ')
        .map((word) => word[0].toUpperCase() + word.slice(1))
        .join(' ');
};
exports.titleize = titleize;
//# sourceMappingURL=inflector.helper.js.map