import { defineConfig } from "vitest/config";

export default defineConfig({
  resolve: {
    preserveSymlinks: true,
  },
  test: {
    pool: "threads",
    include: ["test/**/*.test.ts"],
  },
});
