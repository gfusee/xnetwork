export function dontIndent(string: string): string {
    return string.replace(/(\n)[\r\t\f ]+(\S)/g, '$1$2');
}
