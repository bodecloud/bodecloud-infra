"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.FileController = void 0;
const vscode_1 = require("vscode");
const helpers_1 = require("../helpers");
/**
 * The FileController class.
 *
 * @class
 * @classdesc The class that represents the example controller.
 * @export
 * @public
 * @example
 * const controller = new FileController(config);
 */
class FileController {
    // -----------------------------------------------------------------
    // Constructor
    // -----------------------------------------------------------------
    /**
     * Constructor for the FileController class.
     *
     * @constructor
     * @public
     * @memberof FileController
     */
    constructor(config) {
        this.config = config;
    }
    // -----------------------------------------------------------------
    // Methods
    // -----------------------------------------------------------------
    // Public methods
    /**
     * Generate a new class file.
     *
     * @function generateClass
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateClass(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateClass(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        // Get the type
        let type = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the type name'), 'E.g. class, dto, entity, model...', (type) => {
            if (!/[a-z]+/.test(type)) {
                return 'Invalid format!';
            }
            return;
        });
        if (!type) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `export class ${className}${(0, helpers_1.titleize)(type)} {}
`;
        type = type.length !== 0 ? `.${type}` : '';
        const filename = `${(0, helpers_1.dasherize)(className)}${type}.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new controller file.
     *
     * @function generateController
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateController(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateController(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the controller name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
} from '@nestjs/common';
import { Create${className}Dto } from './dto/create-${(0, helpers_1.dasherize)(className)}.dto';
import { Update${className}Dto } from './dto/update-${(0, helpers_1.dasherize)(className)}.dto';
import { ${className}Service } from './${(0, helpers_1.dasherize)(className)}.service';

@Controller('${(0, helpers_1.dasherize)(className)}s')
export class ${className}Controller {
  constructor(private readonly ${(0, helpers_1.dasherize)(className)}Service: ${className}Service) {}

  @Post()
  create(@Body() create${className}Dto: Create${className}Dto) {
    return this.${(0, helpers_1.dasherize)(className)}Service.create(create${className}Dto);
  }

  @Get()
  findAll() {
    return this.${(0, helpers_1.dasherize)(className)}Service.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.${(0, helpers_1.dasherize)(className)}Service.findOne(+id);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() update${className}Dto: Update${className}Dto) {
    return this.${(0, helpers_1.dasherize)(className)}Service.update(+id, update${className}Dto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.${(0, helpers_1.dasherize)(className)}Service.remove(+id);
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.controller.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
        await this.autoImport(folder, 'controllers', `${className}Controller`, `${(0, helpers_1.dasherize)(className)}.controller`);
    }
    /**
     * Generate a new decorator file.
     *
     * @function generateDecorator
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateDecorator(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateDecorator(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        const entityName = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the decorator name'), 'E.g. user, role, auth...', (name) => {
            if (!/^[A-Za-z-]{3,}$/.test(name)) {
                return 'Invalid format!';
            }
            return;
        });
        if (!entityName) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { SetMetadata } from '@nestjs/common';

export const ${entityName} = (...args: string[]) =>
  SetMetadata('${entityName}-decorator', args);
`;
        const filename = `${(0, helpers_1.dasherize)(entityName)}.decorator.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new dto file.
     *
     * @function generateDto
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateDto(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateDto(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the Dto class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { PartialType } from '@nestjs/mapped-types';
import { Create${className}Dto } from './create-${(0, helpers_1.dasherize)(className)}.dto';

export class Update${className}Dto extends PartialType(Create${className}Dto) {
  name: string;
  age: number;
}
`;
        const filename = `update-${(0, helpers_1.dasherize)(className)}.dto.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new exception filter file.
     *
     * @function generateExceptionFilter
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateExceptionFilter(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateExceptionFilter(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the exception filter name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import {
  ExceptionFilter,
  Catch,
  ArgumentsHost,
  HttpException,
} from '@nestjs/common';
import { Request, Response } from 'express';

@Catch(HttpException)
export class ${className}ExceptionFilter implements ExceptionFilter {
  catch(exception: HttpException, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();
    const status = exception.getStatus();

    response.status(status).json();
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.filter.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
        await this.autoImport(folder, 'providers', `${className}ExceptionFilter`, `${(0, helpers_1.dasherize)(className)}.filter`);
    }
    /**
     * Generate a new exception file.
     *
     * @function generateException
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateException(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateException(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the exception class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { HttpException, HttpStatus } from '@nestjs/common';

export class ${className}Exception extends HttpException {
  constructor() {
    super('${className}Exception', HttpStatus.FORBIDDEN);
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.exception.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new filter file.
     *
     * @function generateFilter
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateFilter(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateFilter(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the filter class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { ArgumentsHost, Catch, ExceptionFilter } from '@nestjs/common';

@Catch()
export class ${className}Filter<T> implements ExceptionFilter {
  catch(exception: T, host: ArgumentsHost) {}
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.filter.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
        await this.autoImport(folder, 'providers', `${className}Filter`, `${(0, helpers_1.dasherize)(className)}.filter`);
    }
    /**
     * Generate a new gateway file.
     *
     * @function generateGateway
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateGateway(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateGateway(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the gateway class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { SubscribeMessage, WebSocketGateway } from '@nestjs/websockets';

@WebSocketGateway()
export class ${className}Gateway {
  @SubscribeMessage('message')
  handleMessage(client: any, payload: any): string {
    return 'Hello world!';
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.gateway.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
        await this.autoImport(folder, 'providers', `${className}Gateway`, `${(0, helpers_1.dasherize)(className)}.gateway`);
    }
    /**
     * Generate a new guard file.
     *
     * @function generateGuard
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateGuard(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateGuard(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the guard class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { Observable } from 'rxjs';

@Injectable()
export class ${className}Guard implements CanActivate {
  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    return true;
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.guard.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new interceptor file.
     *
     * @function generateInterceptor
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateInterceptor(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateInterceptor(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the interceptor class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import {
  CallHandler,
  ExecutionContext,
  Injectable,
  NestInterceptor,
} from '@nestjs/common';
import { Observable } from 'rxjs';

@Injectable()
export class ${className}Interceptor implements NestInterceptor {
  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    return next.handle();
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.interceptor.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new interface file.
     *
     * @function generateInterface
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateInterface(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateInterface(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the interface class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        // Get the type
        let type = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the interface type'), 'E.g. interface, dto, entity, model...', (type) => {
            if (!/[a-z]+/.test(type)) {
                return 'Invalid format!';
            }
            return;
        });
        if (!type) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `export interface ${className}${(0, helpers_1.titleize)(type)} {}
`;
        type = type.length !== 0 ? `.${type}` : '';
        const filename = `${(0, helpers_1.dasherize)(className)}${type}.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new jwt guard file.
     *
     * @function generateJwtGuard
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateJwtGuard(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateJwtGuard(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the jwt guard class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import {
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { AuthGuard } from '@nestjs/passport';

@Injectable()
export class ${className}Guard extends AuthGuard('jwt') {
  canActivate(context: ExecutionContext) {
    return super.canActivate(context);
  }

  handleRequest(err, user, info) {
    if (err || !user) {
      throw err || new UnauthorizedException();
    }
    return user;
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.guard.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new jwt strategy file.
     *
     * @function generateJwtStrategy
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateJwtStrategy(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateJwtStrategy(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        const content = `import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy, 'jwt') {
  constructor(private configService: ConfigService) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKey: configService.get('secret'),
    });
  }

  async validate(payload: any) {
    return { userId: payload.sub, username: payload.username };
  }
}
`;
        const filename = 'jwt.strategy.ts';
        await (0, helpers_1.saveFile)(folder, filename, content);
        await this.autoImport(folder, 'providers', 'JwtStrategy', 'jwt.strategy');
    }
    /**
     * Generate a new logger file.
     *
     * @function generateLogger
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateLogger(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateLogger(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the logger class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { LoggerService } from '@nestjs/common';

export class ${className}Logger implements LoggerService {
  /**
   * Write a 'log' level log.
   */
  log(message: any, ...optionalParams: any[]) {}

