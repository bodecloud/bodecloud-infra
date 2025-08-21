"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.FeedbackController = void 0;
const vscode_1 = require("vscode");
const configs_1 = require("../configs");
/**
 * The FeedbackController class.
 *
 * @class
 * @classdesc The class that represents the feedback controller.
 * @export
 * @public
 * @example
 * const controller = new FeedbackController();
 */
class FeedbackController {
    // -----------------------------------------------------------------
    // Constructor
    // -----------------------------------------------------------------
    /**
     * Constructor for the FeedbackController class.
     *
     * @constructor
     * @public
     * @memberof FeedbackController
     */
    constructor() { }
    // -----------------------------------------------------------------
    // Methods
    // -----------------------------------------------------------------
    // Public methods
    /**
     * The aboutUs method.
     *
     * @function aboutUs
     * @public
     * @memberof FeedbackController
     *
     * @returns {void} - No return value
     */
    aboutUs() {
        vscode_1.env.openExternal(vscode_1.Uri.parse(configs_1.EXTENSION_MARKETPLACE_URL));
    }
    /**
     * The reportIssues method.
     *
     * @function reportIssues
     * @public
     * @memberof FeedbackController
     *
     * @returns {void} - No return value
     */
    reportIssues() {
        vscode_1.env.openExternal(vscode_1.Uri.parse(`${configs_1.EXTENSION_REPOSITORY_URL}/issues`));
    }
    /**
     * The rateUs method.
     *
     * @function rateUs
     * @public
     * @memberof FeedbackController
     *
     * @returns {void} - No return value
     */
    rateUs() {
        vscode_1.env.openExternal(vscode_1.Uri.parse(`${configs_1.EXTENSION_MARKETPLACE_URL}&ssr=false#review-details`));
    }
    /**
     * The supportUs method.
     *
     * @function supportUs
     * @public
     * @async
     * @memberof FeedbackController
     *
     * @returns {Promise<void>} - The promise that resolves with no value
     */
    async supportUs() {
        // Create the actions
        const actions = [
            { title: vscode_1.l10n.t('Become a Sponsor') },
            { title: vscode_1.l10n.t('Donate via PayPal') },
        ];
        // Show the message
        const message = vscode_1.l10n.t('Although {0} is offered at no cost, your support is deeply appreciated if you find it beneficial. Thank you for considering!', configs_1.EXTENSION_DISPLAY_NAME);
        const option = await vscode_1.window.showInformationMessage(message, ...actions);
        // Handle the actions
        switch (option?.title) {
            case actions[0].title:
                vscode_1.env.openExternal(vscode_1.Uri.parse(configs_1.EXTENSION_SPONSOR_URL));
                break;
            case actions[1].title:
                vscode_1.env.openExternal(vscode_1.Uri.parse(configs_1.EXTENSION_PAYPAL_URL));
                break;
        }
    }
}
exports.FeedbackController = FeedbackController;
//# sourceMappingURL=feedback.controller.js.map