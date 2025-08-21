# NestJS Swagger (OpenAPI) Snippets

[![Visual Studio Marketplace Version](https://img.shields.io/visual-studio-marketplace/v/imgildev.vscode-nestjs-swagger-snippets?style=for-the-badge&label=VS%20Marketplace&logo=visual-studio-code)](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-nestjs-swagger-snippets)
[![Visual Studio Marketplace Installs](https://img.shields.io/visual-studio-marketplace/i/imgildev.vscode-nestjs-swagger-snippets?style=for-the-badge&logo=visual-studio-code)](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-nestjs-swagger-snippets)
[![Visual Studio Marketplace Downloads](https://img.shields.io/visual-studio-marketplace/d/imgildev.vscode-nestjs-swagger-snippets?style=for-the-badge&logo=visual-studio-code)](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-nestjs-swagger-snippets)
[![Visual Studio Marketplace Rating](https://img.shields.io/visual-studio-marketplace/r/imgildev.vscode-nestjs-swagger-snippets?style=for-the-badge&logo=visual-studio-code)](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-nestjs-swagger-snippets&ssr=false#review-details)
[![GitHub Repo stars](https://img.shields.io/github/stars/ManuelGil/vscode-nestjs-swagger-snippets?style=for-the-badge&logo=github)](https://github.com/ManuelGil/vscode-nestjs-swagger-snippets)
[![GitHub license](https://img.shields.io/github/license/ManuelGil/vscode-nestjs-swagger-snippets?style=for-the-badge&logo=github)](https://github.com/ManuelGil/vscode-nestjs-swagger-snippets/blob/main/LICENSE)

> Quick, consistent `@nestjs/swagger` decorator and `DocumentBuilder` snippets for faster API documentation in NestJS.

## Overview

This Visual Studio Code extension provides ready-to-use snippets for the `@nestjs/swagger` package and the `DocumentBuilder` API. The snippets help you document controllers, DTOs and API metadata faster and more consistently.

## Requirements

- Visual Studio Code 1.46.0 or later

## Installation

1. Open Visual Studio Code (or a compatible editor).
2. Open the **Extensions** view (`Ctrl+Shift+X` / `⌘+Shift+X`).
3. Search for **NestJS Swagger (OpenAPI) Snippets** or install directly from the [Marketplace page](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-nestjs-swagger-snippets).
4. Click **Install** and reload the editor if prompted.

## Usage

Type part of a snippet name and press `Tab` or `Enter` to expand it.

> Note: Many snippets have two triggers — a namespaced one (e.g. `ns_swagger_deco_api_operation`) and a more natural alias (e.g. `@ApiOperation` or `setTitle`). You can use either.

### Common snippets

| Snippet                                              | Purpose                            |
| ---------------------------------------------------- | ---------------------------------- |
| `ns_swagger_deco_api_basic_auth`                     | `@ApiBasicAuth`                    |
| `ns_swagger_deco_api_bearer_auth`                    | `@ApiBearerAuth`                   |
| `ns_swagger_deco_api_body`                           | `@ApiBody`                         |
| `ns_swagger_deco_api_consumes`                       | `@ApiConsumes`                     |
| `ns_swagger_deco_api_cookie_auth`                    | `@ApiCookieAuth`                   |
| `ns_swagger_deco_api_exclude_controller`             | `@ApiExcludeController`            |
| `ns_swagger_deco_api_operation`                      | `@ApiOperation`                    |
| `ns_swagger_deco_api_param`                          | `@ApiParam`                        |
| `ns_swagger_deco_api_property`                       | `@ApiProperty`                     |
| `ns_swagger_deco_api_query`                          | `@ApiQuery`                        |
| `ns_swagger_deco_api_tags`                           | `@ApiTags`                         |
| `ns_swagger_deco_api_ok_response`                    | `@ApiOkResponse`                   |
| `ns_swagger_deco_api_bad_request_response`           | `@ApiBadRequestResponse`           |
| `ns_swagger_deco_api_unauthorized_response`          | `@ApiUnauthorizedResponse`         |
| `ns_swagger_deco_api_not_found_response`             | `@ApiNotFoundResponse`             |
| `ns_swagger_deco_api_internal_server_error_response` | `@ApiInternalServerErrorResponse`  |
| `ns_swagger_deco_api_default_response`               | `@ApiDefaultResponse`              |
| `ns_swagger_set_title`                               | `setTitle` (DocumentBuilder)       |
| `ns_swagger_set_description`                         | `setDescription` (DocumentBuilder) |
| `ns_swagger_set_version`                             | `setVersion` (DocumentBuilder)     |
| `ns_swagger_set_contact`                             | `setContact` (DocumentBuilder)     |
| `ns_swagger_set_license`                             | `setLicense` (DocumentBuilder)     |
| `ns_swagger_add_server`                              | `addServer` (DocumentBuilder)      |
| `ns_swagger_add_tag`                                 | `addTag` (DocumentBuilder)         |
| `ns_swagger_add_security`                            | `addSecurity` (DocumentBuilder)    |
| `ns_swagger_build`                                   | `build()` (DocumentBuilder)        |

## Contributing

Contributions to the NestJS Swagger (OpenAPI) Snippets are welcome and appreciated. To contribute:

1. Fork the [GitHub repository](https://github.com/ManuelGil/vscode-nestjs-swagger-snippets).
2. Create a new branch for your feature or fix:

   ```bash
   git checkout -b feature/your-feature
   ```

3. Make your changes, commit them, and push to your fork.
4. Submit a Pull Request targeting the `main` branch.

Before contributing, please review the [Contribution Guidelines](https://github.com/ManuelGil/vscode-nestjs-swagger-snippets/blob/main/CONTRIBUTING.md) for coding standards, testing, and commit message conventions. If you encounter a bug or wish to request a new feature, please open an Issue.

## Changelog

For a complete list of changes, see the [CHANGELOG.md](https://github.com/ManuelGil/vscode-nestjs-swagger-snippets/blob/main/CHANGELOG.md).

## Authors

- **Manuel Gil** - _Owner_ - [@ManuelGil](https://github.com/ManuelGil)

For a complete list of contributors, please refer to the [contributors](https://github.com/ManuelGil/vscode-nestjs-swagger-snippets/contributors) page.

## Follow Me

- **GitHub**: [![GitHub followers](https://img.shields.io/github/followers/ManuelGil?style=for-the-badge\&logo=github)](https://github.com/ManuelGil)
- **X (formerly Twitter)**: [![X Follow](https://img.shields.io/twitter/follow/imgildev?style=for-the-badge\&logo=x)](https://twitter.com/imgildev)

## Other Extensions

- **[Auto Barrel](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-auto-barrel)**
  Automatically generates and maintains barrel (`index.ts`) files for your TypeScript projects.

- **[Angular File Generator](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-angular-generator)**
  Generates boilerplate and navigates your Angular (9→20+) project from within the editor, with commands for components, services, directives, modules, pipes, guards, reactive snippets, and JSON2TS transformations.

- **[NestJS File Generator](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-nestjs-generator)**
  Simplifies creation of controllers, services, modules, and more for NestJS projects, with custom commands and Swagger snippets.

- **[NestJS Snippets](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-nestjs-snippets-extension)**
  Ready-to-use code patterns for creating controllers, services, modules, DTOs, filters, interceptors, and more in NestJS.

- **[T3 Stack / NextJS / ReactJS File Generator](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-nextjs-generator)**
  Automates file creation (components, pages, hooks, API routes, etc.) in T3 Stack (Next.js, React) projects and can start your dev server from VSCode.

- **[Drizzle ORM Snippets](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-drizzle-snippets)**
  Collection of code snippets to speed up Drizzle ORM usage, defines schemas, migrations, and common database operations in TypeScript/JavaScript.

- **[CodeIgniter 4 Spark](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-codeigniter4-spark)**
  Scaffolds controllers, models, migrations, libraries, and CLI commands in CodeIgniter 4 projects using Spark, directly from the editor.

- **[CodeIgniter 4 Snippets](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-codeigniter4-snippets)**
  Snippets for accelerating development with CodeIgniter 4, including controllers, models, validations, and more.

- **[CodeIgniter 4 Shield Snippets](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-codeigniter4-shield-snippets)**
  Snippets tailored to CodeIgniter 4 Shield for faster authentication and security-related code.

- **[Mustache Template Engine - Snippets & Autocomplete](https://marketplace.visualstudio.com/items?itemName=imgildev.vscode-mustache-snippets)**
  Snippets and autocomplete support for Mustache templates, making HTML templating faster and more reliable.

## Recommended Browser Extension

For developers who work with `.vsix` files for offline installations or distribution, the complementary [**One-Click VSIX**](https://chromewebstore.google.com/detail/imojppdbcecfpeafjagncfplelddhigc?utm_source=item-share-cb) extension is recommended, available for both Chrome and Firefox.

> **One-Click VSIX** integrates a direct "Download Extension" button into each VSCode Marketplace page, ensuring the file is saved with the `.vsix` extension, even if the server provides a `.zip` archive. This simplifies the process of installing or sharing extensions offline by eliminating the need for manual file renaming.

- [Get One-Click VSIX for Chrome &rarr;](https://chromewebstore.google.com/detail/imojppdbcecfpeafjagncfplelddhigc?utm_source=item-share-cb)
- [Get One-Click VSIX for Firefox &rarr;](https://addons.mozilla.org/es-ES/firefox/addon/one-click-vsix/)

## License

This project is licensed under the **MIT License**. See the [LICENSE](https://github.com/ManuelGil/vscode-nestjs-swagger-snippets/blob/main/LICENSE) file for full details.
