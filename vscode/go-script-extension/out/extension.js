"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const vscode = require("vscode");
const cp = require("child_process");
const path = require("path");
function activate(context) {
    console.log('Go-Script extension is now active!');
    // Register commands
    const runCommand = vscode.commands.registerCommand('go-script.run', runGoScript);
    const buildCommand = vscode.commands.registerCommand('go-script.build', buildGoScript);
    const debugCommand = vscode.commands.registerCommand('go-script.debug', debugGoScript);
    const checkCommand = vscode.commands.registerCommand('go-script.check', checkSyntax);
    context.subscriptions.push(runCommand, buildCommand, debugCommand, checkCommand);
    // Set up diagnostics
    const diagnosticCollection = vscode.languages.createDiagnosticCollection('go-script');
    context.subscriptions.push(diagnosticCollection);
    // Watch for file changes and check syntax
    const watcher = vscode.workspace.createFileSystemWatcher('**/*.gos');
    watcher.onDidChange(uri => checkFileForErrors(uri, diagnosticCollection));
    watcher.onDidCreate(uri => checkFileForErrors(uri, diagnosticCollection));
    context.subscriptions.push(watcher);
    // Check syntax on save
    vscode.workspace.onDidSaveTextDocument(document => {
        if (document.languageId === 'go-script') {
            checkFileForErrors(document.uri, diagnosticCollection);
        }
    });
    // Status bar item
    const statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100);
    statusBarItem.text = "$(play) Go-Script";
    statusBarItem.tooltip = "Click to run Go-Script file";
    statusBarItem.command = 'go-script.run';
    statusBarItem.show();
    context.subscriptions.push(statusBarItem);
}
exports.activate = activate;
async function runGoScript(uri) {
    const filePath = getActiveFilePath(uri);
    if (!filePath)
        return;
    const terminal = getOrCreateTerminal();
    const gosPath = getGosPath();
    terminal.show();
    terminal.sendText(`${gosPath} run "${filePath}"`);
    vscode.window.showInformationMessage(`Running Go-Script: ${path.basename(filePath)}`);
}
async function buildGoScript(uri) {
    const filePath = getActiveFilePath(uri);
    if (!filePath)
        return;
    const terminal = getOrCreateTerminal();
    const gosPath = getGosPath();
    terminal.show();
    terminal.sendText(`${gosPath} build "${filePath}"`);
    vscode.window.showInformationMessage(`Building Go-Script: ${path.basename(filePath)}`);
}
async function debugGoScript(uri) {
    const filePath = getActiveFilePath(uri);
    if (!filePath)
        return;
    const terminal = getOrCreateTerminal();
    const gosPath = getGosPath();
    terminal.show();
    terminal.sendText(`${gosPath} debug "${filePath}"`);
    vscode.window.showInformationMessage(`Debugging Go-Script: ${path.basename(filePath)}`);
}
async function checkSyntax(uri) {
    const filePath = getActiveFilePath(uri);
    if (!filePath)
        return;
    const diagnosticCollection = vscode.languages.createDiagnosticCollection('go-script');
    await checkFileForErrors(vscode.Uri.file(filePath), diagnosticCollection);
}
async function checkFileForErrors(uri, diagnosticCollection) {
    const config = vscode.workspace.getConfiguration('go-script');
    if (!config.get('enableErrorChecking', true)) {
        return;
    }
    const gosPath = getGosPath();
    const filePath = uri.fsPath;
    try {
        const { stdout, stderr } = await execAsync(`${gosPath} build "${filePath}"`);
        // Clear previous diagnostics
        diagnosticCollection.delete(uri);
        if (stderr) {
            const diagnostics = parseErrors(stderr);
            diagnosticCollection.set(uri, diagnostics);
        }
        // Show success message if no errors
        if (!stderr && stdout.includes('Success:')) {
            vscode.window.setStatusBarMessage('$(check) Go-Script: No errors', 3000);
        }
    }
    catch (error) {
        const diagnostics = parseErrors(error.stderr || error.message);
        diagnosticCollection.set(uri, diagnostics);
        vscode.window.setStatusBarMessage('$(error) Go-Script: Compilation errors', 5000);
    }
}
function parseErrors(errorOutput) {
    const diagnostics = [];
    const lines = errorOutput.split('\n');
    for (const line of lines) {
        // Parse Go-Script error format
        const match = line.match(/^(.+):(\d+):(\d+):\s*(error|warning):\s*(.+)$/);
        if (match) {
            const [, , lineStr, colStr, severity, message] = match;
            const lineNum = parseInt(lineStr) - 1; // VS Code uses 0-based line numbers
            const colNum = parseInt(colStr) - 1;
            const range = new vscode.Range(lineNum, colNum, lineNum, colNum + 10);
            const diagnostic = new vscode.Diagnostic(range, message, severity === 'error' ? vscode.DiagnosticSeverity.Error : vscode.DiagnosticSeverity.Warning);
            diagnostics.push(diagnostic);
        }
        // Parse Go-Script compilation failed format
        if (line.includes('Compilation failed:')) {
            const errorLines = lines.slice(lines.indexOf(line) + 1);
            for (const errorLine of errorLines) {
                if (errorLine.match(/^\d+\./)) {
                    const message = errorLine.replace(/^\d+\.\s*/, '');
                    const diagnostic = new vscode.Diagnostic(new vscode.Range(0, 0, 0, 10), message, vscode.DiagnosticSeverity.Error);
                    diagnostics.push(diagnostic);
                }
            }
        }
    }
    return diagnostics;
}
function getActiveFilePath(uri) {
    if (uri) {
        return uri.fsPath;
    }
    const activeEditor = vscode.window.activeTextEditor;
    if (!activeEditor) {
        vscode.window.showErrorMessage('No active Go-Script file');
        return undefined;
    }
    if (activeEditor.document.languageId !== 'go-script') {
        vscode.window.showErrorMessage('Active file is not a Go-Script file');
        return undefined;
    }
    return activeEditor.document.fileName;
}
function getOrCreateTerminal() {
    const existingTerminal = vscode.window.terminals.find(t => t.name === 'Go-Script');
    if (existingTerminal) {
        return existingTerminal;
    }
    return vscode.window.createTerminal('Go-Script');
}
function getGosPath() {
    const config = vscode.workspace.getConfiguration('go-script');
    return config.get('gosPath', 'gos');
}
function execAsync(command) {
    return new Promise((resolve, reject) => {
        cp.exec(command, (error, stdout, stderr) => {
            if (error) {
                reject({ ...error, stderr });
            }
            else {
                resolve({ stdout, stderr });
            }
        });
    });
}
function deactivate() {
    console.log('Go-Script extension deactivated');
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map