  /**
   * Write an 'error' level log.
   */
  error(message: any, ...optionalParams: any[]) {}

  /**
   * Write a 'warn' level log.
   */
  warn(message: any, ...optionalParams: any[]) {}

  /**
   * Write a 'debug' level log.
   */
  debug?(message: any, ...optionalParams: any[]) {}

  /**
   * Write a 'verbose' level log.
   */
  verbose?(message: any, ...optionalParams: any[]) {}
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.logger.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new middleware file.
     *
     * @function generateMiddleware
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateMiddleware(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateMiddleware(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the middleware class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { Injectable, NestMiddleware } from '@nestjs/common';

@Injectable()
export class ${className}Middleware implements NestMiddleware {
  use(req: any, res: any, next: () => void) {
    next();
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.middleware.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new module file.
     *
     * @function generateModule
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateModule(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateModule(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the module class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { Module } from '@nestjs/common';

@Module({
  imports: [],
  controllers: [],
  providers: [],
  exports: []
})
export class ${className}Module {}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.module.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
        await this.autoImport(folder, 'imports', `${className}Module`, `${(0, helpers_1.dasherize)(className)}.module`);
    }
    /**
     * Generate a new pipe file.
     *
     * @function generatePipe
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generatePipe(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generatePipe(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the pipe class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { ArgumentMetadata, Injectable, PipeTransform } from '@nestjs/common';

@Injectable()
export class ${className}Pipe implements PipeTransform {
  transform(value: any, metadata: ArgumentMetadata) {
    return value;
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.pipe.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new provider file.
     *
     * @function generateProvider
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateProvider(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateProvider(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the provider class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { Injectable } from '@nestjs/common';

@Injectable()
export class ${className} {}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.provider.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new resolver file.
     *
     * @function generateResolver
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateResolver(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateResolver(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the resolver class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { Resolver } from '@nestjs/graphql';

@Resolver()
export class ${className}Resolver {}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.resolver.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Generate a new service file.
     *
     * @function generateService
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateService(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateService(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the service class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { Injectable } from '@nestjs/common';
import { Create${className}Dto } from './dto/create-${(0, helpers_1.dasherize)(className)}.dto';
import { Update${className}Dto } from './dto/update-${(0, helpers_1.dasherize)(className)}.dto';

@Injectable()
export class ${className}Service {
  create(create${className}Dto: Create${className}Dto) {
    return 'This action adds a new ${(0, helpers_1.dasherize)(className)}';
  }

  findAll() {
    return \`This action returns all ${(0, helpers_1.dasherize)(className)}s\`;
  }

  findOne(id: number) {
    return \`This action returns a #id ${(0, helpers_1.dasherize)(className)}\`;
  }

  update(id: number, update${className}Dto: Update${className}Dto) {
    return \`This action updates a #id ${(0, helpers_1.dasherize)(className)}\`;
  }

  remove(id: number) {
    return \`This action removes a #id ${(0, helpers_1.dasherize)(className)}\`;
  }
}
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.service.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
        await this.autoImport(folder, 'providers', `${className}Service`, `${(0, helpers_1.dasherize)(className)}.service`);
    }
    /**
     * Generate a new test file.
     *
     * @function generateTest
     * @param {Uri} [path] - The path to the folder.
     * @memberof FileController
     * @public
     * @async
     * @example
     * await generateTest(path);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async generateTest(path) {
        // Get the relative path
        const folderPath = path ? await (0, helpers_1.getRelativePath)(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the test class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        const content = `import { Test } from '@nestjs/testing';
import { ${className}Controller } from './.controller';
import { ${className}Service } from './.service';

describe('${className}Controller', () => {
  let ${className}Controller: ${className}Controller;
  let ${className}Service: ${className}Service;

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
        controllers: [${className}Controller],
        providers: [${className}Service],
      }).compile();

    ${className}Service = moduleRef.get<${className}Service>(${className}Service);
    ${className}Controller = moduleRef.get<${className}Controller>(${className}Controller);
  });

  describe('findAll', () => {
    it('should return an array of ${className}', async () => {
      const result = ['test'];
      jest.spyOn(${className}Service, 'findAll').mockImplementation(() => result);

      expect(await ${className}Controller.findAll()).toBe(result);
    });
  });
});
`;
        const filename = `${(0, helpers_1.dasherize)(className)}.spec.ts`;
        await (0, helpers_1.saveFile)(folder, filename, content);
    }
    /**
     * Creates a new custom element.
     *
     * @function generateCustomElement
     * @param {Uri} [path] - The path to the folder
     * @public
     * @async
     * @memberof FileController
     * @example
     * generateCustomElement();
     *
     * @returns {Promise<void>} - The result of the operation
     */
    async generateCustomElement(path) {
        // Get the relative path
        const folderPath = path ? vscode_1.workspace.asRelativePath(path.path) : '';
        const skipFolderConfirmation = this.config.skipFolderConfirmation;
        let folder;
        if (!folderPath || !skipFolderConfirmation) {
            // Get the path to the folder
            folder = await (0, helpers_1.getPath)(vscode_1.l10n.t('Enter the folder name'), 'Folder name. E.g. src, app...', folderPath, (path) => {
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
        if (this.config.templates.length === 0) {
            const message = vscode_1.l10n.t('The custom components list is empty. Please add custom components to the configuration');
            vscode_1.window.showErrorMessage(message);
            return;
        }
        const items = this.config.templates.map((item) => {
            return {
                label: item.name,
                description: item.description,
                detail: item.type,
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
        const template = this.config.templates.find((item) => item.name === option.label);
        if (!template) {
            const message = vscode_1.l10n.t('The template for the custom component does not exist. Please try again');
            vscode_1.window.showErrorMessage(message);
            return;
        }
        let content = Object(template).template.join('\n');
        // Get the class name
        const className = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the class name'), 'E.g. User, Role, Auth...', (name) => {
            if (!/^[A-Z][A-Za-z]{2,}$/.test(name)) {
                return 'Invalid format! Entity names MUST be declared in PascalCase.';
            }
            return;
        });
        if (!className) {
            const message = vscode_1.l10n.t('Operation cancelled!');
            (0, helpers_1.showError)(message);
            return;
        }
        content = content.replace(/{{ComponentName}}/g, className);
        if (content.includes('{{EntityName}}')) {
            // Get the class name
            const entityName = await (0, helpers_1.getName)(vscode_1.l10n.t('Enter the entity name'), 'E.g. user, role, auth...', (name) => {
                if (!/^[a-z][\w-]+$/.test(name)) {
                    return 'Invalid format! Entity names MUST be declared in camelCase.';
                }
                return;
            });
            if (!entityName) {
                const message = vscode_1.l10n.t('Operation cancelled!');
                (0, helpers_1.showError)(message);
                return;
            }
            content = content.replace(/{{EntityName}}/g, entityName);
        }
        const type = Object(template).type.length !== 0 ? `.${Object(template).type}` : '';
        const filename = `${(0, helpers_1.dasherize)(className)}${type}.ts`;
        (0, helpers_1.saveFile)(folder, filename, content);
    }
    // Private methods
    /**
     * Auto import functionality for files.
     *
     * @function autoImport
     * @param {string} directoryPath - The path to the folder.
     * @param {string} type - The type of the file.
     * @param {string} className - The class name.
     * @param {string} filename - The file name.
     * @memberof FileController
     * @private
     * @async
     * @example
     * await autoImport(directoryPath, type, className, filename);
     *
     * @returns {Promise<void>} The result of the operation.
     */
    async autoImport(directoryPath, type, className, filename) {
        try {
            if (!this.config.autoImport) {
                return; // Auto import is disabled, nothing to do
            }
            let files;
            const excludePatterns = `{${this.config.exclude.join(',')}}`;
            if (filename.includes('module')) {
                const tempPath = directoryPath.substring(0, directoryPath.lastIndexOf('/'));
                files = await vscode_1.workspace.findFiles(`${tempPath}/*.module.ts`, excludePatterns, 1);
                filename = `${directoryPath.substring(directoryPath.lastIndexOf('/') + 1)}/${filename}`;
            }
            else {
                files = await vscode_1.workspace.findFiles(`${directoryPath}/*.module.ts`, excludePatterns, 1);
            }
            if (files.length === 0) {
                const message = vscode_1.l10n.t('No module file found. Skipping auto-import!');
                (0, helpers_1.showError)(message);
                return; // No files found, nothing to do
            }
            const importRegex = new RegExp(`${type}: \\[`, 'g');
            const targetFile = files[0];
            const document = await vscode_1.workspace.openTextDocument(targetFile.path);
            const text = document.getText();
            for (let i = 0; i < document.lineCount; i++) {
                const line = document.lineAt(i).text;
                if (importRegex.test(line)) {
                    const position = document.positionAt(text.indexOf(line) + importRegex.lastIndex);
                    const edit = new vscode_1.WorkspaceEdit();
                    edit.insert(targetFile, position, `${className}, `);
                    edit.insert(targetFile, new vscode_1.Position(0, 0), `import { ${className} } from './${filename}';\n`);
                    await vscode_1.workspace.applyEdit(edit); // Applying edit asynchronously
                    await vscode_1.window.showTextDocument(document);
                    await vscode_1.commands.executeCommand('editor.action.formatDocument'); // Formatting the document
                    await vscode_1.commands.executeCommand('editor.action.organizeImports'); // Organizing the imports
                    await vscode_1.commands.executeCommand('workbench.action.files.saveAll'); // Saving the files
                    const folder = await (0, helpers_1.getRelativePath)(targetFile.path);
                    const message = vscode_1.l10n.t("Auto-import of {0} into '{1}' was successful!", [className, folder]);
                    (0, helpers_1.showMessage)(message);
                    return; // Import added, exiting function
                }
            }
        }
        catch (error) {
            const message = vscode_1.l10n.t('An error occurred during auto-import: {0}', Object(error).message ?? error);
            (0, helpers_1.showError)(message);
        }
    }
}
exports.FileController = FileController;
//# sourceMappingURL=file.controller.js.map