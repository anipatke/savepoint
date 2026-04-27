import { defineConfig } from "vitest/config";

// Keep this TypeScript copy as the typed source of truth for editors and tooling.
export default defineConfig({
  resolve: {
    preserveSymlinks: true,
  },
  test: {
    pool: "threads",
    include: ["test/**/*.test.ts"],
  },
});
