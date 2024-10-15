const coreHost = process.env.RS_CORE_HOST || import.meta.env.VITE_CORE_HOST;

export let coreURLs = {
  files: coreHost + '/v1/files',
  download: coreHost + '/v1/downloads',
  upload: coreHost + '/v1/uploads'
};
