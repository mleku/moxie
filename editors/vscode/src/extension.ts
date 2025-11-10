// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import * as vscode from 'vscode';
import {
    LanguageClient,
    LanguageClientOptions,
    ServerOptions,
    TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient | undefined;

export function activate(context: vscode.ExtensionContext) {
    console.log('Moxie extension activating...');

    // Get configuration
    const config = vscode.workspace.getConfiguration('moxie');
    const moxiePath = config.get<string>('moxiePath', 'moxie');
    const lspEnabled = config.get<boolean>('lsp.enabled', true);

    // Start LSP server if enabled
    if (lspEnabled) {
        startLanguageServer(moxiePath, config, context);
    }

    // Register commands
    registerCommands(context, moxiePath, config);

    // Register format on save
    if (config.get<boolean>('formatOnSave', true)) {
        context.subscriptions.push(
            vscode.workspace.onWillSaveTextDocument(event => {
                const document = event.document;
                if (document.languageId === 'moxie') {
                    event.waitUntil(
                        vscode.commands.executeCommand('editor.action.formatDocument')
                    );
                }
            })
        );
    }

    // Register vet on save
    if (config.get<boolean>('vetOnSave', false)) {
        context.subscriptions.push(
            vscode.workspace.onDidSaveTextDocument(document => {
                if (document.languageId === 'moxie') {
                    vscode.commands.executeCommand('moxie.vet');
                }
            })
        );
    }

    console.log('Moxie extension activated!');
}

export function deactivate(): Thenable<void> | undefined {
    if (!client) {
        return undefined;
    }
    return client.stop();
}

function startLanguageServer(
    moxiePath: string,
    config: vscode.WorkspaceConfiguration,
    context: vscode.ExtensionContext
) {
    // Server options
    const serverOptions: ServerOptions = {
        command: moxiePath,
        args: ['lsp'],
        transport: TransportKind.stdio
    };

    // Client options
    const clientOptions: LanguageClientOptions = {
        documentSelector: [
            { scheme: 'file', language: 'moxie' }
        ],
        synchronize: {
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.{mx,x}')
        }
    };

    // Create and start client
    client = new LanguageClient(
        'moxie-lsp',
        'Moxie Language Server',
        serverOptions,
        clientOptions
    );

    // Start the client
    client.start().then(() => {
        console.log('Moxie LSP client started');
    }).catch(err => {
        console.error('Failed to start Moxie LSP client:', err);
        vscode.window.showErrorMessage(
            `Failed to start Moxie Language Server: ${err.message}`
        );
    });

    context.subscriptions.push(client);
}

function registerCommands(
    context: vscode.ExtensionContext,
    moxiePath: string,
    config: vscode.WorkspaceConfiguration
) {
    // Build command
    context.subscriptions.push(
        vscode.commands.registerCommand('moxie.build', async () => {
            const terminal = vscode.window.createTerminal('Moxie Build');
            terminal.show();
            terminal.sendText(`${moxiePath} build`);
        })
    );

    // Run command
    context.subscriptions.push(
        vscode.commands.registerCommand('moxie.run', async () => {
            const editor = vscode.window.activeTextEditor;
            if (!editor) {
                vscode.window.showErrorMessage('No active file to run');
                return;
            }

            const filePath = editor.document.uri.fsPath;
            const terminal = vscode.window.createTerminal('Moxie Run');
            terminal.show();
            terminal.sendText(`${moxiePath} run ${filePath}`);
        })
    );

    // Test command
    context.subscriptions.push(
        vscode.commands.registerCommand('moxie.test', async () => {
            const terminal = vscode.window.createTerminal('Moxie Test');
            terminal.show();
            terminal.sendText(`${moxiePath} test ./...`);
        })
    );

    // Format command
    context.subscriptions.push(
        vscode.commands.registerCommand('moxie.fmt', async () => {
            const editor = vscode.window.activeTextEditor;
            if (!editor) {
                return;
            }

            await vscode.commands.executeCommand('editor.action.formatDocument');
        })
    );

    // Vet command
    context.subscriptions.push(
        vscode.commands.registerCommand('moxie.vet', async () => {
            const terminal = vscode.window.createTerminal('Moxie Vet');
            terminal.show();
            terminal.sendText(`${moxiePath} vet ./...`);
        })
    );

    // Clean command
    context.subscriptions.push(
        vscode.commands.registerCommand('moxie.clean', async () => {
            const terminal = vscode.window.createTerminal('Moxie Clean');
            terminal.show();
            terminal.sendText(`${moxiePath} clean`);
        })
    );
}
