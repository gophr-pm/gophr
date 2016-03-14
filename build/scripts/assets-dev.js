'use strict';

const fs        = require('fs');
const path      = require('path');
const chalk     = require('chalk');
const watch     = require('watch');
const Imagemin  = require('imagemin');

const SRC_PATH        = path.join(__dirname, '..', '..', 'web', 'src', 'img');
const DIST_PATH       = path.join(__dirname, '..', '..', 'web', 'dist', 'img');
const SRC_IMG_PATH    = path.join('web', 'src', 'img');
const DIST_IMG_PATH   = path.join('web', 'dist', 'img');

function compressImage(imageFilePath) {
  new Promise((resolve, reject) => {
    new Imagemin()
        .src(imageFilePath)
        .dest(DIST_PATH)
        .run((err, files) => {
          if (err) {
            console.log(
                chalk.bold.red('✗ Failed to compress image:'),
                imageFilePath);
            reject(err);
          } else {
            console.log(
                chalk.green('✓ Successfully compressed image:'),
                imageFilePath);
            resolve();
          }
        });
  });
}

function compressAllImages() {
  return new Promise((resolve, reject) => {
    new Imagemin()
        .src(path.join(SRC_PATH, '*.{ico,png,svg,xml,json}'))
        .dest(DIST_PATH)
        .run((err, files) => {
          if (err) {
            console.log(
                chalk.bold.red('✗ Failed to compress image:'),
                imageFilePath);
            reject(err);
          } else {
            console.log(
                chalk.green('✓ All images compressed successfully'));
            resolve(files);
          }
        });
  });
}

function watchImages() {
  console.log('Watching changes to images...');
  watch.createMonitor(SRC_PATH, monitor => {
    monitor.on('created', compressImage);
    monitor.on('changed', compressImage);
    monitor.on('removed', (removedSourceFilePath) => {
      const removedDistFilePath = removedSourceFilePath.replace(
          SRC_IMG_PATH,
          DIST_IMG_PATH);
      fs.unlink(removedDistFilePath, err => {
        if (err) {
          console.log(
              chalk.bold.red('✗ Failed to delete image:'),
              removedDistFilePath);
        } else {
          console.log(
              chalk.green('✓ Successfully deleted image:'),
              removedDistFilePath);
        }
      });
    });
  });
}

compressAllImages().then(watchImages);
