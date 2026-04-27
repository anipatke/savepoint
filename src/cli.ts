import { version } from "./version.js";

function printHelp(): void {
  process.stdout.write(`savepoint v${version}\n`);
  process.stdout.write("This scaffold is not implemented yet.\n");
}

function main(): void {
  const [, , command] = process.argv;

  if (command === "--version" || command === "-v") {
    process.stdout.write(`${version}\n`);
    return;
  }

  printHelp();
}

main();
