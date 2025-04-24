// run/index.js
exports.handler = async (event) => {
    console.log('Função Lambda "dummy" invocada pelo SAM Local.');
    return {
      statusCode: 200,
      body: 'Proxy para container Docker (SAM Local)',
    };
  };