export default {
  newsDeframerWeb: {
    input: {
      target: process.env.OPENAPI_URL,
    },
    output: {
      target: './src/ndf/generated/newsDeframerClient.gen.ts',
      client: 'fetch',
      mode: 'single',
      clean: true,
      override: {
        fetch: {
          useRuntimeFetcher: true,
        },
      },
    },
  },
};
