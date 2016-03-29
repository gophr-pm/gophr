'use strict';

const path    = require('path');
const chalk   = require('chalk');
const rimraf  = require('rimraf');

rimraf(path.join(__dirname, '..', '..', 'dist', '*'), (err) => {
  if (err) {
    console.log(chalk.bold.red('✗ Failed to clean frontend:'), err);
  } else {
    console.log(chalk.green('✓ Successfully cleaned the frontend'));
  }
});
