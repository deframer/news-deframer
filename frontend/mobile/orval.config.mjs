export default {
  newsDeframerMobile: {
    input: {
      target: process.env.OPENAPI_URL,
    },
    output: {
      target: './src/services/generated/newsDeframerClient.gen.ts',
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
