/* eslint-env node */
'use strict';

module.exports = {
  name: 'startup',
  isDevelopingAddon: function() {
    return true;
  },
  contentFor: function(type, config) {
    switch (type) {
      case 'body':
        return `<svg width="168" height="53" xmlns="http://www.w3.org/2000/svg"><g fill="#919FA8" fill-rule="evenodd"><path d="M26.078 32.12a5.586 5.586 0 1 1 5.577-5.599 5.577 5.577 0 0 1-5.577 5.6M37.009 29.328a2.56 2.56 0 1 1 2.56-2.56 2.551 2.551 0 0 1-2.56 2.56M46.916 31.669a2.56 2.56 0 1 1 .051-.21c-.028.066-.028.13-.051.21M44.588 25.068a2.565 2.565 0 0 1-2.672-.992 2.558 2.558 0 0 1-.102-2.845 2.564 2.564 0 0 1 4.676.764c.072.328.081.667.027 1a2.463 2.463 0 0 1-1.925 2.073M53.932 31.402a2.547 2.547 0 0 1-2.95 2.076 2.559 2.559 0 0 1-2.064-2.965 2.547 2.547 0 0 1 2.948-2.077 2.57 2.57 0 0 1 2.128 2.716.664.664 0 0 0-.05.228M51.857 25.103a2.56 2.56 0 1 1 2.108-2.945c.034.218.043.439.027.658a2.547 2.547 0 0 1-2.135 2.287M49.954 40.113a2.56 2.56 0 1 1 .314-1.037c-.02.366-.128.721-.314 1.037M48.974 16.893a2.56 2.56 0 1 1 .97-3.487c.264.446.375.965.317 1.479a2.56 2.56 0 0 1-1.287 2.008"/><path d="M26.526 52.603c-14.393 0-26.06-11.567-26.06-25.836C.466 12.498 12.133.931 26.526.931a25.936 25.936 0 0 1 15.836 5.307l-3.167 4.117A20.962 20.962 0 0 0 17.304 8.23C10.194 11.713 5.7 18.9 5.714 26.763c-.014 7.862 4.48 15.05 11.59 18.534a20.962 20.962 0 0 0 21.89-2.127l3.168 4.123a25.981 25.981 0 0 1-15.836 5.31zM61 30.15V17.948c0-4.962 2.845-7.85 9.495-7.85 2.484 0 5.048.326 7.252.895l-.561 4.433c-2.164-.406-4.688-.691-6.53-.691-3.486 0-4.608 1.22-4.608 4.108v10.412c0 2.888 1.122 4.108 4.607 4.108 1.843 0 4.367-.284 6.53-.691l.562 4.433c-2.204.57-4.768.895-7.252.895C63.845 38 61 35.112 61 30.15zm36.808.04c0 4.068-1.802 7.81-8.493 7.81-6.69 0-8.494-3.742-8.494-7.81v-5.002c0-4.067 1.803-7.81 8.494-7.81 6.69 0 8.493 3.743 8.493 7.81v5.003zm-4.887-5.165c0-2.237-1.002-3.416-3.606-3.416s-3.606 1.18-3.606 3.416v5.328c0 2.237 1.002 3.417 3.606 3.417s3.606-1.18 3.606-3.417v-5.328zm25.79 12.568h-4.887V23.764c0-1.057-.44-1.586-1.563-1.586-1.201 0-3.325.732-5.088 1.668v13.747h-4.887V17.785h3.726l.48 1.668c2.444-1.22 5.53-2.074 7.813-2.074 3.245 0 4.407 2.318 4.407 5.857v14.357zm18.26-5.775c0 3.823-1.162 6.182-7.052 6.182-2.083 0-4.927-.488-6.73-1.139l.68-3.782c1.643.488 3.807.854 5.81.854 2.164 0 2.484-.488 2.484-1.993 0-1.22-.24-1.83-3.405-2.603-4.768-1.18-5.329-2.4-5.329-6.223 0-3.986 1.723-5.735 7.292-5.735 1.803 0 4.166.244 5.85.691l-.482 3.945c-1.482-.284-3.846-.569-5.368-.569-2.124 0-2.484.488-2.484 1.708 0 1.587.12 1.709 2.764 2.4 5.449 1.464 5.97 2.196 5.97 6.264zm4.357-14.033h4.887v13.83c0 1.057.441 1.586 1.563 1.586 1.202 0 3.325-.733 5.088-1.668V17.785h4.888v19.808h-3.726l-.481-1.667c-2.444 1.22-5.529 2.074-7.812 2.074-3.246 0-4.407-2.318-4.407-5.857V17.785zM168 37.593h-4.888V9.691L168 9v28.593z"/></g></svg>`;
      case 'root-class':
        return 'ember-loading';
    }
  },
};
