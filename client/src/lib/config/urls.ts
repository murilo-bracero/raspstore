interface CoreURLs {
  loginPAM: string;
  files: string;
  download: string;
  upload: string;
}

const coreHost = import.meta.env.VITE_CORE_HOST;

export let coreURLs: CoreURLs = {
  loginPAM: coreHost + '/v1/login',
  files: coreHost + '/v1/files',
  download: coreHost + '/v1/downloads',
  upload: coreHost + '/v1/uploads'
};
