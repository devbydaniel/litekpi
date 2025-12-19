import { defineConfig } from 'orval'

export default defineConfig({
  litekpi: {
    input: '../backend/docs/swagger.json',
    output: {
      mode: 'single',
      target: './src/shared/api/generated/api.ts',
      schemas: './src/shared/api/generated/models',
      client: 'react-query',
      override: {
        mutator: {
          path: './src/shared/api/client.ts',
          name: 'customInstance',
        },
      },
    },
  },
})
