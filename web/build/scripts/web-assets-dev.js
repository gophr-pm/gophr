'use strict';

const webAssetsLib = require('./web-assets');

const watchImages       = webAssetsLib.watchImages;
const compressAllImages = webAssetsLib.compressAllImages;

compressAllImages().then(watchImages);
